package internal

import (
	"errors"
	"fmt"
	"github.com/mix-go/xcli/flag"
	logic2 "github.com/go-sweets/cli/internal/logic"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var NewCmd = &cobra.Command{
	Use:     "new",
	Short:   "create a CloudWeGo microservice project",
	Long:    "Create a CloudWeGo microservice project using the sweets-layout template.\n\nUsage:\n  swe-cli new <project-name> [module-name]\n\nExamples:\n  swe-cli new helloworld\n  swe-cli new myservice github.com/myorg/myservice",
	Run:     newRun,
	Version: SkeletonVersion,
}

var (
	sweetsLayoutPath string
	name             string
	moduleName       string
)

func init() {
	// sweetsLayoutPath will be determined at runtime relative to CLI location
	name = flag.Arguments().First().String("hello")
	name = strings.ReplaceAll(name, " ", "")
	moduleName = name // Default module name to project name
}

func findSweetsLayoutPath() (string, error) {
	// First check if cached template exists and is valid
	if cache, err := logic2.GetCachedTemplate(); err == nil {
		if logic2.IsCacheValid(cache, 24*time.Hour) {
			if _, err := os.Stat(cache.Path); err == nil {
				fmt.Printf("Using cached template: %s\n", cache.Path)
				return cache.Path, nil
			}
		}
	}

	// Get the CLI executable directory to find sweets-layout
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}

	// Calculate path to sweets-layout (assume CLI is in cli/ and sweets-layout is sibling directory)
	cliDir := filepath.Dir(exePath)
	repoRoot := filepath.Dir(cliDir)
	sweetsLayoutPath := filepath.Join(repoRoot, "sweets-layout")

	// If CLI is built in cli directory, adjust path
	if filepath.Base(cliDir) == "cli" {
		sweetsLayoutPath = filepath.Join(filepath.Dir(cliDir), "sweets-layout")
	}

	// Check if sweets-layout exists
	if _, err := os.Stat(sweetsLayoutPath); err != nil {
		// Try alternative path - maybe we're in the repository root
		pwd, _ := os.Getwd()
		alternativePath := filepath.Join(pwd, "sweets-layout")
		if _, err := os.Stat(alternativePath); err != nil {
			// Try one more - maybe we're in cli directory
			alternativePath = filepath.Join(pwd, "..", "sweets-layout")
			if _, err := os.Stat(alternativePath); err != nil {
				return "", fmt.Errorf("sweets-layout template not found. Expected at: %s", sweetsLayoutPath)
			}
		}
		sweetsLayoutPath = alternativePath
	}

	// Cache the template path
	if err := logic2.CacheTemplate(SkeletonVersion, sweetsLayoutPath); err != nil {
		fmt.Printf("Warning: Failed to cache template: %v\n", err)
	}

	return sweetsLayoutPath, nil
}

func newRun(_ *cobra.Command, args []string) {
	// Parse arguments for module name
	if len(args) > 0 {
		name = args[0]
		name = strings.ReplaceAll(name, " ", "")
	}
	if len(args) > 1 {
		moduleName = args[1]
	} else {
		moduleName = name // Default to project name
	}

	if name == "" {
		fmt.Println("Project name is required. Usage: swe-cli new <project-name> [module-name]")
		return
	}

	// Find sweets-layout template
	sweetsLayoutPath, err := findSweetsLayoutPath()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Make sure you're running swe-cli from the go-sweets repository\n")
		return
	}

	fmt.Printf("Using template: %s\n", sweetsLayoutPath)

	fmt.Print(" - Generate code")
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dest := fmt.Sprintf("%s/%s", pwd, name)

	// Check if destination already exists
	if _, err := os.Stat(dest); err == nil {
		fmt.Printf("\nProject directory '%s' already exists\n", dest)
		return
	}

	if !logic2.CopyPath(sweetsLayoutPath, dest) {
		panic(errors.New(fmt.Sprintf("copy dir failed srcdir %s to %s", sweetsLayoutPath, dest)))
	}
	fmt.Println(" > ok")

	fmt.Print(" - Processing package name")
	// Replace module name in all files
	if err := logic2.ReplaceAll(dest, "github.com/go-sweets/sweets-layout", moduleName); err != nil {
		panic(errors.New("replace module name failed"))
	}

	// Update go.mod with new module name
	if err := logic2.UpdateGoMod(dest, moduleName); err != nil {
		panic(errors.New("update go.mod failed"))
	}

	// Clean up generated files that shouldn't be in template
	logic2.CleanupTemplate(dest)

	fmt.Println(" > ok")

	fmt.Print(" - Installing dependencies")
	if err := installDependencies(dest); err != nil {
		fmt.Printf("\nWarning: Failed to install dependencies: %v\n", err)
		fmt.Printf("You can install them manually by running 'go mod tidy' in the project directory\n")
	} else {
		fmt.Println(" > ok")
	}

	fmt.Printf("\nProject '%s' generated successfully!\n", name)
	fmt.Printf("Module name: %s\n", moduleName)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", name)
	fmt.Printf("  make init    # Install tools and dependencies\n")
	fmt.Printf("  make api     # Generate protobuf code\n")
	fmt.Printf("  make gen     # Generate Wire dependency injection\n")
	fmt.Printf("  make run     # Run the service\n")
}

func installDependencies(projectDir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	return cmd.Run()
}
