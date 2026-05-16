package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type DbConfig struct {
  Host 	   string
  Port 	   int
  User 	   string
  Password string
  SSLMode	 string
  Database string
};

func GetDbConfig() (*DbConfig, error) {
  var dbConf *DbConfig = nil
  port := os.Getenv("PSQL_PORT")

  portInt, err := strconv.Atoi(port)

  if err == nil {
    dbConf = &DbConfig{
      Host:     os.Getenv("PSQL_HOST"),
      Port:     portInt,
      User:     os.Getenv("PSQL_USERNAME"),
      Password: os.Getenv("PSQL_PASSWORD"),
      SSLMode:  os.Getenv("PSQL_SSL_MODE"),
      Database: os.Getenv("PSQL_DATABASE"),
    }
  }

  return dbConf, err;
}

func (c *DbConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=%s dbname=%s",
    c.Host, c.Port, c.User, c.Password, c.SSLMode, c.Database,
  )
}

func (c *DbConfig) NewDBConnection() (*sql.DB, error) {
	fmt.Printf("%s\n", c.GetConnectionString())
  db, err := sql.Open("postgres", c.GetConnectionString())

  ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()

  err = db.PingContext(ctx)
  if err != nil {
    db = nil
  }

  return db, err
}
