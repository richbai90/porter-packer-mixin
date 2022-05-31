//go:generate packr2
package packer

import (
	"get.porter.sh/porter/pkg/context"
)

const defaultDockerVersion = "20.10.7"
const defaultImagePath = "$HOME"

type Mixin struct {
	*context.Context
	ClientVersion string
	PackerFile string
	BuildArgs string
	DockerVersion string
	TargetOS string
	ImagePath string
	//add whatever other context/state is needed here
}

// New azure mixin client, initialized with useful defaults.
func New() (*Mixin, error) {
	return &Mixin{
		Context:       context.New(),
		ClientVersion: defaultDockerVersion,
		ImagePath: defaultImagePath,
	}, nil

}
