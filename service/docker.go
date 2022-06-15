package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	json2 "encoding/json"
	"fmt"
	"syscall"

	"github.com/IceWhaleTech/CasaOS/model/notify"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/docker"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/env_helper"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"

	//"github.com/containerd/containerd/oci"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	client2 "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerService interface {
	DockerPullImage(imageName string, icon, name string) error
	IsExistImage(imageName string) bool
	DockerContainerCreate(imageName string, m model.CustomizationPostData) (containerId string, err error)
	DockerContainerCopyCreate(info *types.ContainerJSON) (containerId string, err error)
	DockerContainerStart(name string) error
	DockerContainerStats(name string) (string, error)
	DockerListByName(name string) (*types.Container, error)
	DockerListByImage(image, version string) (*types.Container, error)
	DockerContainerInfo(name string) (*types.ContainerJSON, error)
	DockerImageRemove(name string) error
	DockerContainerRemove(name string, update bool) error
	DockerContainerStop(id string) error
	DockerContainerUpdateName(name, id string) (err error)
	DockerContainerUpdate(m model.CustomizationPostData, id string) (err error)
	DockerContainerLog(name string) (string, error)
	DockerContainerCommit(name string)
	DockerContainerList() []types.Container
	DockerNetworkModelList() []types.NetworkResource
	DockerImageInfo(image string) (types.ImageInspect, error)
	GetNetWorkNameByNetWorkID(id string) (string, error)
	ContainerExecShell(container_id string) string
}

type dockerService struct {
	rootDir string
}

func (ds *dockerService) DockerContainerList() []types.Container {
	cli, err := client2.NewClientWithOpts(client2.FromEnv, client2.WithTimeout(time.Second*5))
	if err != nil {
		return nil
	}
	defer cli.Close()
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return containers
	}
	return containers
}

func (ds *dockerService) ContainerExecShell(container_id string) string {
	cli, _ := client2.NewClientWithOpts(client2.FromEnv)
	exec, err := cli.ContainerExecCreate(context.Background(), container_id, types.ExecConfig{
		User: "1000:1000",
		Cmd:  []string{"echo -e \"hellow\nworld\" >> /a.txt"},
	})
	if err != nil {
		os.Exit(5)
	}
	err = cli.ContainerExecStart(context.Background(), exec.ID, types.ExecStartCheck{})
	if err != nil {
		fmt.Println("exec script error ", err)
	}
	return exec.ID
}

//创建默认网络
func DockerNetwork() {

	cli, _ := client2.NewClientWithOpts(client2.FromEnv)
	defer cli.Close()
	d, _ := cli.NetworkList(context.Background(), types.NetworkListOptions{})

	for _, resource := range d {
		if resource.Name == docker.NETWORKNAME {
			return
		}
	}
	cli.NetworkCreate(context.Background(), docker.NETWORKNAME, types.NetworkCreate{})
}

//根据网络id获取网络名
func (ds *dockerService) GetNetWorkNameByNetWorkID(id string) (string, error) {
	cli, _ := client2.NewClientWithOpts(client2.FromEnv)
	defer cli.Close()
	filter := filters.NewArgs()
	filter.Add("id", id)
	d, err := cli.NetworkList(context.Background(), types.NetworkListOptions{Filters: filter})
	if err == nil && len(d) > 0 {
		return d[0].Name, nil
	}
	return "", err
}

//拉取镜像
func DockerPull() {

	cli, _ := client2.NewClientWithOpts(client2.FromEnv)
	defer cli.Close()

	authConfig := types.AuthConfig{}
	encodedJSON, err := json2.Marshal(authConfig)
	fmt.Println(err)

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	reader, err := cli.ImagePull(context.Background(), "swr.cn-north-4.myhuaweicloud.com/root/swr-demo-2048:latest", types.ImagePullOptions{RegistryAuth: authStr})

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	fmt.Println(buf.String())

}

