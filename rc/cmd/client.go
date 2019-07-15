package cmd

import (
	"github.com/blushft/rc"
	"github.com/spf13/viper"
)

var (
	rcClient *rc.Client
)

func clientOpts() []rc.ClientOption {
	return []rc.ClientOption{
		rc.ServerURL(viper.GetString("server-url")),
		rc.CredFromJson(viper.GetString("cred-file")),
		rc.Debug(true),
		rc.Realtime(true),
	}
}

func createClient() {
	opts := clientOpts()
	rcClient = rc.New(opts...)
}
