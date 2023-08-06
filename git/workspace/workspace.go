package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

const (
	defaultAuthorName  = "bot"
	defaultAuthorEmail = "bot@example.gom"
)

// Workspace は git のレポジトリを clone したディレクトリ。clone したファイルの作業場となる
// Workspace への操作はすべて interface を通じて行わなければならない
type Workspace interface {
	CreateBranch(branch string) error
	UpdateFile(fileRelPath string, content []byte) error
	CommitFiles(branch, message string, fileRelPaths []string) error
	GitPush(branch string) error
	Clear() error
}

type workspace struct {
	repoURL string
	dir     string
	repo    *git.Repository
	auth    transport.AuthMethod
}

func NewWorkspace(ctx context.Context, repoURL string, auth transport.AuthMethod) (Workspace, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create a tempdir: %w", err)
	}
	repo, err := git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		Auth:     auth,
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to clone a repository %s: %w", repoURL, err)
	}
	return &workspace{
		repoURL: repoURL,
		dir:     dir,
		repo:    repo,
		auth:    auth,
	}, nil
}

func (w *workspace) CreateBranch(branch string) error {
	headRef, err := w.repo.Head()
	if err != nil {
		return err
	}
	ref := plumbing.NewHashReference(
		plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		headRef.Hash())
	return w.repo.Storer.SetReference(ref)
}

func (w *workspace) UpdateFile(fileRelPath string, content []byte) error {
	if err := os.WriteFile(filepath.Join(w.dir, fileRelPath), content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fileRelPath, err)
	}
	return nil
}

func (w *workspace) CommitFiles(branch, message string, fileRelPaths []string) error {
	wt, err := w.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to commit files: %w", err)
	}
	for _, fp := range fileRelPaths {
		if _, err := wt.Add(fp); err != nil {
			return fmt.Errorf("failed to commit file %s: %w", fp, err)
		}
	}
	if err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Keep:   true,
	}); err != nil {
		return err
	}
	_, err = wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  defaultAuthorName,
			Email: defaultAuthorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}

func (w *workspace) GitPush(branch string) error {
	ref := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))
	return w.repo.Push(&git.PushOptions{
		Auth:     w.auth,
		Progress: os.Stdout,
		RefSpecs: []config.RefSpec{
			config.RefSpec(ref + ":" + ref),
		},
		Force: true,
	})
}

// Clear removes all contents in repo. Make it clean for not pushing something in progress
func (w *workspace) Clear() error {
	if err := os.RemoveAll(w.dir); err != nil {
		return fmt.Errorf("failed to clear working dir: %w", err)
	}
	return nil
}
