package main

import (
	"fmt"
	"os"
)

// var rootCmd = &cobra.Command{
// 	Use:   "vaultenv",
// 	Short: "vault env is a secret processor using vault as its source",
// 	Long:  `an application to process secrets from vault and make it available for applications to use.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// Do Stuff Here
// 	},
// }

// Execute -
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
