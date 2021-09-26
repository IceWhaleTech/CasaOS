package docker_base

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	client2 "github.com/docker/docker/client"
	"time"
)

//创建一个mysql数据库
func MysqlCreate(mysqlConfig MysqlConfig, dbId string, cpuShares int64, memory int64) (string, error) {
	const imageName = "mysql"
	const imageVersion = "8"
	const imageNet = "oasis"

	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return "", err
	}
	defer cli.Close()
	_, err = cli.ImagePull(context.Background(), imageName+":"+imageVersion, types.ImagePullOptions{})

	isExist := true
	//检查到镜像才继续
	for isExist {
		filter := filters.NewArgs()
		filter.Add("before", imageName+":"+imageVersion)
		list, e := cli.ImageList(context.Background(), types.ImageListOptions{Filters: filter})
		if e == nil && len(list) > 0 {
			isExist = false
		}
		time.Sleep(time.Second)
	}

	var envArr = []string{"MYSQL_ROOT_PASSWORD=" + mysqlConfig.DataBasePassword, "MYSQL_DATABASE=" + mysqlConfig.DataBaseDB}

	res := container.Resources{}
	if cpuShares > 0 {
		res.CPUShares = cpuShares
	}
	if memory > 0 {
		res.Memory = memory << 20
	}

	rp := container.RestartPolicy{}

	rp.Name = "always"

	config := &container.Config{
		Image:  imageName,
		Labels: map[string]string{"version": imageVersion, "author": "official"},
		Env:    envArr,
	}
	hostConfig := &container.HostConfig{Resources: res, RestartPolicy: rp, NetworkMode: container.NetworkMode(imageNet)}

	containerCreate, err := cli.ContainerCreate(context.Background(),
		config,
		hostConfig,
		&network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{imageNet: {NetworkID: ""}}},
		nil,
		dbId)

	containerId := containerCreate.ID

	//启动容器
	err = cli.ContainerStart(context.Background(), dbId, types.ContainerStartOptions{})
	if err != nil {
		return containerId, err
	}

	return containerId, nil

}

func MysqlDelete(dbId string) error {

	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()
	err = cli.ContainerStop(context.Background(), dbId, nil)
	if err != nil {
		return err
	}

	err = cli.ContainerRemove(context.Background(), dbId, types.ContainerRemoveOptions{})
	return err

}
