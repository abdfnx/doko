package core

import (
	"fmt"
	"runtime"
	"context"

	"github.com/abdfnx/doko/docker"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

type info struct {
	*tview.TextView
	Docker *dockerInfo
	Host   *hostInfo
}

type dockerInfo struct {
	HostName      string
	ServerVersion string
	APIVersion    string
	KernelVersion string
	OSType        string
	Architecture  string
	Endpoint      string
	Containers    int
	Images        int
	MemTotal      string
}

type hostInfo struct {
	OSType       string
	Architecture string
}

func newHostInfo() *hostInfo {
	return &hostInfo{
		OSType:       runtime.GOOS,
		Architecture: runtime.GOARCH,
	}
}

func newDockerInfo() *dockerInfo {
	info, err := docker.Client.Info(context.TODO())
	if err != nil {
		return nil
	}

	var engineVersion string
	if v, err := docker.Client.ServerVersion(context.TODO()); err != nil {
		engineVersion = ""
	} else {
		engineVersion = v.APIVersion
	}

	return &dockerInfo{
		HostName:      info.Name,
		ServerVersion: info.ServerVersion,
		APIVersion:    engineVersion,
		KernelVersion: info.KernelVersion,
		OSType:        info.OSType,
		Architecture:  info.Architecture,
		Endpoint:      docker.Client.DaemonHost(),
		Containers:    info.Containers,
		Images:        info.Images,
		MemTotal:      fmt.Sprintf("%dMB", info.MemTotal/1024/1024),
	}
}

func newInfo() *info {
	i := &info{
		TextView: tview.NewTextView(),
		Docker:   newDockerInfo(),
		Host:     newHostInfo(),
	}

	i.display()

	return i
}

func (i *info) display() {
	dockerEngine := fmt.Sprintf("engine version: %s", i.Docker.APIVersion)
	dockerVersion := fmt.Sprintf("server version: %s", i.Docker.ServerVersion)
	dockerEndpoint := fmt.Sprintf("endpoint: %s", i.Docker.Endpoint)

	i.SetTextColor(tcell.ColorYellow)
	i.SetText(fmt.Sprintf(" doko | %s | %s | %s", dockerEngine, dockerVersion, dockerEndpoint))
}
