package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	// https://github.com/spf13/cobra
	// https://github.com/spf13/viper
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "cobra",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
)

func main() {
	// initializeViper()

	Execute()
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

// func initializeViper() {
// 	viper.SetConfigName("vault") // name of config file (without extension)
// 	viper.SetConfigType("yaml")  // REQUIRED if the config file does not have the extension in the name
// 	viper.AddConfigPath(".")     // optionally look for config in the working directory
// 	err := viper.ReadInConfig()  // Find and read the config file
// 	if err != nil {              // Handle errors reading the config file
// 		panic(fmt.Errorf("Fatal error config file: %s \n", err))
// 	}
// }

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")

	// rootCmd.AddCommand(addCmd)
	// rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(renderCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
		er("unable to find config file to process. exiting.")
	}
}
