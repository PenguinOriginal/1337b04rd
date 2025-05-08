package utils

import "fmt"

// Displates the usage instructions for CLI
func PrintUsage() {
	fmt.Printf(`$ ./1337b04rd --help
	hacker board

	Usage:
	1337b04rd [--port <N>]  
	1337b04rd --help

	Options:
	--help       Show this screen.
	--port N     Port number.`)
}
