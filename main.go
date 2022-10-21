package main

import (
	"fmt"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
