package cmd

import (
	"fmt"

	"github.com/agrimel-0/rio-server/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start serving the remote-io edge server",
	Long: `Start serving the remote-io edge server on your current device.
	Serve uses the config file to automatically initialize the required pins and port.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get the necessary values from the config file
		if err := viper.ReadInConfig(); err != nil {
			return
		}

		var config server.Config

		if err := viper.Unmarshal(&config); err != nil {
			fmt.Println(err)
			return
		}

		err := server.Start(config)
		if err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
