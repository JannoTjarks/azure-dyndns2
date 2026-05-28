package main

import (
	"embed"

	"github.com/JannoTjarks/azure-dyndns2/cmd"
	"github.com/JannoTjarks/azure-dyndns2/internal/server"
)

//go:embed docs
var openApiFiles embed.FS

func main() {
	server.OpenApiFiles = openApiFiles

	cmd.Execute()
}
