package container

import (
	"context"
	"fmt"
	"net/netip"
	"slices"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

const ContainerSetupDelay time.Duration = time.Second * 3

func StartPostgresContainer(moby *client.Client, containerName string) error {
  if !containerExists(moby, containerName) {
    _, err := createPostgresContainer(moby, containerName)
    if err != nil {
      return err
    }
  }

  time.Sleep(ContainerSetupDelay)

  return startContainer(moby, containerName)
}

func StopAndKillContainer(moby *client.Client, containerName string) error {
  stopResult, err := moby.ContainerStop(
    context.Background(),
    containerName,
    client.ContainerStopOptions{},
  )
  if err != nil {
    panic(err)
  }

  fmt.Printf("stopResult: %v\n", stopResult)

  time.Sleep(ContainerSetupDelay)

  _, err = moby.ContainerRemove(
    context.Background(),
    containerName,
    client.ContainerRemoveOptions{
      RemoveVolumes: true,
      Force: true,
    },
  )

  return nil
}

func containerExists(moby *client.Client, containerName string) bool {
  res, err := moby.ContainerList(
    context.Background(),
    client.ContainerListOptions{ All: true },
  )

  fmt.Printf("err: %v\n", err)

  if containerName[0] != '/' {
    containerName = "/" + containerName
  }
  for i := range len(res.Items) {
    item := res.Items[i]

    if slices.Contains(item.Names, containerName) {
      return true
    }
  }

  return false
}

func createPostgresContainer(moby *client.Client, containerName string) (string, error) {
  portMap := network.PortMap{}
  port, err := network.ParsePort("5432/tcp")

  portBind := network.PortBinding{
    HostIP: netip.AddrFrom4([4]byte{0, 0, 0, 0}),
    HostPort: "5432",
  }
  portMap[port] = []network.PortBinding{portBind}

  createResult, err := moby.ContainerCreate(
    context.Background(),
    client.ContainerCreateOptions{
      Config: &container.Config{
        Env: []string{ "POSTGRES_PASSWORD=pwd" },
      },
      HostConfig: &container.HostConfig{
        PortBindings: portMap,
      },
      Name: containerName,
      Image: "postgres:latest",
    },
  )

  return createResult.ID, err
}

func startContainer(moby *client.Client, containerName string) error {
  _, err := moby.ContainerStart(
    context.Background(),
    containerName,
    client.ContainerStartOptions{},
  )

  return err
}
