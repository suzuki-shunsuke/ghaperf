package github

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/go-github/v79/github"
	"github.com/suzuki-shunsuke/ghtkn-go-sdk/ghtkn"
	"golang.org/x/oauth2"
)

type Client struct {
	actions ActionsService
	http    *http.Client
}

type InputNew struct {
	GHTKNEnabled bool
	AccessToken  string
}

type (
	Response                  = github.Response
	Jobs                      = github.Jobs
	ListWorkflowJobsOptions   = github.ListWorkflowJobsOptions
	WorkflowJob               = github.WorkflowJob
	WorkflowRun               = github.WorkflowRun
	ListOptions               = github.ListOptions
	TaskStep                  = github.TaskStep
	WorkflowRunAttemptOptions = github.WorkflowRunAttemptOptions
	ListWorkflowRunsOptions   = github.ListWorkflowRunsOptions
	WorkflowRuns              = github.WorkflowRuns
)

func New(ctx context.Context, logger *slog.Logger, input *InputNew) (*Client, error) {
	httpClient, err := newHTTPClient(ctx, logger, input)
	if err != nil {
		return nil, err
	}
	gh := github.NewClient(httpClient)
	return &Client{
		actions: gh.Actions,
		// This is used to download logs with redirect URLs.
		// The authentication fails if httpClient is used, so http.DefaultClient is used.
		// > 401 InvalidAuthenticationInfo - Server failed to authenticate the request. Please refer to the information in the www-authenticate header.
		http: http.DefaultClient,
	}, nil
}

func newHTTPClient(ctx context.Context, logger *slog.Logger, input *InputNew) (*http.Client, error) {
	ts, err := newTokenSource(logger, input)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(ctx, ts), nil
	// return makeRetryable(oauth2.NewClient(ctx, ts), logger), nil
}

var errAccessTokenRequired = errors.New("access token is required")

func newTokenSource(logger *slog.Logger, input *InputNew) (oauth2.TokenSource, error) {
	if input.AccessToken != "" {
		return oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: input.AccessToken},
		), nil
	}
	if input.GHTKNEnabled {
		client := ghtkn.New()
		return client.TokenSource(logger, &ghtkn.InputGet{}), nil
	}
	return nil, errAccessTokenRequired
}

/*
func makeRetryable(client *http.Client, logger *slog.Logger) *http.Client {
	c := retryablehttp.NewClient()
	c.HTTPClient = client
	c.Logger = logger
	return c.StandardClient()
}
*/
