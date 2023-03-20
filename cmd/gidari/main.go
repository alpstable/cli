// Copyright 2023 The Gidari CLI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0

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
