# pq-timeouts
A Postgres driver for Go that wraps lib/pq to provide read and write timeouts.

## Why?

[lib/pq](https://github.com/lib/pq) is an excellent Postgres driver written in pure Go, but it only offers support for the default 
Postgres timeouts of `connect_timeout` and `statement_timeout`. These work well, but in a high availability 
situation, they might not be enough. `statement_timeout` only works if the connection to Postgres is alive and well. 
`connect_timeout` only provides a timeout during initial connection. Once the connection is in the pool, `connect_timeout` 
doesn't apply. If the database goes down, or the network dies, the open connections will hang. Without a read or
write timeout on the connection, there is no way to recover quickly. pq-timeouts provides a way to add a timeout to 
every write and read to and from the database.

## Install

```
go get github.com/Kount/pq-timeouts
```

## Using pq-timeouts 

pq-timeouts adds two new connection string parameters: `read_timeout` and `write_timeout`. Otherwise, usage is exactly the same
as [lib/pq](https://github.com/lib/pq): 
```go
import (
  "database/sql"

  _ "github.com/Kount/pq-timeouts"
}

func main() {
  // Note: read_timeout and write_timeout are specified in milliseconds.
  db, err := sql.Open(
    "pq-timeouts", 
    "user=pqtest dbname=pqtest read_timeout=500 write_timeout=1000 sslmode=verify-full"
  )
  if err != nil {
    log.Fatal(err)
  }

  age := 21
  rows, err := db.Query("SELECT name FROM users WHERE age =$1", age)
  ...
}
```
