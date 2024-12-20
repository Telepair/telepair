# Telepair Tools

The `tools` command provides various utilities for development and testing purposes.

```bash
Usage:
  telepair tools [command]

Available Commands:
  api          API Proxy
  api-template API proxy template
```

## API

```bash
./telepair tools api --help
API Proxy tool for testing HTTP requests.

Examples:
  # Simple GET request
  ./telepair tools api https://httpbin.org/get -H "Accept: application/json"

  # POST request with headers
  ./telepair tools api https://httpbin.org/post -X POST -H "Accept: application/json" --data '{"name":"John", "age":30}'

  # Request with timeout
  ./telepair tools api https://httpbin.org/post -X POST -t 30s

Usage:
  telepair tools api [url] [flags]

Flags:
  -d, --data string      HTTP request body
  -H, --header strings   HTTP headers (can be specified multiple times)
  -h, --help             help for api
  -X, --method string    HTTP method (GET, POST, etc.) (default "GET")
  -t, --timeout string   Timeout for the request (default "30s")
```

## API Template

```bash
➜ ./telepair tools api-template --help
API proxy template tool for testing HTTP requests.

Examples:
  # Simple request
  ./telepair tools api-template eip

  # Use a custom template file
  ./telepair tools api-template geo -t ./configs/apis.yaml

  # Request with variables
  ./telepair tools api-template weather -v '{"city": "beijing", "lang": "zh"}'

Usage:
  telepair tools api-template [name] [flags]

Flags:
  -h, --help              help for api-template
  -t, --template string   Template file, yaml or json (default "./configs/apis.yaml")
  -v, --values string     Values for the template, json format
```
