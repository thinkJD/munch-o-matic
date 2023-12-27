package cmd

import (
	"fmt"
	"log"
	"os"

	"munch-o-matic/client"
	"munch-o-matic/core"

	"github.com/common-nighthawk/go-figure"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfg AppConfig
var cli *client.RestClient

var rootCmd = &cobra.Command{
	Use:   "munch-o-matic",
	Short: "TasteNext API Client",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatal(err)
		}

		cli, err = client.NewClient(cfg.Client)
		if err != nil {
			fmt.Println("Failed to init client:", err)
			fmt.Println("Please check your munch-o-matic configuration")
			log.Fatal(err)
		}

		// Update configuration
		viper.Set("client.SessionCredentials.SessionID", cli.SessionID)
		viper.Set("client.SessionCredentials.UserId", cli.UserId)
		viper.Set("client.SessionCredentials.CustomerId", cli.CustomerId)
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

// Contains each package config
type AppConfig struct {
	Client client.Config
	Core   core.Config
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.SetUsageTemplate(customUsageTemplate())
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

	// Read configuration from environment variables
	viper.SetEnvPrefix("mom")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %v\n\n", viper.ConfigFileUsed())
	}
}

func customUsageTemplate() string {
	asciiArt := figure.NewFigure("munch-o-matic", "ogre", true).String()
	return asciiArt + `
Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}
