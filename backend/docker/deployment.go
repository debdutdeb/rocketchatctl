package docker

import (
	"context"
	"log"
	"sync"

	"github.com/debdutdeb/rocketchatctl/deployment"
	"github.com/debdutdeb/rocketchatctl/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type dockerBackend struct {
	client  *client.Client
	images  []string
	network types.NetworkResource
}

const networkName = "rocketchatctl_rocketchat"

func NewDockerBackend() (*dockerBackend, error) {
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	d := dockerBackend{
		client: c,
	}
	return &d, nil
}

func (d *dockerBackend) PullArtifacts(ctx context.Context, opts deployment.HostDeploymentOptions) error {
	// TODO: handle UseExistingMongo
	// the main challenge is insufficient options, i.e. what if
	// an existing mongodb server is running on port != 3306?
	var wg sync.WaitGroup
	// TODO handle traefik
	d.images = []string{
		"registry.rocket.chat/rocketchat/rocket.chat:" + opts.Version,
		"docker.io/bitnami/mongodb:" + opts.MongoVersion,
	}
	log.Printf("images: %v\n", d.images)
	errors := utils.MultiError{}
	for _, image := range d.images {
		wg.Add(1)
		go func(image string, wg *sync.WaitGroup, errors *utils.MultiError) {
			log.Println("pulling image: " + image)
			r, err := d.client.ImagePull(ctx, image, types.ImagePullOptions{})
			if err != nil {
				errors.Append(err)
			}
			defer r.Close()
			// FIXME handle the progression right
			for _, err := r.Read([]byte{}); err == nil; {
			}
			wg.Done()
		}(image, &wg, &errors)
	}
	wg.Wait()
	if errors.Has() {
		return errors
	}
	return nil
}

func (d *dockerBackend) StartDeployment(ctx context.Context, opts deployment.HostDeploymentOptions) error {
	return nil
}

func (d *dockerBackend) ConfigureHost(ctx context.Context, opts deployment.HostDeploymentOptions) error {
	// create unit files
	// enable services
	return nil
}

func (d *dockerBackend) createNetwork(ctx context.Context) error {
	if network, err := d.client.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{}); err == nil {
		d.network = network
		return nil
	}
	resp, err := d.client.NetworkCreate(ctx, networkName, types.NetworkCreate{})
	if err != nil {
		return err
	}
	network, err := d.client.NetworkInspect(ctx, resp.ID, types.NetworkInspectOptions{})
	if err != nil {
		return err
	}
	d.network = network
	return nil
}

func (d dockerBackend) runContainer(ctx context.Context) error {
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"80/tcp": {
				{
					HostIP:   "0.0.0.0",
					HostPort: "3000",
				},
			},
		},
	}
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: make(map[string]*network.EndpointSettings),
	}
	networkConfig.EndpointsConfig["bridge"] = &network.EndpointSettings{
		NetworkID: d.network.ID,
	}
	var wg sync.WaitGroup
	var errors utils.MultiError
	for _, image := range d.images {
		wg.Add(1)
		go func(image string, ctx context.Context, wg *sync.WaitGroup, errors *utils.MultiError) {
			container, err := d.client.ContainerCreate(ctx, &container.Config{
				Image: image,
			}, hostConfig, networkConfig, &v1.Platform{}, "")
			if err != nil {
				errors.Append(err)
			}
			d.client.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
			wg.Done()
		}(image, ctx, &wg, &errors)
	}
	wg.Wait()
	return nil
}
