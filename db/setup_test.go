package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/moby/moby/client"

	gContainer "go-integration-tests/container"
	"go-integration-tests/db"
)

// TODO: Fix issue with docker connection
func TestSetupDB(t *testing.T) {
  err := godotenv.Load(envFolder)
  if err != nil {
    t.Errorf("Error: %v", err)
  }

  dbConfig, err := db.GetDbConfig()
  if err != nil {
    t.Errorf("Error: %v", err)
  }

  dbConnection, err := dbConfig.NewDBConnection()
  if err != nil {
    t.Errorf("Error: %v", err)
  }

  moby, err := client.New(
    client.FromEnv,
  )
  if err != nil {
    t.Errorf("Error: %v", err)
  }
  defer moby.Close()

  const containerName string = "setup-dev-container"
  err = gContainer.StartPostgresContainer(moby, containerName)
  if err != nil {
    t.Errorf("Error: %v", err)
  }

  time.Sleep(gContainer.ContainerSetupDelay)

  db.SetupDB(dbConnection)

  for _, table := range db.SetupTables {
    fmt.Printf("Table %s\n", table)
    if(!db.CheckTableExists(dbConnection, "public", table)) {
      t.Errorf("Error: table %s not created", table)
    }
  }
}
