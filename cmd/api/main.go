package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ItayTurniansky/musicwide/internal/server"
	// 1. ADD THIS IMPORT (The generated code)
	db "github.com/ItayTurniansky/musicwide/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	// 1. Get the connection string
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgresql://admin:secretpassword@localhost:5432/musicwide?sslmode=disable"
	}

	// 2. Open the connection
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// 3. Ping to verify
	if err := conn.Ping(); err != nil {
		log.Fatal("database is not responding:", err)
	}

	fmt.Println("Successfully connected to the Database!")

	// 4. PREPARE THE DATABASE OBJECT (This is the missing link!)
	// We wrap the raw connection 'conn' with the generated code 'db.New'
	database := db.New(conn)

	// 5. START THE SERVER
	// Pass the 'database' object into NewServer
	s := server.NewServer(database)

	port := "8080"
	fmt.Printf("MusicWide Server running on port %s\n", port)

	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