//拉取镜像
func DockerEx() {

	cli, _ := client2.NewClientWithOpts(client2.FromEnv)
	defer cli.Close()

	importResponse, err := cli.ImageImport(context.Background(), types.ImageImportSource{
		Source:     strings.NewReader("source"),
		SourceName: "image_source",
	}, "repository_name:imported", types.ImageImportOptions{
		Tag:     "imported",
		Message: "A message",
		Changes: []string{"change1", "change2"},
	})

	response, err := ioutil.ReadAll(importResponse)
	if err != nil {
		fmt.Println(err)
	}
	importResponse.Close()
	println(string(response))
	if string(response) != "response" {
		fmt.Printf("expected response to contain 'response', got %s", string(response))
	}
}

//func DockerContainerSize() {
//	cli, err := client2.NewClientWithOpts(client2.FromEnv)
//	//but := bytes.Buffer{}
//	d, err := cli.ContainerExecCreate(context.Background(), "c3adcef92bae648890941ac00e6c4024d7f2959c2e629f0b581d6a19d77b5eda")
//	fmt.Println(d)
//	st, _ := ioutil.ReadAll(d.Body)
//	fmt.Println(string(st))
//	if err != nil {
//		fmt.Print(err)
//	}
//
//}

func (ds *dockerService) DockerImageInfo(image string) (types.ImageInspect, error) {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return types.ImageInspect{}, err
	}
	inspect, _, err := cli.ImageInspectWithRaw(context.Background(), image)
	if err != nil {
		return inspect, err
	}
	return inspect, nil
}

func MsqlExec(container string) error {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	ctx := context.Background()
	// 执行/bin/bash命令
	ir, err := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"date"},
		Tty:          true,
		Env:          []string{"aaa=ddd"},
	})
	err = cli.ContainerExecStart(ctx, ir.ID, types.ExecStartCheck{})

	fmt.Println(err)

	return err
}

func Exec(container, row, col string) (hr types.HijackedResponse, err error) {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	ctx := context.Background()
	// 执行/bin/bash命令
	ir, err := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Env:          []string{"COLUMNS=" + col, "LINES=" + row},
		Cmd:          []string{"/bin/bash"},
		Tty:          true,
	})
	if err != nil {
		return
	}
	// 附加到上面创建的/bin/bash进程中
	hr, err = cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return
	}
	return
}

func DockerLog() {
	//cli, err := client2.NewClientWithOpts(client2.FromEnv)
	//ctx := context.Background()
	//ir, err := cli.ContainerLogs(ctx, "79c6fa382c330b9149e2d28d24f4d2c231cdb8cfc0710c2d268ccee13c5b24f8", types.ContainerLogsOptions{})
	//str, err := ioutil.ReadAll(ir)
	//fmt.Println(string(str))
	//fmt.Println(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, _ := client2.NewClientWithOpts(client2.FromEnv)
	reader, err := client.ContainerLogs(ctx, "79c6fa382c330b9149e2d28d24f4d2c231cdb8cfc0710c2d268ccee13c5b24f8", types.ContainerLogsOptions{})
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}

func DockerLogs() {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	i, err := cli.ContainerLogs(context.Background(), "79c6fa382c330b9149e2d28d24f4d2c231cdb8cfc0710c2d268ccee13c5b24f8", types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true,
		Tail:       "40",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer i.Close()

	hdr := make([]byte, 8)
	for {
		_, err := i.Read(hdr)
		if err != nil {
			log.Fatal(err)
		}
		var w io.Writer
		switch hdr[0] {
		case 1:
			w = os.Stdout
		default:
			w = os.Stderr
		}
		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		_, err = i.Read(dat)
		fmt.Fprint(w, string(dat))
	}
}

//正式内容

//检查镜像是否存在
func (ds *dockerService) IsExistImage(imageName string) bool {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return false
	}
	defer cli.Close()
	filter := filters.NewArgs()
	filter.Add("reference", imageName)

	list, err := cli.ImageList(context.Background(), types.ImageListOptions{Filters: filter})

	if err == nil && len(list) > 0 {
		return true
	}

	return false
}

