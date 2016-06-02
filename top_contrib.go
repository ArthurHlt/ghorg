package main

import (
	"github.com/google/go-github/github"
	"sort"
	"time"
)

type TopContributors struct {
	client       *github.Client
	contributors Contributors
	org          string
}

var contributorsInOrg map[string]bool = make(map[string]bool)

func (this *TopContributors) Load() error {
	page := 1

	for ; ; {
		opt := &github.RepositoryListByOrgOptions{
			Type: "all",
			ListOptions: github.ListOptions{PerPage: 10, Page: page},
		}
		repos, _, err := this.client.Repositories.ListByOrg(this.org, opt)
		if err != nil {
			return err
		}
		this.analyseRepos(repos)
		if len(repos) == 0 {
			break
		}
		page ++
	}
	this.updateReposNumberForContributors()
	return nil
}
func (this *TopContributors) updateReposNumberForContributors() {
	for index, _ := range this.contributors {
		this.contributors[index].NumberRepos = len(this.contributors[index].OwnedRepos)
	}
}
func (this *TopContributors) analyseRepos(repos []github.Repository) {
	for _, repo := range repos {
		this.updateContributorFromRepp(*repo.Name)
	}
}
func (this *TopContributors) GetTopContributors() Contributors {

	return this.contributors
}
func (this *TopContributors) updateContributorFromRepp(repoName string) {
	i := 0
	var contributorStats []github.ContributorStats
	var err error
	for ; ; {
		contributorStats, _, err = this.client.Repositories.ListContributorsStats(this.org, repoName)
		if err != nil && i < 5 {
			log("Retrying for repo '" + this.org + "/" + repoName + "', details: " + err.Error())
			i++
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			log("Warning: repo '" + this.org + "/" + repoName + "' skipped: " + err.Error())
			return
		}
		break;
	}

	contributorsInRepo := make(Contributors, 0)
	for _, contributorStat := range contributorStats {
		author := *contributorStat.Author.Login
		isMember := this.isContributorInOrg(author)
		if err != nil {
			log("Warning: repo '" + this.org + "/" + repoName + "' skipped: " + err.Error())
			return
		}
		commits := *contributorStat.Total
		contributorsInRepo = append(contributorsInRepo, Contributor{author, commits, make([]string, 0), 0, ""})
		if !isMember {
			continue
		}
		contributorIndex := this.getContributorIndexFromName(author)
		if contributorIndex == -1 {
			this.contributors = append(this.contributors, Contributor{author, commits, make([]string, 0), 0, this.getUserCompany(author)})
		} else {
			this.contributors[contributorIndex].Commits += commits
		}
	}
	sort.Sort(contributorsInRepo)
	if len(contributorsInRepo) > 0 && this.isContributorExists(contributorsInRepo[0].Name) {
		contributorIndex := this.getContributorIndexFromName(contributorsInRepo[0].Name)
		this.contributors[contributorIndex].OwnedRepos = append(this.contributors[contributorIndex].OwnedRepos, repoName)
	}
	log("Done for repo '" + this.org + "/" + repoName + "'.")
}
func (this *TopContributors) isContributorExists(name string) bool {
	for _, contributor := range this.contributors {
		if contributor.Name == name {
			return true
		}
	}
	return false
}
func (this *TopContributors) isContributorInOrg(name string) bool {
	if _, ok := contributorsInOrg[name]; ok {
		return contributorsInOrg[name]
	}
	isMember, _, err := this.client.Organizations.IsMember(this.org, name)
	if err != nil {
		return false
	}
	contributorsInOrg[name] = isMember
	return contributorsInOrg[name]
}
func (this *TopContributors) getUserCompany(login string) string {
	user, _, err := this.client.Users.Get(login)
	if err != nil {
		log("Warning: can't get company for user '" + login + "', details: " + err.Error())
		return "independent"
	}
	company := user.Company
	if company == nil {
		return "independent"
	}
	return *company
}
func (this *TopContributors) getContributorIndexFromName(name string) int {
	for index, contributor := range this.contributors {
		if contributor.Name == name {
			return index
		}
	}
	return -1
}

