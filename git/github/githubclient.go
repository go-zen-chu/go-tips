package github

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client is a interface that handle about github
type Client interface {
	Clone(repoURI string, dir string) (*git.Repository, error)
	Commit(r *git.Repository, msg string) error
	Push(r *git.Repository) error
	PullRequest(owner, repo, title, head, body, baseBranch string) (string, error)
	ListRepoIssuesSince(owner, repo string, since time.Time, state string, labels []string) ([]*github.Issue, error)
	ListRepoIssues(owner, repo string, state string, labels []string) ([]*github.Issue, error)
}

type ghclient struct {
	client *github.Client
	ctx    context.Context
	user   string
	mail   string
	token  string
}

// NewGitHubClient create GitHubClient implementation
func NewGitHubClient(baseURL string, token string, user string, mail string) (Client, error) {
	// validation
	if baseURL == "" {
		return nil, errors.New("need to set baseURL")
	}
	uploadURL := path.Join(baseURL, "upload")
	if token == "" {
		return nil, errors.New("need to set token")
	}
	if user == "" || mail == "" {
		return nil, errors.New("need to set user, mail for git operation")
	}

	// initialize github client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	cli, err := github.NewClient(baseURL, uploadURL, tc)
	if err != nil {
		return nil, fmt.Errorf("creating github enterprise client: %s", err)
	}
	c := &ghclient{
		client: cli,
		ctx:    ctx,
		user:   user,
		mail:   mail,
		token:  token,
	}
	return c, nil
}

// Clone is function of 'git clone'
func (c *ghclient) Clone(repoURI string, dir string) (*git.Repository, error) {
	o := &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: c.user,
			Password: c.token,
		},
		URL: repoURI,
	}
	return git.PlainClone(dir, false, o)
}

// Commit is function of 'git commit'
func (c *ghclient) Commit(repo *git.Repository, msg string) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	o := &git.CommitOptions{
		Author: &object.Signature{
			Name:  c.user,
			Email: c.mail,
			When:  time.Now(),
		},
	}
	_, err = w.Commit(msg, o)
	return err
}

// Push is function of 'git push'
func (c *ghclient) Push(repo *git.Repository) error {
	o := &git.PushOptions{
		Auth: &http.BasicAuth{
			Username: c.user,
			Password: c.token,
		},
		// TODO: add refspec
	}
	return repo.Push(o)
}

// PullRequest is function which create new pull request
func (c *ghclient) PullRequest(owner, repo, title, head, body, baseBranch string) (string, error) {
	npr := github.NewPullRequest{
		Title: &title,
		Head:  &head,
		Base:  &baseBranch,
		Body:  &body,
	}
	pr, _, err := c.client.PullRequests.Create(c.ctx, owner, repo, &npr)
	if err != nil {
		return "", fmt.Errorf("creating PullRequest : %s", err)
	}
	prURL := pr.GetHTMLURL()
	log.Printf("PullRequest created: %s", prURL)
	return prURL, nil
}

// private function for handling pagination
func listRepoIssues(listFunc func(pageIdx int) ([]*github.Issue, *github.Response, error)) ([]*github.Issue, error) {
	// github pagination :https://developer.github.com/v3/guides/traversing-with-pagination/
	maxTry := 20 // limit requests for safety
	pageIdx := 1
	issues := make([]*github.Issue, 0)
	for ; maxTry > 0; maxTry-- {
		iss, resp, err := listFunc(pageIdx)
		if err != nil {
			return nil, fmt.Errorf("list issues from repo: %s, pageIdx %d, lastPageIdx %d", err, pageIdx, resp.LastPage)
		}
		issues = append(issues, iss...)
		// last page index is 0 when no more pagination
		if resp.LastPage == 0 {
			break
		}
		pageIdx = resp.NextPage
	}
	if maxTry == 0 {
		return issues, fmt.Errorf("list issues reached to max try: %d", maxTry)
	}
	return issues, nil
}

// ListRepoIssuesSince lists issues since
func (c *ghclient) ListRepoIssuesSince(owner, repo string, since time.Time, state string, labels []string) ([]*github.Issue, error) {
	return listRepoIssues(func(pageIdx int) ([]*github.Issue, *github.Response, error) {
		return c.client.Issues.ListByRepo(c.ctx, owner, repo, &github.IssueListByRepoOptions{
			State:  state,
			Labels: labels,
			Since:  since, // since get updated issues: https://developer.github.com/v3/issues/#list-repository-issues
			ListOptions: github.ListOptions{
				Page:    pageIdx,
				PerPage: 30,
			},
		})
	})
}

// ListRepoIssues lists issues
func (c *ghclient) ListRepoIssues(owner, repo string, state string, labels []string) ([]*github.Issue, error) {
	return listRepoIssues(func(pageIdx int) ([]*github.Issue, *github.Response, error) {
		return c.client.Issues.ListByRepo(c.ctx, owner, repo, &github.IssueListByRepoOptions{
			State:  state,
			Labels: labels,
			ListOptions: github.ListOptions{
				Page:    pageIdx,
				PerPage: 30,
			},
		})
	})
}
