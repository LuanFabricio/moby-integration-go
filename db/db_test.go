package db_test

import (
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/moby/moby/client"

  "go-integration-tests/db"
	gContainer "go-integration-tests/container"
)

const envFolder string = "../.env"

func _TestGetDbConfig(t *testing.T) {
  err := godotenv.Load(envFolder)
  if err != nil {
    t.Errorf("Error: %v", err)
  }

  db_config, err := db.GetDbConfig()

  if err != nil {
    t.Errorf("Error: %v\n", err)
  }

  if db_config == nil {
    t.Error("Error: db_config returns a nil value")
  }
}

func TestGetConnectionString(t *testing.T) {
  err := godotenv.Load(envFolder)
  if err != nil {
    t.Errorf("Error: %v\n", err)
  }

  dbConfig, err := db.GetDbConfig()
  if err != nil {
    t.Errorf("Error: %v\n", err)
  }

  connString := dbConfig.GetConnectionString()

  valuesToCheck := []string{
    "host", "port", "user", "password", "ssl", "dbname",
  }

  for _, val := range valuesToCheck {
    if !strings.Contains(connString, val) {
      t.Errorf("Error: connString should contain a %s parameter", val)
    }
  }
}

func TestNewDBConnection(t *testing.T) {
  err := godotenv.Load(envFolder)
  if err != nil {
    t.Errorf("Error: %v", err)
  }

  dbConfig, err := db.GetDbConfig()
  if err != nil {
    t.Errorf("Error: %v", err)
  }

  const containerName string = "test-container"
  moby, err := client.New(
		client.FromEnv,
	)
	if err != nil {
    t.Errorf("Error: %v", err)
	}
  defer moby.Close()

  err = gContainer.StartPostgresContainer(moby, containerName)
  if err != nil {
    t.Errorf("Error: %v", err)
  }

  time.Sleep(gContainer.ContainerSetupDelay)
  _, err = dbConfig.NewDBConnection()
  if err != nil {
    t.Errorf("Error: %v", err)
  }

  err = gContainer.StopAndKillContainer(moby, containerName)
  if err != nil {
    t.Errorf("Error: %v", err)
  }

  time.Sleep(gContainer.ContainerSetupDelay)
}
