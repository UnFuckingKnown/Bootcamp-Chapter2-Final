package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	
)

var Conn *pgx.Conn

func DatabaseConnect() {
	var err error
	databaseUrl := "postgres://postgres:1@localhost:5432/Personal-web"

	Conn, err = pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v", err)
		os.Exit(1)
	}
	
	fmt.Println("Success connect to database")
}
