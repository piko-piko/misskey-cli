/*
Copyright © 2022 mikuta0407
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/mikuta0407/misskey-cli/config"
	"github.com/spf13/cobra"
)

var cfgFile string
var instanceName string
var plainPrint bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "misskey-cli",
	Short: "Misskey CLI Client",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func configFile() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "misskey-cli.toml")
}

func init() {

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", configFile(), "config file")
	rootCmd.PersistentFlags().StringVarP(&instanceName, "instance", "i", "", "connect instance name(not host name)")
	rootCmd.PersistentFlags().BoolVarP(&plainPrint,"plain", "p", false, "Toggle plain print(for use with pipes)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

var configs config.Config
