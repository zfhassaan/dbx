package main

import (
	"dbx/cmd"
)

/*
main is the entry point for the dbx CLI application.
It invokes the Execute function to run the root Cobra command.
*/
func main() {
	cmd.PrintBanner()
	cmd.Execute()
}
