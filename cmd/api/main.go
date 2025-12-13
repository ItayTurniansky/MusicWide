package main

import (
	"fmt"
	"net/http"
	"github.com/ItayTurniansky/musicwide/internal/server"
)

func main() {
	s := server.NewServer()

	port := "8080"
	fmt.Printf("MusicWide Server running on port %s\n", port)

	err := http.ListenAndServe(":"+port, s.Router)
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}