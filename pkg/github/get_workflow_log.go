package github

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (c *Client) download(ctx context.Context, link string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return nil, fmt.Errorf("create a http request: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download a file: %w", err)
	}
	return resp, nil
}

func (c *Client) GetWorkflowRunLogs(ctx context.Context, owner, repo string, runID int64, attempt int) ([]*zip.File, error) {
	link, res, err := c.actions.GetWorkflowRunAttemptLogs(ctx, owner, repo, runID, attempt, maxRedirects)
	if err != nil {
		if res == nil {
			return nil, fmt.Errorf("get workflow run logs redirect URL: %w", err)
		}
		if res.StatusCode == http.StatusGone {
			return nil, fmt.Errorf("get workflow run logs redirect URL: %w", slogerr.With(ErrLogHasGone, "status_code", res.StatusCode))
		}
		return nil, fmt.Errorf("get workflow run logs redirect URL: %w", err)
	}
	resp, err := c.download(ctx, link.String())
	if err != nil {
		return nil, fmt.Errorf("download workflow run logs: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read error response body: %w", slogerr.With(err, "status_code", resp.StatusCode))
		}
		return nil, fmt.Errorf("download workflow run logs: %w", slogerr.With(errInvalidStatusCode, "status_code", resp.StatusCode, "response_body", string(b)))
	}
	return readZip(resp.Body)
}

func readZip(body io.ReadCloser) ([]*zip.File, error) {
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	body.Close()
	readerAt := bytes.NewReader(buf)
	zr, err := zip.NewReader(readerAt, int64(len(buf)))
	if err != nil {
		return nil, fmt.Errorf("create a zip reader: %w", err)
	}
	files := make([]*zip.File, 0, len(zr.File))
	for _, f := range zr.File {
		if strings.HasSuffix(f.Name, "/") {
			continue // skip directories
		}
		if filepath.Dir(f.Name) != "." {
			continue // skip files in sub directories
		}
		files = append(files, f)
	}
	return files, nil
}
