package cmd

import (
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/app"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/http"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		Start()
	},
}

func Start() {
	telemetryClearFunc := app.WithTelemetry()

	defer telemetryClearFunc()

	http.NewServer().Serve()

	app.WithGracefulShutdown()

	app.Wait()
}
