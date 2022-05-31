package packer

import (
	"get.porter.sh/porter/pkg/exec/builder"
	yaml "gopkg.in/yaml.v2"
	"fmt"
	"path/filepath"
)

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the packer mixin in porter.yaml
// mixins:
// - packer:
//	  clientVersion: "v0.0.0"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`
	PackerFile string `yaml:"packerFile"`
	BuildArgs string `yaml:"buildArgs,omitempty"`
	TargetOS string `yaml:"targetOS"`
	ImagePath string `yaml:"imagePath,omitempty"`
}

// Install docker in the container. We will pass the host docker daemon to the container
// Pull host_exec to facilitate executing commands on the host
// add install and uninstall scripts
var dockerfileLines = `ARG DOCKER_VERSION=%s
RUN apt-get update && apt-get install -y curl && \
	curl -o docker.tgz https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VERSION}.tgz && \
    tar -xvf docker.tgz && \
    mv docker/docker /usr/bin/docker && \
    chmod +x /usr/bin/docker && \
    rm docker.tgz && \
	docker pull richbai90/host_exec:%s && \
	mkdir -p /tmp/packer/ /tmp/user_scripts && \
	echo $'%s' > /tmp/packer/Dockerfile && \
	docker build /tmp/packer -t packer && \
	echo $'%s' > /install.sh && \
	echo $'%s' > /uninstall.sh && \
	chmod +x /install.sh /uninstall.sh && \
	rm -rf /tmp/packer
`

// The install container should create a new Dockerfile that installs packer and qemu
var packerDockerfileLines = `FROM jkz0/qemu \n\
RUN apt-get -y install curl && \ \n\
	curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add - && \ \n\
	apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main" && \ \n\
	apt-get update && apt-get install packer \n\
	COPY %s /tmp/packer \n\
	WORKDIR /tmp/packer \n\
	ENTRYPOINT packer build %s && mv /tmp/packer/ouput*/* /image \n\
`

var install_sh = `#! /bin/sh \n\
rm -rf /tmp/user_scripts/*
echo echo 'docker cp vm:/image %s' > /tmp/user_scripts/copy_img.sh
docker run -d -name vm packer \n\
docker run --rm richbai90/host_exec -v /tmp/user_scripts:/user_scripts \n\
`

var uninstall_sh = `#! /bin/sh \n\
docker volume rm image \n\
docker container stop vm \n\
docker container rm vm \n\
docker image rm packer \n\
rm -rf /tmp/user_scripts/* \n\
echo 'rm -rf %s/image' > /tmp/user_scripts/rm_img.sh \n\
docker run --rm richbai90/host_exec -v /tmp/user_scripts:/user_scripts \n\
`

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build() error {

	// Create new Builder.
	var input BuildInput

	packerDir, _ := filepath.Split(m.PackerFile)

	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	suppliedClientVersion := input.Config.ClientVersion

	if suppliedClientVersion != "" {
		m.ClientVersion = suppliedClientVersion
	}

	if input.Config.ClientVersion != "" {
		m.DockerVersion = input.Config.ClientVersion
	}

	// when the packer container is run it should build the vm and move the output to the image folder for extraction with a port cmd build step
	packerDockerfileLines := fmt.Sprintf(packerDockerfileLines, packerDir, m.PackerFile)
	install_sh := fmt.Sprintf(install_sh, m.ImagePath)
	uninstall_sh := fmt.Sprintf(uninstall_sh, m.ImagePath)

	fmt.Fprintf(m.Out, dockerfileLines, m.DockerVersion, m.TargetOS, packerDockerfileLines, install_sh, uninstall_sh)

	// Example of pulling and defining a client version for your mixin
	// fmt.Fprintf(m.Out, "\nRUN curl https://get.helm.sh/helm-%s-linux-amd64.tar.gz --output helm3.tar.gz", m.ClientVersion)

	return nil
}
