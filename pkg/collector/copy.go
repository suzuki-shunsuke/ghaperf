package collector

import (
	"errors"
	"fmt"
	"io"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

const maxFileSize = 1073741824 // 1GB
var errFileTooLarge = errors.New("file is too large")

func copySafe(dst io.Writer, src io.Reader) error {
	writeCount, err := io.CopyN(dst, src, maxFileSize)
	// io.CopyN returns io.EOF error when the file size is less than maxFileSize
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("copy a log file: %w", err)
	}
	if writeCount >= maxFileSize {
		return slogerr.With(errFileTooLarge, "file_size", writeCount, "max_file_size", maxFileSize) //nolint:wrapcheck
	}
	return nil
}