//安装镜像
func (ds *dockerService) DockerPullImage(imageName string, icon, name string) error {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()
	out, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	if err != nil {

		return err
	}
	//io.Copy()
	buf := make([]byte, 2048*4)
	for {
		n, err := out.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		if len(icon) > 0 && len(name) > 0 {
			notify := notify.Application{}
			notify.Icon = icon
			notify.Name = name
			notify.State = "PULLING"
			notify.Type = "INSTALL"
			notify.Finished = false
			notify.Success = true
			notify.Message = string(buf[:n])
			MyService.Notify().SendInstallAppBySocket(notify)
		}

	}
	return err
}
func (ds *dockerService) DockerContainerCopyCreate(info *types.ContainerJSON) (containerId string, err error) {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return "", err
	}
	defer cli.Close()
	container, err := cli.ContainerCreate(context.Background(), info.Config, info.HostConfig, &network.NetworkingConfig{info.NetworkSettings.Networks}, nil, info.Name)
	if err != nil {
		return "", err
	}
	return container.ID, err
}

//param imageName 镜像名称
//param containerDbId 数据库的id
//param port 容器内部主端口
//param mapPort 容器主端口映射到外部的端口
//param tcp 容器其他tcp端口
//param udp 容器其他udp端口
func (ds *dockerService) DockerContainerCreate(imageName string, m model.CustomizationPostData) (containerId string, err error) {
	if len(m.NetworkModel) == 0 {
		m.NetworkModel = "bridge"
	}

	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return "", err
	}
	defer cli.Close()
	ports := make(nat.PortSet)
	portMaps := make(nat.PortMap)
	//	ports[nat.Port(fmt.Sprint(m.PortMap)+"/tcp")] = struct{}{}
	//	if net != "host" {
	//		portMaps[nat.Port(fmt.Sprint(m.Port)+"/tcp")] = []nat.PortBinding{{HostIP: "", HostPort: m.PortMap}}
	//	}
	//port := ""
	for _, portMap := range m.Ports {
		// if portMap.CommendPort == m.PortMap && portMap.Protocol == "tcp" || portMap.Protocol == "both" {
		// 	port = portMap.ContainerPort
		// }
		if portMap.Protocol == "tcp" {

			tContainer, _ := strconv.Atoi(portMap.ContainerPort)
			if tContainer > 0 {
				ports[nat.Port(portMap.ContainerPort+"/tcp")] = struct{}{}
				if m.NetworkModel != "host" {
					portMaps[nat.Port(portMap.ContainerPort+"/tcp")] = []nat.PortBinding{{HostPort: portMap.CommendPort}}
				}
			}
		} else if portMap.Protocol == "both" {

			tContainer, _ := strconv.Atoi(portMap.ContainerPort)
			if tContainer > 0 {
				ports[nat.Port(portMap.ContainerPort+"/tcp")] = struct{}{}
				if m.NetworkModel != "host" {
					portMaps[nat.Port(portMap.ContainerPort+"/tcp")] = []nat.PortBinding{{HostPort: portMap.CommendPort}}
				}
			}

			uContainer, _ := strconv.Atoi(portMap.ContainerPort)
			if uContainer > 0 {
				ports[nat.Port(portMap.ContainerPort+"/udp")] = struct{}{}
				if m.NetworkModel != "host" {
					portMaps[nat.Port(portMap.ContainerPort+"/udp")] = []nat.PortBinding{{HostPort: portMap.CommendPort}}
				}
			}

		} else {
			uContainer, _ := strconv.Atoi(portMap.ContainerPort)
			if uContainer > 0 {
				ports[nat.Port(portMap.ContainerPort+"/udp")] = struct{}{}
				if m.NetworkModel != "host" {
					portMaps[nat.Port(portMap.ContainerPort+"/udp")] = []nat.PortBinding{{HostPort: portMap.CommendPort}}
				}
			}
		}

	}

	var envArr []string
	var showENV []string
	showENV = append(showENV, "casaos")
	for _, e := range m.Envs {
		showENV = append(showENV, e.Name)
		if strings.HasPrefix(e.Value, "$") {
			envArr = append(envArr, e.Name+"="+env_helper.ReplaceDefaultENV(e.Value, MyService.System().GetTimeZone()))
			continue
		}
		if len(e.Value) > 0 {
			if e.Value == "port_map" {
				envArr = append(envArr, e.Name+"="+m.PortMap)
				continue
			}
			envArr = append(envArr, e.Name+"="+e.Value)
		}
	}

	res := container.Resources{}
	if m.CpuShares > 0 {
		res.CPUShares = m.CpuShares
	}
	if m.Memory > 0 {
		res.Memory = m.Memory << 20
	}
	for _, p := range m.Devices {
		if len(p.Path) > 0 {
			res.Devices = append(res.Devices, container.DeviceMapping{PathOnHost: p.Path, PathInContainer: p.ContainerPath, CgroupPermissions: "rwm"})
		}
	}
	hostConfingBind := []string{}
	// volumes bind
	volumes := []mount.Mount{}
	for _, v := range m.Volumes {
		path := v.Path
		if len(path) == 0 {
			path = docker.GetDir(m.Label, v.Path)
			if len(path) == 0 {
				continue
			}
		}
		path = strings.ReplaceAll(path, "$AppID", m.Label)
		//reg1 := regexp.MustCompile(`([^<>/\\\|:""\*\?]+\.\w+$)`)
		//result1 := reg1.FindAllStringSubmatch(path, -1)
		//if len(result1) == 0 {
		err = file.IsNotExistMkDir(path)
		if err != nil {
			loger.Error("Failed to create a folder", zap.Any("err", err))
			continue
		}
		//}
		//  else {
		// 	err = file.IsNotExistCreateFile(path)
		// 	if err != nil {
		// 		ds.log.Error("mkdir error", err)
		// 		continue
		// 	}
		// }

		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: path,
			Target: v.ContainerPath,
		})

		hostConfingBind = append(hostConfingBind, v.Path+":"+v.ContainerPath)
	}

	rp := container.RestartPolicy{}

	if len(m.Restart) > 0 {
		rp.Name = m.Restart
	}
	// healthTest := []string{}
	// if len(port) > 0 {
	// 	healthTest = []string{"CMD-SHELL", "curl -f http://localhost:" + port + m.Index + " || exit 1"}
	// }

	// health := &container.HealthConfig{
	// 	Test:        healthTest,
	// 	StartPeriod: 0,
	// 	Retries:     1000,
	// }
	// fmt.Print(health)
	if len(m.HostName) == 0 {
		m.HostName = m.Label
	}
	config := &container.Config{
		Image:  imageName,
		Labels: map[string]string{"origin": m.Origin, m.Origin: m.Origin, "casaos": "casaos"},
		Env:    envArr,
		//	Healthcheck: health,
		Hostname: m.HostName,
		Cmd:      m.Cmd,
	}

	config.Labels["web"] = m.PortMap
	config.Labels["icon"] = m.Icon
	config.Labels["desc"] = m.Description
	config.Labels["index"] = m.Index
	config.Labels["custom_id"] = m.CustomId
	config.Labels["show_env"] = strings.Join(showENV, ",")
	config.Labels["protocol"] = m.Protocol
	config.Labels["host"] = m.Host
	//config.Labels["order"] = strconv.Itoa(MyService.App().GetCasaOSCount() + 1)
	hostConfig := &container.HostConfig{Resources: res, Mounts: volumes, RestartPolicy: rp, NetworkMode: container.NetworkMode(m.NetworkModel), Privileged: m.Privileged, CapAdd: m.CapAdd}
	//if net != "host" {
	config.ExposedPorts = ports
	hostConfig.PortBindings = portMaps
	//}

	containerDb, err := cli.ContainerCreate(context.Background(),
		config,
		hostConfig,
		&network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{m.NetworkModel: {NetworkID: "", Aliases: []string{}}}},
		nil,
		m.Label)
	if err != nil {
		return "", err
	}
	return containerDb.ID, err
}

