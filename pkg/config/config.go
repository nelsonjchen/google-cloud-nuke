package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/rebuy-de/aws-nuke/v2/pkg/types"

	"gopkg.in/yaml.v2"
)

type ResourceTypes struct {
	Targets  types.Collection `yaml:"targets"`
	Excludes types.Collection `yaml:"excludes"`
}

type Project struct {
	Filters       Filters       `yaml:"filters"`
	ResourceTypes ResourceTypes `yaml:"resource-types"`
	Presets       []string      `yaml:"presets"`
}

type Nuke struct {
	ProjectBlacklist []string                     `yaml:"project-blocklist"`
	Projects         map[string]Project           `yaml:"projects"`
	ResourceTypes    ResourceTypes                `yaml:"resource-types"`
	Presets          map[string]PresetDefinitions `yaml:"presets"`
	FeatureFlags     FeatureFlags                 `yaml:"feature-flags"`
}

type FeatureFlags struct {
	DisableDeletionProtection DisableDeletionProtection `yaml:"disable-deletion-protection"`
}

type DisableDeletionProtection struct {
	ComputeEngineInstance bool `yaml:"CEInstance"`
}

type PresetDefinitions struct {
	Filters Filters `yaml:"filters"`
}

func Load(path string) (*Nuke, error) {
	var err error

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := new(Nuke)
	err = yaml.UnmarshalStrict(raw, config)
	if err != nil {
		return nil, err
	}

	if err := config.resolveDeprecations(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Nuke) ResolveBlocklist() []string {
	return c.ProjectBlacklist
}

func (c *Nuke) HasBlocklist() bool {
	var blocklist = c.ResolveBlocklist()
	return blocklist != nil && len(blocklist) > 0
}

func (c *Nuke) InBlocklist(searchID string) bool {
	for _, blocklistID := range c.ResolveBlocklist() {
		if blocklistID == searchID {
			return true
		}
	}

	return false
}

func (c *Nuke) ValidateProject(projectId string, name string) error {
	if !c.HasBlocklist() {
		return fmt.Errorf("the config file contains an empty blocklist. " +
			"For safety reasons you need to specify at least one project ID. " +
			"This should be your production project(s)")
	}

	if c.InBlocklist(projectId) {
		return fmt.Errorf("you are trying to nuke the project with the ID %s, "+
			"but it is blocklisted. Aborting", projectId)
	}

	if strings.Contains(strings.ToLower(name), "prod") {
		return fmt.Errorf("you are trying to nuke an project with the alias '%s', "+
			"but it has the substring 'prod' in it. Aborting", name)
	}

	if _, ok := c.Projects[projectId]; !ok {
		return fmt.Errorf("your project ID '%s' isn't listed in the config. "+
			"Aborting", projectId)
	}

	return nil
}

func (c *Nuke) Filters(projectId string) (Filters, error) {
	project := c.Projects[projectId]
	filters := project.Filters

	if filters == nil {
		filters = Filters{}
	}

	if project.Presets == nil {
		return filters, nil
	}

	for _, presetName := range project.Presets {
		notFound := fmt.Errorf("Could not find filter preset '%s'", presetName)
		if c.Presets == nil {
			return nil, notFound
		}

		preset, ok := c.Presets[presetName]
		if !ok {
			return nil, notFound
		}

		filters.Merge(preset.Filters)
	}

	return filters, nil
}

func (c *Nuke) resolveDeprecations() error {
	return nil
}
