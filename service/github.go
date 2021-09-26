package service

import (
	"context"
	"fmt"
	"github.com/google/go-github/v36/github"
	"github.com/tidwall/gjson"
)

type GithubService interface {
	GetManifestJsonByRepo() (image, tcp, udp string)
}

type githubService struct {
	cl *github.Client
}

func (g *githubService) GetManifestJsonByRepo() (image, tcp, udp string) {
	c, _, _, e := g.cl.Repositories.GetContents(context.Background(), "a624669980", "o_test_json", "/OasisManifest.json", &github.RepositoryContentGetOptions{})
	if e != nil {
		fmt.Println(e)
	}
	str, e := c.GetContent()
	if e != nil {
		fmt.Println(e)
	}
	image = gjson.Get(str, "dockerImage").String()
	tcp = gjson.Get(str, "tcp_ports").Raw
	udp = gjson.Get(str, "udp_ports").Raw
	return
}

func GetNewGithubService(cl *github.Client) GithubService {
	return &githubService{cl: cl}
}
