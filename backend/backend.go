package backend

import (
	"context"

	"github.com/debdutdeb/rocketchatctl/deployment"
)

type Backend interface {
	PullArtifacts(ctx context.Context, opts deployment.HostDeploymentOptions) error
	StartDeployment(ctx context.Context, opts deployment.HostDeploymentOptions) error
	ConfigureHost(ctx context.Context, opts deployment.HostDeploymentOptions) error
}
