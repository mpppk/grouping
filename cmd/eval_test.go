package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mpppk/grouping/cmd"
	"github.com/spf13/afero"
)

func TestEval(t *testing.T) {
	cases := []struct {
		command string
		want    string
	}{
		{command: "eval --file ../testdata/dup_groups.csv", want: "2\n"},
		{command: "eval --file ../testdata/no_dup_groups.csv", want: "0\n"},
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		rootCmd, err := cmd.NewRootCmd(afero.NewMemMapFs())
		if err != nil {
			t.Errorf("failed to create rootCmd: %s", err)
		}
		rootCmd.SetOut(buf)
		cmdArgs := strings.Split(c.command, " ")
		rootCmd.SetArgs(cmdArgs)
		if err := rootCmd.Execute(); err != nil {
			t.Errorf("failed to execute rootCmd: %s", err)
		}

		get := buf.String()
		if c.want != get {
			t.Errorf("unexpected response: want:%q, get:%q", c.want, get)
		}
	}
}
