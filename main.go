package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/moby/moby/client"

	"go-integration-tests/db"
  g_container "go-integration-tests/container"
);

func main() {
	apiClient, err := client.New(
		client.FromEnv,
	)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	// showContainers(apiClient)
	// showImages(apiClient)

  containerName := "auto-integrated-tests-database"
  err = g_container.StartPostgresContainer(
    apiClient,
    containerName,
  )

  fmt.Printf("container: %v\n", containerName)

  if err != nil {
    panic(err)
  }

  time.Sleep(g_container.ContainerSetupDelay)

  runSetup()

  err = g_container.StopAndKillContainer(apiClient, containerName)

  if err != nil {
    panic(err)
  }
}

func runSetup() {
  err := godotenv.Load()
  if err != nil {
    panic(err)
  }

  dbConfig, err := db.GetDbConfig()
  fmt.Printf("db: %v\n", dbConfig)
  if err != nil {
    panic(err)
  }

  dbConnection, err := dbConfig.NewDBConnection()
  if err != nil {
    panic(err)
  }

  db.SetupDB(dbConnection)

}

func showContainers(apiClient client.APIClient) {
	result, err := apiClient.ContainerList(
		context.Background(),
		client.ContainerListOptions{All: true},
	)
	if err != nil {
		panic(err)
	}

	for _, ctr := range result.Items {
		name := strings.Join(ctr.Names, ", ")
		fmt.Printf("%s | %s | %s | (status: %s)\n", ctr.ID, name, ctr.Image, ctr.Status)
	}
}

func showImages(apiClient client.APIClient) {
	result, err := apiClient.ImageList(
		context.Background(),
		client.ImageListOptions{},
	)
	if err != nil {
		panic(err)
	}

	for _, img := range result.Items {
		tags := strings.Join(img.RepoTags, ", ")
		fmt.Printf("%s | %s | %d\n", img.ID, tags, img.Containers)
	}
}
