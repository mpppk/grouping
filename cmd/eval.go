package cmd

import (
	"fmt"

	"github.com/mpppk/grouping/cmd/option"
	"github.com/mpppk/grouping/domain"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newEvalCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "eval",
		Short: "evaluate groups",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewEvalCmdConfigFromViper(args)
			if err != nil {
				return err
			}

			groupsList, err := domain.ParseGroupFile(conf.File)
			if err != nil {
				return fmt.Errorf("failed to parse group file from %s: %w", conf.File, err)
			}

			cnt, err := domain.CountDupMemberPairs(groupsList)
			if err != nil {
				return fmt.Errorf("failed to count dup member pairs: %w", err)
			}

			cmd.Println(cnt)

			return nil
		},
	}

	registerEvalCommandFlags := func(cmd *cobra.Command) error {
		flags := []option.Flag{
			&option.StringFlag{
				BaseFlag: &option.BaseFlag{
					Name:  "file",
					Usage: "file",
				},
			},
		}
		return option.RegisterFlags(cmd, flags)
	}

	if err := registerEvalCommandFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newEvalCmd)
}