//删除容器
func (ds *dockerService) DockerContainerRemove(name string, update bool) error {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()
	err = cli.ContainerRemove(context.Background(), name, types.ContainerRemoveOptions{})

	//路径处理
	if !update {
		path := docker.GetDir(name, "/config")
		if !file.CheckNotExist(path) {
			file.RMDir(path)
		}
	}

	if err != nil {
		return err
	}

	return err
}

//删除镜像
func (ds *dockerService) DockerImageRemove(name string) error {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()
	imageList, err := cli.ImageList(context.Background(), types.ImageListOptions{})

	imageId := ""

Loop:
	for _, ig := range imageList {
		for _, i := range ig.RepoTags {
			if i == name {
				imageId = ig.ID
				break Loop
			}
		}
	}
	_, err = cli.ImageRemove(context.Background(), imageId, types.ImageRemoveOptions{})
	return err
}

func DockerImageRemove(name string) error {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()
	imageList, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	imageId := ""

Loop:
	for _, ig := range imageList {
		fmt.Println(ig.RepoDigests)
		fmt.Println(ig.Containers)
		for _, i := range ig.RepoTags {
			if i == name {
				imageId = ig.ID
				break Loop
			}
		}
	}
	_, err = cli.ImageRemove(context.Background(), imageId, types.ImageRemoveOptions{})
	return err
}

