package dns

import (
	"fmt"

	"github.com/FirstShift/namecheapcli/pkg/namecheap"
	"github.com/FirstShift/namecheapcli/pkg/utils"
	nc "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
	"github.com/spf13/cobra"
)

func NewDeleteDNSRecorcCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [name] [type]",
		Short: "Delete a DNS record",
		Long:  "Delete a DNS record",
		Run: func(cmd *cobra.Command, args []string) {
			config := utils.GetGlobalConfig(cmd)
			domain, _ := cmd.Flags().GetString("domain")

			if domain == "" {
				cmd.Help()
				return
			}

			client := nc.NewClient(&nc.ClientOptions{
				ApiKey:     config.ApiKey,
				ApiUser:    config.UserName,
				UserName:   config.UserName,
				UseSandbox: config.UseSandbox,
				ClientIp:   config.ClientIp,
			})

			service := namecheap.New(client)

			name := args[0]
			recordType := args[1]

			if name == "" || recordType == "" {
				cmd.Help()
				return
			}

			dryRun, _ := cmd.Flags().GetBool("dry-run")

			record := &namecheap.NamecheapDNSRecord{
				Host: name,
				Type: recordType,
			}

			err := service.DeleteRecord(domain, record, dryRun)
			if err != nil {
				fmt.Printf("error deleting record: %s\n", err)
			}
		},
	}

	// Add flags
	cmd.Flags().BoolP("dry-run", "", false, "Dry run")

	return cmd
}
