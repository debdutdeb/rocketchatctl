package cmd

import (
	"context"

	"github.com/debdutdeb/rocketchatctl/actions"
	"github.com/debdutdeb/rocketchatctl/backend/docker"
	"github.com/debdutdeb/rocketchatctl/deployment"
	"github.com/spf13/cobra"
)

const dockerBackend = true // TODO change this backend

func installCommand() *cobra.Command {
	opts := deployment.HostDeploymentOptions{}
	install := &cobra.Command{
		Use:     "install",
		Aliases: []string{"deploy"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO changeme
			if dockerBackend {
				backend, err := docker.NewDockerBackend()
				if err != nil {
					return err
				}
				return actions.InstallAllResources(context.TODO(), backend, opts)
			}
			/*
				step by step
				1. verify versions (Rockret.Chat & MongoDB)
				2. pull images
				3. Create netwoek
				4. Start containers
				5. Set up systemd services
			*/
			return nil
		},
	}
	install.Flags().StringVar(&opts.RootUrl, "root-url", "http://localhost:3000", "Root url (the url you use to navigate to your instance) value")
	// FIXME
	install.Flags().StringVar(&opts.Version, "version", "latest", "RocketChat server version to install")
	// FIXME
	install.Flags().StringVar(&opts.MongoVersion, "mongo-version", "", "MongoDb version to use")
	install.Flags().StringVar(&opts.RegistrationToken, "reg-token", "", "This value can be obtained from https://cloud.rocket.chat to automatically register your workspace during startup")
	install.Flags().Int16Var(&opts.Port, "port", 3000, "port for the RocketChat server")
	install.Flags().BoolVar(&opts.UseExistingMongo, "use-mongo", false, "in case mongo installed, and storage engine configured is wiredTiger, skip mongo installation but uses systems mongo for RocketChat server database")
	install.Flags().BoolVar(&opts.BindLoopback, "bind-loopback", false, "value=(true|false) set to false to prevent from bind RocketChat server to loopback interface when installing a webserver")
	return install
}
