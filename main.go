package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"fmt"
	"os"
	"strconv"
	"sort"
	"github.com/urfave/cli"
	"github.com/olekukonko/tablewriter"
	"strings"
	"github.com/gocarina/gocsv"
)

const TOKEN_ENV_VAR_KEY = "GH_TOKEN"

var ghToken string
var verbose bool
var reposDetails bool
var csvFile string
var noMarkdown bool

func main() {
	app := cli.NewApp()
	app.Name = "ghorg"
	app.Version = "1.0.0"
	app.Usage = "Analyse a github organization to find top contributors by number of commits and number of repos"
	app.Commands = []cli.Command{
		{
			Name:    "analyse",
			Aliases: []string{"a"},
			Usage:   "Analyse a github org",
			ArgsUsage: "name-of-the-org-to-scan",
			Action:  Analyse,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "repos-details, rd",
					Destination: &reposDetails,
					Usage: "To show in table all repos names owned by user",
				},
				cli.StringFlag{
					Name: "csv",
					Value: "",
					Destination: &csvFile,
					Usage: "file name to ouput in csv format",
				},
				cli.BoolFlag{
					Name: "no-markdown",
					Destination: &noMarkdown,
					Usage: "To hide markdown result",
				},
			},
		},
	}
	app.ErrWriter = os.Stderr
	app.HideHelp = false
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "ghtoken, g",
			Value: "",
			Destination: &ghToken,
			Usage: "(Mandatory) Your github token",
			EnvVar: TOKEN_ENV_VAR_KEY,
		},
		cli.BoolFlag{
			Name: "verbose, vvv",
			Destination: &verbose,
			Usage: "To use it in verbose mode",
		},
	}

	app.Run(os.Args)
}
func createCsvFile(contributors Contributors) {
	if csvFile == "" {
		return
	}
	f, err := os.Create(csvFile)
	fatalIf(err)
	defer f.Close()
	csvWriter := gocsv.DefaultCSVWriter(f)
	csvWriter.Comma = ';'
	gocsv.MarshalCSV(contributors, csvWriter)
	fatalIf(err)
}
func Analyse(c *cli.Context) error {
	args := c.Args()
	if len(args) < 1 {
		fatal("You need to pass an organization name, see usage with help command.")
	}
	org := args.First()
	topContributors := &TopContributors{NewGithubClient(), make(Contributors, 0), org}
	err := topContributors.Load()
	fatalIf(err)
	contributors := topContributors.GetTopContributors()
	sort.Sort(contributors)
	createCsvFile(contributors)
	if noMarkdown {
		return nil
	}
	fmt.Println("# Top contributors by commits")
	fmt.Println("-----------------------------")
	fmt.Println("")
	createTableCommits(contributors)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("# Top contributors by number of repositories")
	fmt.Println("--------------------------------------------")
	fmt.Println("")
	contributorsByRepos := ContributorsByRepos(contributors)
	sort.Sort(contributorsByRepos)
	createTableNumberRepos(contributorsByRepos)
	return nil
}
func createTableCommits(contributors Contributors) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Rank", "Name", "Company", "Commits"})
	for index, contributor := range contributors {
		rank := index + 1
		table.Append([]string{strconv.Itoa(rank) + "#", contributor.Name, contributor.Company, strconv.Itoa(contributor.Commits)})
	}
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.Render()
}
func createTableNumberRepos(contributors ContributorsByRepos) {
	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"Rank", "Name", "Company", "Number of repos"}
	if reposDetails {
		headers = append(headers, "Repos names")
	}
	table.SetHeader(headers)
	for index, contributor := range contributors {
		rank := index + 1
		data := []string{strconv.Itoa(rank) + "#", contributor.Name, contributor.Company, strconv.Itoa(contributor.NumberRepos)}
		if reposDetails {
			data = append(data, strings.Join(contributor.OwnedRepos, "\n"))
		}
		table.Append(data)
	}
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.Render()
}
func NewGithubClient() *github.Client {
	if ghToken == "" {
		fatal("You need to set the env var GH_TOKEN with a github token.")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}

func fatalIf(err error) {
	if err != nil {
		fatal(err.Error())
	}
}
func fatal(message string) {
	fmt.Fprintln(os.Stdout, message)
	os.Exit(1)
}
func log(message string) {
	if !verbose {
		return
	}
	fmt.Fprintln(os.Stderr, message)
}