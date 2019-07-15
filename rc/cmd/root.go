package cmd

import (
	"log"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "rc - rocket.chat terminal client",
}

func init() {
	cobra.OnInitialize(configure)
	rootCmd.PersistentFlags().StringP("rc-server-url", "s", "http://localhost:3000", "rocket.chat server url")
	rootCmd.PersistentFlags().StringP("rc-username", "u", "", "rocket.chat username")
	rootCmd.PersistentFlags().StringP("rc-password", "p", "", "rocket.chat password")
	rootCmd.PersistentFlags().StringP("user-id", "i", "", "rocket.chat user id (for use with token)")
	rootCmd.PersistentFlags().StringP("token", "t", "", "rocket.chat user access token")
	rootCmd.PersistentFlags().StringP("cred-file", "f", "token.json", "path to a json file with user credentials (u & p or token)")

	viper.BindPFlag("server-url", rootCmd.PersistentFlags().Lookup("server-url"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("user-id", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("cred-file", rootCmd.PersistentFlags().Lookup("cred-file"))

	rootCmd.AddCommand(integrationsCmd)
}

func configure() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(".rccfg")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("unable to load config file: %v", err)
	}
}

func run(cmd *cobra.Command, args []string) error {
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
