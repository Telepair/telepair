package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/telepair/telepair/core/proxy/api"
	"github.com/telepair/telepair/pkg/httpclient"
)

// APITemplateCmd represents the api template command
var APITemplateCmd = &cobra.Command{
	Use:   "api-template [name]",
	Short: "API proxy template",
	Long: `API proxy template tool for testing HTTP requests.

Examples:
  # Simple request
  ./telepair tools api-template eip

  # Use a custom template file
  ./telepair tools api-template geo -t ./configs/apis.yaml

  # Request with variables
  ./telepair tools api-template weather -v '{"city": "beijing", "lang": "zh"}'
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		template, _ := cmd.Flags().GetString("template")
		fileType := filepath.Ext(template)
		if fileType != ".yaml" && fileType != ".yml" && fileType != ".json" {
			log.Fatalf("Unsupported file type: %s", fileType)
		}
		fileType = strings.TrimPrefix(fileType, ".")
		data, err := os.ReadFile(template)
		if err != nil {
			log.Fatalf("Failed to read template file (%s): %v", template, err)
		}
		values, _ := cmd.Flags().GetString("values")
		v := make(map[string]string)
		if values != "" {
			if err := json.Unmarshal([]byte(values), &v); err != nil {
				log.Fatalf("Failed to parse values: %v", err)
			}
		}
		fmt.Printf("Running API template with name <%s> template <%s>\n", name, template)

		if err := api.RegisterAPITemplateData(fileType, data); err != nil {
			log.Fatalf("Failed to register template: %v", err)
		}
		resp, err := api.DoTemplate(name, v)
		if err != nil {
			log.Fatalf("Failed to do template: %v", err)
		}
		mediaType, body, err := httpclient.ParseResponse(resp)
		if err != nil {
			log.Fatalf("Failed to parse response: %v", err)
		}
		fmt.Printf("Response: \n\tContent-Type: %s\n\tBody: \n", mediaType)
		fmt.Println("--------------------------------")
		fmt.Println(string(body))
		fmt.Println("--------------------------------")
	},
}

func init() {
	APITemplateCmd.Flags().StringP("template", "t", "./configs/apis.yaml", "Template file, yaml or json")
	APITemplateCmd.Flags().StringP("values", "v", "", "Values for the template, json format")
}
