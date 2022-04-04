package cmd

import (
	"WebPathScanner/controller"
	"WebPathScanner/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for both one url or multiple urls from file",
	Long:  "Scan for individual url or import urls from file",
}

var scanUrlCmd = &cobra.Command{
	Use:   "url",
	Short: "Scan for target url",
	Long:  "Scan for individual url",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		core.GlobalConfig.TargetType = "url"
		core.GlobalConfig.Target = args[0]
		log.Infoln("Scanning for single url..")
		controller.PipeLine()
	},
}

var scanFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Scan for multiple urls",
	Long:  "Scan for multiple urls by importing urls from file",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		core.GlobalConfig.TargetType = "file"
		log.Infoln("Scanning targets from file..")
		controller.PipeLine()
		//fileName := args[0]
		//utility.LoadUrl(true, fileName, threadNum, saveResult, saveFileName)
	},
}
