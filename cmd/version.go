package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of WebPathScanner",
	Long:  "All software has versions. This is WebPathScanner's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("WebPathScanner v1.0 -- HEAD")
	},
}
