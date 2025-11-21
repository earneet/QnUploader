package main

import (
	"log"
	"os"

	"qiniu-uploader/internal/cli"
)

func main() {
	app := cli.NewApp()

	if err := app.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}