package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/harshitsinghai/gogist/models"
	"github.com/harshitsinghai/gogist/utils"
	"github.com/joho/godotenv"
	"github.com/levigross/grequests"
	"github.com/urfave/cli"
)

var githubTokenKey string
var githubAPI = "https://api.github.com/"
var requestOptions *grequests.RequestOptions

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	githubTokenKey = os.Getenv("GithubTokenKey")
	requestOptions = &grequests.RequestOptions{Headers: map[string]string{"Accept": "application/vnd.github.v3+json"}, Auth: []string{"token ", githubTokenKey}}
}

func getResp(username string) []models.Repo {
	var repos []models.Repo

	repoURL := fmt.Sprintf(githubAPI+"users/%s/repos", username)
	resp, err := grequests.Get(repoURL, requestOptions)
	if err != nil {
		fmt.Println((err))
	}
	resp.JSON(&repos)
	return repos
}

func createGist(gist models.Gist) *models.GistResponse {

	postBody, _ := json.Marshal(gist)
	requestOptionsCopy := requestOptions
	requestOptionsCopy.JSON = string(postBody)

	var gistResponse *models.GistResponse

	resp, err := grequests.Post(githubAPI+"gists", requestOptionsCopy)
	if err != nil {
		log.Println("Create request failed for Github API")
	}

	resp.JSON(&gistResponse)
	return gistResponse
}

func createGistFromFolder(description string, root string) *models.GistResponse {

	myFiles := make(map[string]models.File)
	var filesPath []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		filesPath = append(filesPath, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, path := range filesPath[1:] {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println("Error when reading the file")
		}
		fileName := strings.Split(path, "/")[1]
		myFiles[fileName] = models.File{Content: string(content)}
		fmt.Println("Creating gist for ", fileName)
	}
	fmt.Println()

	gist := models.Gist{
		Description: description,
		Files:       myFiles,
		Public:      true,
	}

	return createGist(gist)
}

func createGistFromFiles(args cli.Args) *models.GistResponse {

	description := args.Get(0)
	myFiles := make(map[string]models.File)

	for i := 1; i < args.Len(); i++ {
		content, err := ioutil.ReadFile(args.Get(i))
		if err != nil {
			log.Println("Error when reading the file")
		}
		myFiles[args.Get(i)] = models.File{Content: string(content)}
	}

	gist := models.Gist{
		Description: description,
		Files:       myFiles,
		Public:      true,
	}

	return createGist(gist)
}

func createTimeline(username string) {
	fmt.Println("Generating timeline.html....")
	repoDetails := getResp(username)

	sort.Slice(repoDetails, func(i, j int) bool {
		return repoDetails[i].CreatedAt.Before(repoDetails[j].CreatedAt)
	})

	utils.GenerateTimeline(repoDetails)
	fmt.Println("Generated timeline.html....")
}

func main() {

	app := &cli.App{
		Name:    "gogist",
		Version: "1.0",
		Commands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Usage:   "Fetch all the repo name for the given github username. [Usage]: goTool fetch user_name",
				Action: func(c *cli.Context) error {
					if c.Args().Len() > 0 {
						// Github API
						username := c.Args().Get(0)
						repoDetails := getResp(username)
						for _, repo := range repoDetails {
							fmt.Println(repo.Name)
						}
					} else {
						log.Println("Please give a username. See -h to see help")
					}
					return nil
				},
			},
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "Creates a gist of the given file(s). [Usage]: goTool 'description' sample1.txt sample2.txt",
				Action: func(c *cli.Context) error {
					if c.Args().Len() > 0 {
						// Github API Logic
						gistResponse := createGistFromFiles(c.Args())
						log.Println("Created gist of all the file(s)")

						fmt.Println("URL ", gistResponse.URL)
						fmt.Println("Description ", gistResponse.Description)
						// log.Println(resp.String())
					} else {
						log.Println("Please give sufficient arguments. See -h to see help")
					}
					return nil
				},
			},
			{
				Name:    "create-from-dir",
				Aliases: []string{"dir"},
				Usage:   "Creates a gist from the given text. [Usage]: goTool 'description' ./folder_name",
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 2 {
						// Github API Logic
						description := c.Args().Get(0)
						rootDir := c.Args().Get(1)

						gistResponse := createGistFromFolder(description, rootDir)
						fmt.Println("URL ", gistResponse.URL)
						fmt.Println("Description ", gistResponse.Description)
						// log.Println(resp.String())
						// log.Println("Done")
					} else {
						log.Println("Please give sufficient arguments. See -h to see help")
					}
					return nil
				},
			},
			{
				Name:    "create-timeline",
				Aliases: []string{"timeline"},
				Usage:   "Creates a timeline.html file based on your github repo. [Usage]: goTool create-timeline user_name",
				Action: func(c *cli.Context) error {
					if c.Args().Len() > 0 {
						// Github API
						username := c.Args().Get(0)
						createTimeline(username)
					} else {
						log.Println("Please give a username. See -h to see help")
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
