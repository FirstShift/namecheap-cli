package dns

import (
	"strconv"

	"github.com/FirstShift/namecheapcli/pkg/namecheap"
	"github.com/FirstShift/namecheapcli/pkg/utils"
	nc "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
	"github.com/spf13/cobra"
)

// ae60f5dd724aa48dbb4909fffb7f0dc8-1019ae06a0f0287f.elb.us-west-1.amazonaws.com.

func NewCreateDNSRecorcCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name] [type] [value] [ttl]",
		Short: "Create a DNS record",
		Long:  "Create a DNS record",
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
			value := args[2]

			if name == "" || recordType == "" || value == "" {
				cmd.Help()
				return
			}

			dryRun, _ := cmd.Flags().GetBool("dry-run")

			// convert ttl to int
			ttl, err := strconv.Atoi(args[3])
			if err != nil {
				panic(err)
			}

			record := &namecheap.NamecheapDNSRecord{
				Host:  name,
				Type:  recordType,
				Value: value,
				TTL:   ttl,
			}

			err = service.CreateRecord(domain, record, dryRun)
			if err != nil {
				panic(err)
			}
		},
	}

	// Add flags
	cmd.Flags().BoolP("dry-run", "", false, "Dry run")

	return cmd
}
