package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start app")
	},
}
