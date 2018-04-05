package plugins

import (
	"os"
	"os/exec"

	"github.com/apex/log"
	"github.com/tympanix/supper/types"
	"gopkg.in/yaml.v2"
)

// Plugin is a YAML formatted struct describing external functionality
type Plugin struct {
	NameYaml string `yaml:"name"`
	ExecYaml string `yaml:"exec"`
}

// Run executes the plugin
func (p *Plugin) Run(s types.LocalSubtitle) error {
	cmd := exec.Command(shell[0], shell[1], escape(p.ExecYaml, s.Path()))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.WithError(err).WithField("plugin", p.Name()).
			Debugf("Plugin debug\n%s", string(out))
	}
	return err
}

// Name returns the name of the plugin
func (p *Plugin) Name() string {
	return p.NameYaml
}

// Config is a struct which contains YAML formatted configuration
type Config struct {
	PluginList []Plugin `yaml:"plugins,omitempty"`
}

// Plugins returns a list of plugins
func (c *Config) Plugins() []types.Plugin {
	pl := make([]types.Plugin, 0)
	for _, p := range c.PluginList {
		pl = append(pl, &p)
	}
	return pl
}

// Load reads a configuration file and enterprets the content
func Load(path string) (*Config, error) {
	c := Config{}

	if path == "" {
		return &c, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dec := yaml.NewDecoder(file)

	if err := dec.Decode(&c); err != nil {
		return nil, err
	}

	return &c, nil
}