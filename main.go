package main

import (
	"os"

	"github.com/neticdk/go-jsonnetic/cmd"
)

var version = "HEAD"

func main() {
	os.Exit(cmd.Execute(version))
}
