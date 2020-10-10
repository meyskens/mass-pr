package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5/config"

	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/mb0/glob"
	gocopy "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(NewAddDirectroryCmd())
}

type addDirectroyCmdOptions struct {
	Files string
	To    string
}

// NewCloneCmd generates the `clone` command
func NewAddDirectroryCmd() *cobra.Command {
	s := addDirectroyCmdOptions{}
	c := &cobra.Command{
		Use:   "add-directory",
		Short: "Clones a repo, adds directory, creates PR",
		Long:  `Clones a repo, adds directory, creates PR`,
		RunE:  s.RunE,
	}

	c.Flags().StringVarP(&s.Files, "files", "f", "", "Directory to add")
	c.Flags().StringVarP(&s.To, "to", "g", "", "Directory where to add")

	c.MarkFlagRequired("prefix")
	c.MarkFlagRequired("org")
	c.MarkFlagRequired("files")

	viper.BindPFlags(c.Flags())

	return c
}

func (s *addDirectroyCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	globber, err := glob.New(glob.Default())
	if err != nil {
		return err
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	gh := github.NewClient(tc)

	page := 0
	hasMore := true

	for hasMore {
		repos, _, err := gh.Repositories.ListByOrg(ctx, org, &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{Page: page, PerPage: 20},
		})
		if err != nil {
			return err
		}

		for _, repo := range repos {
			if match, _ := globber.Match(prefix, *repo.Name); match {
				fmt.Println(*repo.Name)
				err := s.createBranch(*repo.CloneURL)
				fmt.Println("Branch Created", *repo.Name)
				if err != nil {
					log.Println(err)
				} else {
					title := "Add CSS Part 1"
					base := "master"
					head := "add-css"
					body := `Add start files for CSS Part 1`
					_, _, err = gh.PullRequests.Create(ctx, org, *repo.Name, &github.NewPullRequest{
						Title: &title,
						Head:  &head,
						Base:  &base,
						Body:  &body,
					})
				}
				if err != nil {
					log.Println(err)
				}

				fmt.Println("PR Created", *repo.Name)
				// anti-rate-limit
				time.Sleep(10 * time.Second)
			}
		}

		if len(repos) == 0 {
			hasMore = false
		}
		page++
	}

	return nil
}

func (s *addDirectroyCmdOptions) createBranch(url string) error {
	dir, err := ioutil.TempDir(os.TempDir(), "mass-pr-clone")
	if err != nil {
		return err
	}

	log.Println("Cloning to", url, dir)

	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		Auth: &http.BasicAuth{
			Username: "iloveoctocats", // yes, this can be anything except an empty string
			Password: authToken,
		},
		URL: url,
	})
	if err != nil {
		return fmt.Errorf("Error on clone: %w", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("Error on get worktree: %w", err)
	}

	err = gocopy.Copy(s.Files, path.Join(dir, s.To))
	if err != nil {
		return fmt.Errorf("Error on copy: %w", err)
	}

	err = w.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return err
	}

	// TODO: not hard code me
	author := object.Signature{Name: "Maartje Eyskens", Email: "maartje@eyskens.me", When: time.Now()}
	commit, err := w.Commit("Add CSS Part 1", &git.CommitOptions{Author: &author, Committer: &author})
	if err != nil {
		return err
	}

	//TODO: make branch name configurable
	ref := plumbing.NewHashReference("refs/heads/add-css", commit)
	err = r.Storer.SetReference(ref)
	if err != nil {
		return fmt.Errorf("Error on branch create: %w", err)
	}

	pushOptions := git.PushOptions{
		RefSpecs: []config.RefSpec{"refs/heads/add-css:refs/heads/add-css"},
		Auth: &http.BasicAuth{
			Username: "iloveoctocats", // yes, this can be anything except an empty string
			Password: authToken,
		},
	}

	err = r.Push(&pushOptions)
	if err != nil {
		return fmt.Errorf("Error on push: %w", err)
	}

	return nil
}
