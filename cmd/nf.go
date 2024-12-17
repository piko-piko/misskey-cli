/*
Copyright © 2022 mikuta0407
*/
package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/mikuta0407/misskey-cli/misskey"
)

// nfCmd represents the tl command
var nfCmd = &cobra.Command{
	Use:   "nf",
	Short: "Show notifications",
	Long:  `Show notifications command`,
	Run: func(cmd *cobra.Command, args []string) {
		client := misskey.NewClient(instanceName, cfgFile)
		if err := client.GetNotifications(plainPrint, limit, sinceId); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}

var sinceId string

func init() {
	rootCmd.AddCommand(nfCmd)
	nfCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Limit display items")
	nfCmd.Flags().StringVarP(&sinceId, "sinceId", "s", "", "Since (noteId)")
}
