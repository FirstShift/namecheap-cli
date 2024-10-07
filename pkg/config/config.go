package config

import "github.com/spf13/viper"

type NamecheapCliConfig struct {
	ApiKey     string
	UserName   string
	ClientIp   string
	UseSandbox bool
}

func GetConfig() *NamecheapCliConfig {
	return &NamecheapCliConfig{
		ApiKey:     viper.GetString("NAMECHEAP_API_KEY"),
		UserName:   viper.GetString("NAMECHEAP_USERNAME"),
		ClientIp:   viper.GetString("NAMECHEAP_CLIENT_IP"),
		UseSandbox: viper.GetBool("NAMECHEAP_SANDBOX"),
	}
}
