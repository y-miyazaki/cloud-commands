package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"
	"github.com/y-miyazaki/cloud-commands/pkg/command"
)

// ShowManifests is az acr show-manifests command result.
type ShowManifests struct {
	Digest    string   `json:"digest"`
	Tags      []string `json:"tags"`
	Timestamp string   `json:"timestamp"`
}

func run(args []string) error {
	app := cli.NewApp()
	app.Name = "azacrdelete"
	app.Usage = "This command remove old image in Azure Container Registry."
	app.Version = "v1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "resource-group",
			Usage: "Name of resource group.",
		},
		cli.StringFlag{
			Name:     "name",
			Usage:    "The name of the container registry.",
			Required: true,
		},
		cli.StringFlag{
			Name:  "repository",
			Usage: "The name of the repository.",
		},
		cli.StringFlag{
			Name:  "subscription",
			Usage: "Name or ID of subscription.",
		},
		cli.IntFlag{
			Name:  "keep",
			Value: 20,
			Usage: "Keep image manifest, if 20 is set, the latest 20 manifests will remain.",
		},
		cli.BoolFlag{
			Name:  "all-repository",
			Usage: "remove all repository from container registry",
		},
	}
	app.Action = func(c *cli.Context) error {
		var err error
		resourceGroup := c.String("resource-group")
		subscription := c.String("subscription")
		name := c.String("name")
		repository := c.String("repository")
		keep := c.Int("keep")
		allRepository := c.Bool("all-repository")

		// repository delete image tags.
		if allRepository {
			var repositories []string
			out, err := command.OutputStr("az acr repository list --name "+name+" --output table", false)
			if err != nil {
				return err
			}
			r := regexp.MustCompile(`^(Result|-)`)
			repositories = strings.Split(out, "\n")
			for _, repo := range repositories {
				if r.MatchString(repo) || repo == "" {
					continue
				}
				removeImageTags(resourceGroup, subscription, name, repo, keep)
			}
		} else {
			err = removeImageTags(resourceGroup, subscription, name, repository, keep)
		}
		return err
	}
	return app.Run(args)
}

func main() {
	_ = run(os.Args)
}

// removeImageTags removes image tags, can specify image tags to keep in keep option.
func removeImageTags(resourceGroup string, subscription string, name string, repository string, keep int) error {
	showManifestsCommand := []string{"acr", "repository", "show-manifests", "--orderby", "time_desc"}
	if resourceGroup != "" {
		showManifestsCommand = append(showManifestsCommand, "--resource-group")
		showManifestsCommand = append(showManifestsCommand, resourceGroup)
	}
	if subscription != "" {
		showManifestsCommand = append(showManifestsCommand, "--subscription")
		showManifestsCommand = append(showManifestsCommand, subscription)
	}
	if name != "" {
		showManifestsCommand = append(showManifestsCommand, "--name")
		showManifestsCommand = append(showManifestsCommand, name)
	}
	if repository != "" {
		showManifestsCommand = append(showManifestsCommand, "--repository")
		showManifestsCommand = append(showManifestsCommand, repository)
	}
	// get repository image tag lists.
	out, err := command.Output("az", false, showManifestsCommand...)
	if err != nil {
		return err
	}
	var showManifests []ShowManifests
	if err = json.Unmarshal([]byte(out), &showManifests); err != nil {
		return err
	}

	repositoryDeleteCommand := []string{"acr", "repository", "delete", "--yes"}
	if name != "" {
		repositoryDeleteCommand = append(repositoryDeleteCommand, "--name")
		repositoryDeleteCommand = append(repositoryDeleteCommand, name)
	}
	if repository != "" {
		repositoryDeleteCommand = append(repositoryDeleteCommand, "--image")
	}
	var tmpCommand []string
	fmt.Println("#--------------------------------------------------------------")
	fmt.Printf("# %s remove image tags(keep %d)\n", repository, keep)
	fmt.Println("#--------------------------------------------------------------")
	for i, p := range showManifests {
		if keep <= i {
			fmt.Printf("%s : %s\n", p.Digest, p.Tags)
			tmpCommand = append(repositoryDeleteCommand, repository+"@"+p.Digest)
			_, err := command.Output("az", true, tmpCommand...)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("")
	return nil
}
