package main

import (
	"context"
	"log"

	"github.com/jailtonjunior94/order-aws/cmd/consumer"
	"github.com/jailtonjunior94/order-aws/cmd/server"
	"github.com/jailtonjunior94/order-aws/configs"

	"github.com/spf13/cobra"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	root := &cobra.Command{
		Use:   "outbox",
		Short: "Outbox",
	}

	server := &cobra.Command{
		Use:   "api",
		Short: "Outbox API",
		Run: func(cmd *cobra.Command, args []string) {
			server.Run(context.Background())
		},
	}

	consumers := &cobra.Command{
		Use:   "consumers",
		Short: "Outbox Consumers",
		Run: func(cmd *cobra.Command, args []string) {
			consumer.Run(context.Background(), config)
		},
	}

	root.AddCommand(server, consumers)
	if err := root.Execute(); err != nil {
		log.Fatalf("main: %v", err)
	}
}
