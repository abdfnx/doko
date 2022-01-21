package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

// Networks get networks
func (d *Docker) Networks(opt types.NetworkListOptions) ([]types.NetworkResource, error) {
	return d.NetworkList(context.TODO(), opt)
}

// InspectNetwork inspect network
func (d *Docker) InspectNetwork(name string) (types.NetworkResource, error) {
	return d.NetworkInspect(context.TODO(), name, types.NetworkInspectOptions{})
}

// DeleteNetwork -> delete network
func (d *Docker) DeleteNetwork(name string) error {
	return d.NetworkRemove(context.TODO(), name)
}
