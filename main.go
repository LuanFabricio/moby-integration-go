package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/moby/moby/client"

	"go-integration-tests/db"
  gContainer "go-integration-tests/container"
);

func main() {
	apiClient, err := client.New(
		client.FromEnv,
	)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

  const containerName string = "auto-integrated-tests-database"
  err = gContainer.StartPostgresContainer(
    apiClient,
    containerName,
  )

  fmt.Printf("container: %s\n", containerName)

  if err != nil {
    panic(err)
  }

  time.Sleep(gContainer.ContainerSetupDelay)

  runSetup()

  err = gContainer.StopAndKillContainer(apiClient, containerName)

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
