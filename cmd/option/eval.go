package option

import (
	"fmt"
	"github.com/spf13/viper"
)

// EvalCmdConfig is config for sum command
type EvalCmdConfig struct {
	File string
}

// NewEvalCmdConfigFromViper generate config for sum command from viper
func NewEvalCmdConfigFromViper(args []string) (*EvalCmdConfig, error) {
	var conf EvalCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("failed to create sum cmd config: %w", err)
	}

	return &conf, nil
}

func (c *EvalCmdConfig) validate() error {
	return nil
}
