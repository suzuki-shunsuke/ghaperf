package github

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

// GetWorkflowJobLogs downloads the logs for a workflow job.
// It returns an io.ReadCloser which must be closed by the caller.
func (c *Client) GetWorkflowJobLogs(ctx context.Context, owner, repo string, jobID int64) (io.ReadCloser, error) {
	link, _, err := c.actions.GetWorkflowJobLogs(ctx, owner, repo, jobID, maxRedirects)
	if err != nil {
		return nil, fmt.Errorf("get workflow job logs redirect URL: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create a http request: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download workflow job logs: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read error response body: %w", slogerr.With(err, "status_code", resp.StatusCode))
		}
		return nil, fmt.Errorf("download workflow job logs: %w", slogerr.With(errInvalidStatusCode, "status_code", resp.StatusCode, "response_body", string(b)))
	}
	return resp.Body, nil
}
