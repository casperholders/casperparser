package main

import (
	"casperParser/cmd"
	"github.com/spf13/cobra/doc"
	"log"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./")
	if err != nil {
		log.Fatal(err)
	}
}
