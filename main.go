package main

import (
	"github.com/go-sweets/cli/internal"
	"github.com/go-sweets/cli/internal/gen"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:     "swe-cli",
	Short:   "swe-cli: CloudWeGo microservice project generator",
	Long:    `swe-cli: A CLI tool for generating CloudWeGo-based microservice projects using Hertz and Kitex.`,
	Version: internal.CLIVersion,
}

func init() {
	rootCmd.AddCommand(internal.NewCmd)
	rootCmd.AddCommand(internal.UpgradeCmd)
	rootCmd.AddCommand(gen.CmdGen)
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
