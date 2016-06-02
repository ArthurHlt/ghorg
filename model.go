package main

type Contributor struct {
	Name        string     `csv:"name"`
	Commits     int        `csv:"commits"`
	OwnedRepos  []string   `csv:"-"`
	NumberRepos int        `csv:"repos"`
	Company     string     `csv:"company"`
}

type Contributors []Contributor

func (c Contributors) Len() int {
	return len(c)
}
func (c Contributors) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c Contributors) Less(i, j int) bool {
	return c[i].Commits > c[j].Commits
}

type ContributorsByRepos []Contributor

func (c ContributorsByRepos) Len() int {
	return len(c)
}
func (c ContributorsByRepos) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c ContributorsByRepos) Less(i, j int) bool {
	return c[i].NumberRepos > c[j].NumberRepos
}