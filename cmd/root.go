package cmd

import "github.com/spf13/cobra"

func Execute() error {
	rootCommand := cobra.Command{
		Use:   "rocketchatctl",
		Short: "Do stuff with Rocket.Chat",
	}
	rootCommand.AddCommand(
		installCommand(),
	)
	return rootCommand.Execute()
}
