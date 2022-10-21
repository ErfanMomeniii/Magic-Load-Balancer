package cmd

import (
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/app"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		app.Start()
	},
}
