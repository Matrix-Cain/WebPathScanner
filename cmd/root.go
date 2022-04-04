package cmd

import (
	"WebPathScanner/core"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "fofa",
	Short: `
                   ___________                            
    ______________ __  /___  /________ _________ ________ 
    ___  __ \  __ '/  __/_  __ \_  __ '__ \  __ '/__  __ \
	__  /_/ / /_/ // /_ _  / / /  / / / / / /_/ /__  /_/ /
	_  .___/\__,_/ \__/ /_/ /_//_/ /_/ /_/\__,_/ _  .___/ 
	/_/                                          /_/ 
`,
}

func init() {
	rootCmd.AddCommand(versionCmd)

	scanCmd.PersistentFlags().StringVarP(&core.GlobalConfig.FileName, "name", "n", "", "Naming results file")
	scanCmd.PersistentFlags().StringVarP(&core.GlobalConfig.ConfigFilePath, "config", "c", "", "Specify path for custom toml config file")
	scanCmd.PersistentFlags().StringVarP(&core.GlobalConfig.Proxy, "proxy", "", "", "Specify proxy for requests")
	scanCmd.PersistentFlags().BoolVarP(&core.GlobalConfig.Save, "save", "", false, "Saving results to a file")
	scanCmd.PersistentFlags().IntVarP(&core.GlobalConfig.ThreadNum, "thread", "", 10, "Specify the worker pool size")
	scanCmd.PersistentFlags().IntVarP(&core.GlobalConfig.Mode, "mode", "", 0, "Specify the scan mode") //default dict mode
	scanCmd.PersistentFlags().IntVarP(&core.GlobalConfig.RandomSleep, "sleep", "", 0, "Specify sleep time between requests, the delay time will be 0~sleep time you specify")

	scanCmd.AddCommand(scanUrlCmd)

	scanCmd.AddCommand(scanFileCmd)

	rootCmd.AddCommand(scanCmd)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
