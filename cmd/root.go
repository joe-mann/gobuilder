package cmd

import (
	"fmt"
	"os"

	"github.com/joe-mann/gobuilder/cmd/internal/builder"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool
var excludeDirs []string

var logger logrus.Logger

var rootCmd = &cobra.Command{
	Use:   "gobuilder",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		build()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	logger = logrus.Logger{
		Out:       os.Stdout,
		Formatter: &formatter{},
		Level:     logrus.InfoLevel,
	}

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gobuilder.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")

	rootCmd.PersistentFlags().StringSliceVar(&excludeDirs, "exclude", []string{".git", "vendor", builder.BuildDir}, "directory names to exclude")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gobuilder" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gobuilder")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}