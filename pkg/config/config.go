package config

import (
	"fmt"
	"path"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"gopkg.in/yaml.v3"
)

type Config struct {
	// job name glob patterns to include
	JobNames []string `json:"job_names,omitempty" yaml:"job_names,omitempty"`
	// job name glob patterns to exclude
	ExcludedJobNames []string `json:"excluded_job_names,omitempty" yaml:"excluded_job_names,omitempty"`
	// original job name glob pattern => normalized job name
	JobNameMappings map[string]string `json:"job_name_mappings,omitempty" yaml:"job_name_mappings,omitempty"`
}

func (c *Config) Validate() error {
	for _, pattern := range c.JobNames {
		if _, err := path.Match(pattern, "test"); err != nil {
			return fmt.Errorf("invalid job name pattern: %w", slogerr.With(err, "pattern", pattern))
		}
	}
	for _, pattern := range c.ExcludedJobNames {
		if _, err := path.Match(pattern, "test"); err != nil {
			return fmt.Errorf("invalid job name pattern: %w", slogerr.With(err, "pattern", pattern))
		}
	}
	for pattern := range c.JobNameMappings {
		if _, err := path.Match(pattern, "test"); err != nil {
			return fmt.Errorf("invalid job name pattern: %w", slogerr.With(err, "pattern", pattern))
		}
	}
	return nil
}

func (c *Config) Include(name string) bool {
	if len(c.ExcludedJobNames) > 0 {
		for _, pattern := range c.ExcludedJobNames {
			matched, err := path.Match(pattern, name)
			if err != nil {
				continue
			}
			if matched {
				return false
			}
		}
		return true
	}
	if len(c.JobNames) == 0 {
		return true
	}
	for _, pattern := range c.JobNames {
		matched, err := path.Match(pattern, name)
		if err != nil {
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

func (c *Config) NormalizeJobName(name string) string {
	for pattern, mapped := range c.JobNameMappings {
		if matched, err := path.Match(pattern, name); err != nil {
			continue
		} else if matched {
			return mapped
		}
	}
	return name
}

func Read(fs afero.Fs, path string, cfg *Config) error {
	b, err := afero.ReadFile(fs, path)
	if err != nil {
		return fmt.Errorf("read a configuration file: %w", err)
	}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return fmt.Errorf("unmarshal a configuration file: %w", err)
	}
	return nil
}
