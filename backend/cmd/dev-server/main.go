package main

import (
	"os"

	"careco/backend/entrypoint/server"
)

func main() { os.Exit(server.Start(build)) }
