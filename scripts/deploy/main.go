package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v27/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var (
	githubToken = os.Getenv("GITHUB_TOKEN")
)

func strptr(s string) *string { return &s }

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	rerr, ok := err.(*github.ErrorResponse)
	if !ok {
		return false
	}
	return rerr.Response.StatusCode == http.StatusNotFound
}

func deploy(tag, path string) error {
	if len(githubToken) == 0 {
		return errors.New("GITHUB_TOKEN not set")
	}
	f, err := os.Open(path)
	if err != nil {
		errors.Wrap(err, "failed to open "+path)
	}
	defer f.Close()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	release, _, err := client.Repositories.GetReleaseByTag(ctx, "ueokande", "logbook", tag)
	if isNotFoundError(err) {
		release, _, err = client.Repositories.CreateRelease(ctx, "ueokande", "logbook", &github.RepositoryRelease{
			TagName: strptr(tag),
			Name:    strptr("Release " + tag),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create a release for "+tag)
		}
		fmt.Fprintln(os.Stderr, "Created release on", release.GetHTMLURL())
	} else if err != nil {
		return err
	}

	opt := &github.UploadOptions{
		Name: filepath.Base(path),
	}
	asset, _, err := client.Repositories.UploadReleaseAsset(ctx, "ueokande", "logbook", release.GetID(), opt, f)
	if err != nil {
		return errors.Wrap(err, "failed to upload assets "+opt.Name)
	}
	fmt.Fprintln(os.Stderr, "Uploaded on", asset.GetBrowserDownloadURL())
	return nil
}

func main() {
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s TAG ASSET_FILE\n", os.Args[0])
	}
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(2)
	}
	err := deploy(flag.Arg(0), flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %+v\n", os.Args[0], err)
		os.Exit(1)
	}
}
