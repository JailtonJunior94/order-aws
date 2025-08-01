package main

import (
	"context"
	"log"

	"github.com/jailtonjunior94/order-aws/cmd/server"

	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()

	root := &cobra.Command{
		Use:   "outbox",
		Short: "Outbox",
	}

	server := &cobra.Command{
		Use:   "api",
		Short: "Outbox API",
		Run: func(cmd *cobra.Command, args []string) {
			server := server.NewApiServer()
			server.Run(ctx)
		},
	}

	consumers := &cobra.Command{
		Use:   "consumers",
		Short: "Outbox Consumers",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	root.AddCommand(server, consumers)
	if err := root.Execute(); err != nil {
		log.Fatalf("main: %v", err)
	}
}
