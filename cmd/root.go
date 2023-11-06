package cmd

import (
	"fmt"
	"log"
	"os"

	"munch-o-matic/client"
	. "munch-o-matic/client/types"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfg Config
var cli *client.RestClient

var rootCmd = &cobra.Command{
	Use:   "munch-o-matic",
	Short: "TasteNext API Client",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatal(err)
		}

		cli, err = client.NewClient(cfg)
		if err != nil {
			fmt.Println("Failed to init client:", err)
			fmt.Println("Please check your munch-o-matic configuration")
			log.Fatal(err)
		}

		// Update configuration
		viper.Set("SessionCredentials.SessionID", cli.SessionID)
		viper.Set("SessionCredentials.UserId", cli.UserId)
		viper.Set("SessionCredentials.CustomerId", cli.CustomerId)
		err = viper.WriteConfig()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.munch-o-matic.yaml)")
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

		// Search config in home directory with name ".munch-o-matic" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".munch-o-matic")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %v\n\n", viper.ConfigFileUsed())
	}
}
