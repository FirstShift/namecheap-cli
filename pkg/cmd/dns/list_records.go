package dns

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/FirstShift/namecheapcli/pkg/namecheap"
	"github.com/FirstShift/namecheapcli/pkg/utils"
	"github.com/jedib0t/go-pretty/table"
	nc "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
	"github.com/spf13/cobra"
)

func ListDNSRecordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List DNS records",
		Long:  "List DNS records for a domain",
		Run: func(cmd *cobra.Command, args []string) {
			config := utils.GetGlobalConfig(cmd)
			domain, _ := cmd.Flags().GetString("domain")
			outputFormat, _ := cmd.Flags().GetString("output")

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

			records, err := service.ListRecords(domain)
			if err != nil {
				panic(err)
			}

			switch outputFormat {
			case "json":
				PrintJSON(records)
			case "table":
				PrintTable(records)
			case "text":
				PrintText(records)
			default:
				fmt.Println("Invalid output format, please use json, table, or text")
			}
		},
	}

	return cmd
}

func PrintJSON(records []*namecheap.NamecheapDNSRecord) {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(data))

}

func PrintTable(records []*namecheap.NamecheapDNSRecord) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Host", "Type", "Value", "TTL"})

	for _, record := range records {
		t.AppendRow([]interface{}{record.Host, record.Type, record.Value, record.TTL})
	}

	t.SetStyle(table.StyleLight)
	t.Render()
}

func PrintText(records []*namecheap.NamecheapDNSRecord) {
	for _, record := range records {
		fmt.Printf("%s %s %s %d\n", record.Host, record.Type, record.Value, record.TTL)
	}
}
