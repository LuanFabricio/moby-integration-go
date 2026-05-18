package db

import "database/sql"

type Client struct {
  id *int
  name string
}

type Order struct{
  id *int
  client_id int
  value float32
}

var SetupTables = map[string]string{
  "clients": `
    create table clients(
        id integer
          primary key
          generated always as identity,
        name varchar(50),
        created_at timestamp default now()
    );
  `,
  "orders": "create table orders(id serial, client_id int, value money, created_at timestamp);",
}


func CheckTableExists(db *sql.DB, schema string, table string) bool {
  const query string = `
    SELECT EXISTS (
      SELECT 1
      FROM pg_tables
      WHERE schemaname = $1
      AND tablename = $2
    );`;
  row := db.QueryRow(query, schema, table)
  var exists bool = false
  if err := row.Scan(&exists); err != nil {
    panic(err)
  }

  return exists
}

func SetupDB(db *sql.DB) {
  createTables(db)

  insertTables(db)
}

func createTables(db *sql.DB) {
  for key, val := range SetupTables {
    if (!CheckTableExists(db, "public", key)) {
      _, err := db.Exec(val)
      if err != nil {
        panic(err);
      }
    }
  }
}

func insertTables(db *sql.DB) {
  clients := []Client{
    {
      name: "client1",
    },
    {
      name: "client2",
    },
  }

  for i := range len(clients) {
    client_id, err := insertClientIfNotExists(db, &clients[i])
    if err != nil && err != sql.ErrNoRows {
      panic(err)
    }
    clients[i].id = &client_id
  }

  orders := []Order{
    {
      client_id: *clients[0].id,
      value: 42.42,
    },
    {
      client_id: *clients[1].id,
      value: 10.42,
    },
  }

  for i := range len(orders) {
    insertOrder(db, &orders[i])
  }
}

func insertClientIfNotExists(db *sql.DB, client *Client) (int, error) {
    row := db.QueryRow(`
        SELECT
          id
        FROM clients
        where name = $1
        LIMIT 1;
      `,
      client.name,
    )

    var client_id int
    err := row.Scan(&client_id)
    if err == nil || err != sql.ErrNoRows {
      return client_id, err
    }

    row = db.QueryRow(`
      insert into clients (name)
        values ($1)
        returning id;`,
      client.name,
    )

    row.Scan(&client_id)
    return client_id, nil
}

func insertOrder(db *sql.DB, order *Order) int {
  row := db.QueryRow(`
    insert into orders(client_id, value)
      values ($1, $2)
      returning id;`,
    order.client_id, order.value,
  )

  row.Scan(&order.id)

  return *order.id
}
