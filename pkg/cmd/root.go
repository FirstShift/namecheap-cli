package cmd

import (
	"fmt"

	"github.com/FirstShift/namecheapcli/pkg/cmd/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	NamecheapAPIKey   string
	NamecheapAPIUser  string
	NamecheapClientIP string
	NamecheapSandbox  bool
	OutputFormat      string
)

// NewBaseCommand creates a new base command
func NewBaseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namecheap",
		Short: "Namecheap CLI",
		Long:  "A CLI for managing Firstshift's Namecheap account",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add global flags
	cmd.PersistentFlags().StringVarP(&NamecheapAPIKey, "key", "k", "", "Namecheap API Key. Can also be set via NAMECHEAP_API_KEY")
	cmd.PersistentFlags().StringVarP(&NamecheapAPIUser, "user", "u", "", "Namecheap API User. Can also be set via NAMECHEAP_API_USER")
	cmd.PersistentFlags().StringVarP(&NamecheapClientIP, "ip", "i", "", "Namecheap Client IP. Can also be set via NAMECHEAP_CLIENT_IP")
	cmd.PersistentFlags().BoolVarP(&NamecheapSandbox, "sandbox", "s", false, "Use Namecheap Sandbox. Can also be set via NAMECHEAP_SANDBOX")
	cmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "text", "Output format (table, json, text)")

	// DNS Commands
	createDnsCmd := dns.NewCreateDNSRecorcCmd()
	deleteDnsCmd := dns.NewDeleteDNSRecorcCmd()
	listDnsCmd := dns.ListDNSRecordsCmd()
	baseDnsCmd := dns.NewDNSCommand([]*cobra.Command{createDnsCmd, deleteDnsCmd, listDnsCmd})

	cmd.AddCommand(baseDnsCmd)

	return cmd
}

// init initializes the command
func init() {
	cobra.OnInitialize(InitConfig)
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Execute runs the command
func Execute() {
	baseCmd := NewBaseCommand()

	err := baseCmd.Execute()
	if err != nil {
		panic(err)
	}
}
