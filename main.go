package main

import (
	"munch-o-matic/cmd"
)

// btoi is a helper function to convert bool to int
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func main() {
	cmd.Execute()
}
