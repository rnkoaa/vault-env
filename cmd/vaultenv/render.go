package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "render values from vault",
	Long:  `render values from vault using a configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rendering vault values to file.")
		os.Exit(0)
	},
}
