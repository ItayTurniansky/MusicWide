package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ItayTurniansky/musicwide/internal/server"

	// This is the driver. The underscore "_" means "load it but don't use it directly"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Get the connection string from Environment Variables (from docker-compose)
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		// Fallback for local testing if env is missing
		dbSource = "postgresql://admin:secretpassword@localhost:5432/musicwide?sslmode=disable"
	}

	// 2. Open the connection
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// 3. Check if the database is actually alive (Ping)
	if err := conn.Ping(); err != nil {
		log.Fatal("database is not responding:", err)
	}

	fmt.Println("Successfully connected to the Database!")

	// 4. Start the Server (Passing the DB connection will happen later)
	s := server.NewServer()
	port := "8080"
	fmt.Printf("MusicWide Server running on port %s\n", port)

	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