//停止镜像
func (ds *dockerService) DockerContainerStop(id string) error {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()
	err = cli.ContainerStop(context.Background(), id, nil)
	return err
}

//启动容器
func (ds *dockerService) DockerContainerStart(name string) error {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()
	err = cli.ContainerStart(context.Background(), name, types.ContainerStartOptions{})
	return err
}

//查看日志
func (ds *dockerService) DockerContainerLog(name string) (string, error) {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return "", err
	}
	defer cli.Close()
	body, err := cli.ContainerLogs(context.Background(), name, types.ContainerLogsOptions{ShowStderr: true, ShowStdout: true})
	if err != nil {
		return "", err
	}

	defer body.Close()
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func DockerContainerStats1() error {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()
	dss, err := cli.ContainerStats(context.Background(), "dockermysql", false)
	if err != nil {
		return err
	}
	defer dss.Body.Close()
	sts, err := ioutil.ReadAll(dss.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(sts))
	return nil
}

//获取容器状态
func (ds *dockerService) DockerContainerStats(name string) (string, error) {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return "", err
	}
	defer cli.Close()
	dss, err := cli.ContainerStats(context.Background(), name, false)
	if err != nil {
		return "", err
	}
	defer dss.Body.Close()
	sts, err := ioutil.ReadAll(dss.Body)
	if err != nil {
		return "", err
	}
	return string(sts), nil
}

//备份容器
func (ds *dockerService) DockerContainerCommit(name string) {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		fmt.Println(err)
	}
	defer cli.Close()
	d, err := cli.ContainerInspect(context.Background(), name)
	dss, err := cli.ContainerCommit(context.Background(), name, types.ContainerCommitOptions{Reference: "test", Config: d.Config})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dss)
}

func (ds *dockerService) DockerListByName(name string) (*types.Container, error) {
	cli, _ := client2.NewClientWithOpts(client2.FromEnv)
	defer cli.Close()
	filter := filters.NewArgs()
	filter.Add("name", name)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{Filters: filter})
	if err != nil {
		return &types.Container{}, err
	}
	if len(containers) == 0 {
		return &types.Container{}, errors.New("not found")
	}
	return &containers[0], nil
}

func (ds *dockerService) DockerListByImage(image, version string) (*types.Container, error) {
	cli, _ := client2.NewClientWithOpts(client2.FromEnv)
	defer cli.Close()
	filter := filters.NewArgs()
	filter.Add("ancestor", image+":"+version)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{Filters: filter})
	if err != nil {
		return nil, err
	}
	if len(containers) == 0 {
		return nil, nil
	}
	return &containers[0], nil
}

