package config

import (
	_ "embed"
	"fmt"
	"regexp"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"gopkg.in/yaml.v3"
)

type Config struct {
	JobNames         []*regexp.Regexp
	ExcludedJobNames []*regexp.Regexp
	JobNameMappings  map[*regexp.Regexp]string
}

type RawConfig struct {
	// job name regular expressions to include
	JobNames []string `json:"job_names,omitempty" yaml:"job_names,omitempty"`
	// job name regular expressions to exclude
	ExcludedJobNames []string `json:"excluded_job_names,omitempty" yaml:"excluded_job_names,omitempty"`
	// original job name regular expression => normalized job name
	JobNameMappings map[string]string `json:"job_name_mappings,omitempty" yaml:"job_name_mappings,omitempty"`
}

func (c *Config) Include(name string) bool {
	if len(c.ExcludedJobNames) > 0 {
		for _, pattern := range c.ExcludedJobNames {
			if pattern.MatchString(name) {
				return false
			}
		}
		return true
	}
	if len(c.JobNames) == 0 {
		return true
	}
	for _, pattern := range c.JobNames {
		if pattern.MatchString(name) {
			return true
		}
	}
	return false
}

func (c *Config) NormalizeJobName(name string) string {
	for pattern, mapped := range c.JobNameMappings {
		if matched := pattern.MatchString(name); matched {
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
	rCfg := &RawConfig{}
	if err := yaml.Unmarshal(b, rCfg); err != nil {
		return fmt.Errorf("unmarshal a configuration file: %w", err)
	}
	cfg.JobNames = make([]*regexp.Regexp, len(rCfg.JobNames))
	for i, pattern := range rCfg.JobNames {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("compile job name pattern: %w", slogerr.With(err, "pattern", pattern))
		}
		cfg.JobNames[i] = re
	}
	cfg.ExcludedJobNames = make([]*regexp.Regexp, len(rCfg.ExcludedJobNames))
	for i, pattern := range rCfg.ExcludedJobNames {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("compile excluded job name pattern: %w", slogerr.With(err, "pattern", pattern))
		}
		cfg.ExcludedJobNames[i] = re
	}
	cfg.JobNameMappings = make(map[*regexp.Regexp]string, len(rCfg.JobNameMappings))
	for original, mapped := range rCfg.JobNameMappings {
		re, err := regexp.Compile(original)
		if err != nil {
			return fmt.Errorf("compile job name mapping: %w", slogerr.With(err, "original", original, "mapped", mapped))
		}
		cfg.JobNameMappings[re] = mapped
	}
	return nil
}

//go:embed ghaperf.yaml
var initConfigContent []byte

const filePermission = 0o644

func Init(fs afero.Fs, path string) error {
	if f, err := afero.Exists(fs, path); err != nil {
		return fmt.Errorf("check existence of a config file: %w", err)
	} else if f {
		return nil
	}
	if err := afero.WriteFile(fs, path, initConfigContent, filePermission); err != nil {
		return fmt.Errorf("create a config file: %w", err)
	}
	return nil
}
