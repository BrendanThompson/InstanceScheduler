package main

import (
	"flag"
	"os"

	"instancescheduler/azure"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	tagsConfigPath := flag.String("config", "./tags.yaml", "path for tags config file")
	flag.Parse()

	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = log.With().Caller().Logger()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	}

	client, err := azure.NewComputeClient(subscriptionID, *tagsConfigPath)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to get compute client")
	}

	instances, err := client.ListInstances()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to get list of instances from Azure")
	}

	client.AssessInstancesAndAction(instances)
}
