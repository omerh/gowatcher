package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/omerh/gowatcher/pkg/aws"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	log.Println("Starting watcher ....")

	tickerEnvWait, ok := os.LookupEnv("TIME_INTERVAL")
	if !ok {
		tickerEnvWait = "60s"
	}
	tickerWait, err := time.ParseDuration(tickerEnvWait)
	if err != nil {
		log.Printf("Ticking with %d", tickerWait)
		log.Println(err)
	}

	log.Printf("Ticking with %v", tickerWait)

	ticker := time.NewTicker(tickerWait)
	for range ticker.C {
		log.Println("Checking for running containers...")
		runtime()
	}
}

func runtime(){
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// Get all running containers
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	// Set filter to count running container images
	filter, ok := os.LookupEnv("FILTER")

	if !ok {
		log.Panic("Missing environment variable FILTER to filter containers")
	} else {
		log.Printf("Using filter %v", filter)
	}

	log.Printf("Total %d unfiltered containers running on host", len(containers))

	// Create filtered container list to check if the filtered containers are running
	var filteredContainers []string
	for _, container := range containers {
		if strings.Contains(container.Image, filter) {
			log.Printf("%v contains %v", container.Image, filter)
			filteredContainers = append(filteredContainers, container.ID)
		}
	}

	log.Printf("There are %d running containers running on the host", len(filteredContainers))

	if len(filteredContainers) == 0 {
		log.Printf("Filtered out %d running containers", len(filteredContainers))

		// Call terminate instance process for clean termination form Auto scale group
		aws.TerminateInstance()
	}
}
