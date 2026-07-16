package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ctxqa/ctxqa/pkg/renderer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ctxqa [command]",
	Short: "ctxqa - Repository Code Auditor",
	Long: `ctxqa is a terminal CLI tool that performs comprehensive code audits
by reviewing all source code files in a repository and generating
structured reports with progress analysis, issue detection,
code quality scoring, and improvement suggestions.`,
}

func Execute() {
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if license, _ := cmd.Flags().GetBool("license"); license {
			fmt.Print(licenseText)
			os.Exit(0)
		}
	}

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		renderer.PrintWelcome(version)
		if len(args) == 0 {
			return cmd.Help()
		}
		return nil
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("ctxqa version " + version + " (" + runtime.GOOS + "/" + runtime.GOARCH + ")\n")

	rootCmd.Flags().BoolP("license", "l", false, "Show license information")
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}

const licenseText = `Apache License
Version 2.0, January 2004
http://www.apache.org/licenses/

Copyright (c) 2026 ctxqa contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`
