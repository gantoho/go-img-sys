package main

import (
	"github.com/gantoho/go-img-sys/internal/app"
)

func main() {
	srv := app.New()
	defer srv.Close()

	srv.Start()
}
