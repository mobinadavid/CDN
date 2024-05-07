package cmd

import (
	"cdn/cmd/app"
	"cdn/src/config"
	"github.com/spf13/cobra"
	"log"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cdn",
	Short: "CDN API project.",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(app.AppCmd)
}

func initConfig() {
	// Initialize configuration
	err := config.Init()
	if err != nil {
		log.Fatalf("Config Service: Failed to Initialize. %v", err)
	}
	log.Println("Config Service: Initialized Successfully.")
}
