package github

import (
	"context"
	"github.com/google/go-github/v36/github"
	"golang.org/x/oauth2"
)

func GetGithubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ""},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client

	// list all repositories for the authenticated user
	//repos, _, err := client.Repositories.List(ctx, "", nil)

	//fmt.Print(err)
	//fmt.Print(repos)

}
