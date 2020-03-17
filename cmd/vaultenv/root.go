package main

import (
	"fmt"
	"os"
)

// Execute -
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
