package main

import (
	"errors"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

type CLI struct {
	Serve ServeCmd `cmd:"" help:"Start server."`
}

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal("could not load .env file")
	}

	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("config-sentinel"),
		kong.Description("A microservice to apply rule based modifications to json/yaml files."),
		kong.UsageOnError(),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
