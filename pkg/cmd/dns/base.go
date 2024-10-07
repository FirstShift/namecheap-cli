package dns

import "github.com/spf13/cobra"

var domain string

func NewDNSCommand(subCmds []*cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "DNS management",
		Long:  "DNS management",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	// Add base flags
	cmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "Top level domain to manage")

	// Add subcommands
	for _, subCmd := range subCmds {
		cmd.AddCommand(subCmd)
	}

	return cmd
}