//获取容器详情
func (ds *dockerService) DockerContainerInfo(name string) (*types.ContainerJSON, error) {

	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return &types.ContainerJSON{}, err
	}
	defer cli.Close()
	d, err := cli.ContainerInspect(context.Background(), name)
	if err != nil {
		return &types.ContainerJSON{}, err
	}
	return &d, nil
}

//更新容器
//param shares cpu优先级
//param containerDbId 数据库的id
//param port 容器内部主端口
//param mapPort 容器主端口映射到外部的端口
//param tcp 容器其他tcp端口
//param udp 容器其他udp端口
func (ds *dockerService) DockerContainerUpdate(m model.CustomizationPostData, id string) (err error) {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()
	//重启策略
	rp := container.RestartPolicy{
		Name:              "",
		MaximumRetryCount: 0,
	}
	if len(m.Restart) > 0 {
		rp.Name = m.Restart
	}
	res := container.Resources{}

	if m.Memory > 0 {
		res.Memory = m.Memory * 1024 * 1024
		res.MemorySwap = -1
	}
	if m.CpuShares > 0 {
		res.CPUShares = m.CpuShares
	}
	for _, p := range m.Devices {
		res.Devices = append(res.Devices, container.DeviceMapping{PathOnHost: p.Path, PathInContainer: p.ContainerPath, CgroupPermissions: "rwm"})
	}
	_, err = cli.ContainerUpdate(context.Background(), id, container.UpdateConfig{RestartPolicy: rp, Resources: res})
	if err != nil {
		return err
	}

	return
}

//更新容器名称
//param name 容器名称
//param id 老的容器名称
func (ds *dockerService) DockerContainerUpdateName(name, id string) (err error) {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerRename(context.Background(), id, name)
	if err != nil {
		return err
	}
	return
}

//获取网络列表
func (ds *dockerService) DockerNetworkModelList() []types.NetworkResource {

	cli, _ := client2.NewClientWithOpts(client2.FromEnv)
	defer cli.Close()
	networks, _ := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	return networks
}
func NewDockerService() DockerService {
	return &dockerService{rootDir: command2.ExecResultStr(`source ./shell/helper.sh ;GetDockerRootDir`)}
}

//   ---------------------------------------test------------------------------------
//func ServiceCreate() {
//	cli, err := client2.NewClientWithOpts(client2.FromEnv)
//	r, err := cli.ServiceCreate(context.Background(), swarm.ServiceSpec{}, types.ServiceCreateOptions{})
//	if err != nil {
//		fmt.Println("error", err)
//	}
//
//
//}

func Containerd() {
	// create a new client connected to the default socket path for containerd
	cli, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		fmt.Println("111")
		fmt.Println(err)
	}
	defer cli.Close()

	// create a new context with an "example" namespace
	ctx := namespaces.WithNamespace(context.Background(), "default")

	// pull the redis image from DockerHub
	image, err := cli.Pull(ctx, "docker.io/library/busybox:latest", containerd.WithPullUnpack)
	if err != nil {
		fmt.Println("222")
		fmt.Println(err)
	}

	// create a container
	container, err := cli.NewContainer(
		ctx,
		"test1",
		containerd.WithImage(image),
		containerd.WithNewSnapshot("redis-server-snapshot1", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)

	if err != nil {
		fmt.Println(err)
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	// create a task from the container
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		fmt.Println(err)
	}
	defer task.Delete(ctx)

	// make sure we wait before calling start
	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// call start on the task to execute the redis server
	if err = task.Start(ctx); err != nil {
		fmt.Println(err)
	}

	fmt.Println("执行完成等待")
	// sleep for a lil bit to see the logs
	time.Sleep(3 * time.Second)

	// kill the process and get the exit status
	if err = task.Kill(ctx, syscall.SIGTERM); err != nil {
		fmt.Println(err)
	}

	// wait for the process to fully exit and print out the exit status

	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("redis-server exited with status: %d\n", code)

}
