/*
Copyright © 2024 Liys <liys87x@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/telepair/telepair/core/proxy/api"
	"github.com/telepair/telepair/pkg/httpclient"
)

// toolsCmd represents the tools command
var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Tools for development",
}

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api [url]",
	Short: "API Proxy",
	Long: `API Proxy tool for testing HTTP requests.

Examples:
  # Simple GET request
  ./telepair tools api https://httpbin.org/get -H "Accept: application/json"

  # POST request with headers
  ./telepair tools api https://httpbin.org/post -X POST -H "Accept: application/json" --data '{"name":"John", "age":30}'

  # Request with timeout
  ./telepair tools api https://httpbin.org/post -X POST -t 30s`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		method, _ := cmd.Flags().GetString("method")
		url := args[0]
		headers, _ := cmd.Flags().GetStringSlice("header")
		data, _ := cmd.Flags().GetString("data")
		timeout, _ := cmd.Flags().GetString("timeout")
		fmt.Printf("Running API command with:\n\tMethod: %s\n\tURL: %s\n\tHeaders: %v\n\tTimeout: %s\n",
			method, url, headers, timeout)
		c := api.API{
			Name:   "cli-tools-api",
			Method: method,
			URL:    url,
			Body:   data,
			Config: api.Config{
				Timeout: cast.ToDuration(timeout),
			},
		}
		if len(headers) > 0 {
			c.Headers = make(map[string]string)
			for _, header := range headers {
				parts := strings.SplitN(header, ":", 2)
				if len(parts) == 2 {
					c.Headers[parts[0]] = parts[1]
				}
			}
		}
		if err := c.Parse(); err != nil {
			log.Fatalf("Error parsing API: %v", err)
		}
		resp, err := c.Do()
		if err != nil {
			log.Fatalf("Error doing API: %v", err)
		}
		mediaType, body, err := httpclient.ParseResponse(resp)
		if err != nil {
			log.Fatalf("Error parsing response: %v", err)
		}
		fmt.Printf("Response: \n\tContent-Type: %s\n\tBody: \n", mediaType)
		fmt.Println("--------------------------------")
		fmt.Println(string(body))
		fmt.Println("--------------------------------")
	},
}

func init() {
	rootCmd.AddCommand(toolsCmd)
	toolsCmd.AddCommand(apiCmd)

	apiCmd.Flags().StringP("method", "X", "GET", "HTTP method (GET, POST, etc.)")
	apiCmd.Flags().StringSliceP("header", "H", []string{}, "HTTP headers (can be specified multiple times)")
	apiCmd.Flags().StringP("data", "d", "", "HTTP request body")
	apiCmd.Flags().StringP("timeout", "t", "30s", "Timeout for the request")
}
