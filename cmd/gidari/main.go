package main

import (
	"log"

	"github.com/alpstable/cli/internal/operation"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:     "gidari",
		Short:   "gidari is a web-to-storage proxy",
		Version: "pre-alpha",
	}

	operation.AddWebServerCommand(cmd)

	if err := cmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
