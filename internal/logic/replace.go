package logic

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"encoding/json"
	"time"
)

func ReplaceAll(root, old, new string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// 替换内容
			text, err := ReadAll(path)
			if err != nil {
				return err
			}
			str := string(text)
			reg := regexp.MustCompile(old)
			str = reg.ReplaceAllString(str, new)
			if err := WriteToFile(path, str); err != nil {
				return err
			}
		}

		return err
	})
	return err
}

func ReplaceMod(root string) error {
	path := fmt.Sprintf("%s/go.mod", root)
	text, err := ReadAll(path)
	if err != nil {
		return err
	}
	str := string(text)
	reg := regexp.MustCompile(`(replace \([\s\S]*?\))`)
	str = reg.ReplaceAllString(str, "")
	if err := WriteToFile(path, str); err != nil {
		return err
	}
	return nil
}

// TemplateCache represents cached template metadata
type TemplateCache struct {
	Version   string    `json:"version"`
	Path      string    `json:"path"`
	CachedAt  time.Time `json:"cached_at"`
	Checksum  string    `json:"checksum"`
}

// GetCacheDir returns the cache directory for templates
func GetCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(homeDir, ".swe-cli", "cache")
	return cacheDir, nil
}

// CacheTemplate stores template information in cache
func CacheTemplate(version, path string) error {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	cache := TemplateCache{
		Version:  version,
		Path:     path,
		CachedAt: time.Now(),
		Checksum: "", // TODO: Add checksum calculation
	}

	cacheFile := filepath.Join(cacheDir, "sweets-layout.json")
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return WriteToFile(cacheFile, string(data))
}

// GetCachedTemplate retrieves cached template information
func GetCachedTemplate() (*TemplateCache, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}

	cacheFile := filepath.Join(cacheDir, "sweets-layout.json")
	data, err := ReadAll(cacheFile)
	if err != nil {
		return nil, err
	}

	var cache TemplateCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

// IsCacheValid checks if the cached template is still valid
func IsCacheValid(cache *TemplateCache, maxAge time.Duration) bool {
	return time.Since(cache.CachedAt) < maxAge
}

func UpdateGoMod(root, moduleName string) error {
	path := fmt.Sprintf("%s/go.mod", root)
	text, err := ReadAll(path)
	if err != nil {
		return err
	}
	str := string(text)

	// Replace module name in the first line
	lines := strings.Split(str, "\n")
	if len(lines) > 0 {
		lines[0] = fmt.Sprintf("module %s", moduleName)
	}

	// Remove any replace directives (commented out)
	reg := regexp.MustCompile(`//replace \([\s\S]*?\)`)
	updatedStr := reg.ReplaceAllString(strings.Join(lines, "\n"), "")

	if err := WriteToFile(path, updatedStr); err != nil {
		return err
	}
	return nil
}

func CleanupTemplate(root string) error {
	// Remove build artifacts and temporary files
	filesToRemove := []string{
		"sweets-app",
		"test-build",
		"http",
		"rpc",
		"bin/",
		".git/",
	}

	for _, file := range filesToRemove {
		path := filepath.Join(root, file)
		if _, err := os.Stat(path); err == nil {
			_ = os.RemoveAll(path)
		}
	}

	return nil
}
