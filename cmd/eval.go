package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/mpppk/cli-template/cmd/option"
	"github.com/spf13/afero"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func newEvalCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "eval",
		Short: "evaluate groups",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// groupごとの
			conf, err := option.NewEvalCmdConfigFromViper(args)
			if err != nil {
				return err
			}

			groupsList, err := parseGroupFile(conf.File)
			if err != nil {
				return fmt.Errorf("failed to parse group file from %s: %w", conf.File, err)
			}

			fmt.Println(groupsList)

			return nil
		},
	}

	if err := registerEvalCommandFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

type MemberID int
type Member struct {
	ID   MemberID
	Name string
}

type GroupID int
type Group struct {
	ID      GroupID
	members []*Member
}

type Groups map[GroupID] *Group

func NewGroups() Groups {
	return map[GroupID] *Group{}
}

func (g Groups) addGroup(member *Member, id GroupID) {
	if _, ok := g[id]; !ok {
		g[id] = &Group{ID: id}
	}
	g[id].members = append(g[id].members, member)
}

func (g Groups) addGroups(member *Member, idList ...GroupID) {
	for _, id := range idList {
		g.addGroup(member, id)
	}
}

func eval(members []*Member, groups []*Group) int {
	m := map[MemberID]map[MemberID]struct{}{}
	for _, group := range groups {
		for i, member := range group.members {
			for j := i + 1; j < len(group.members); j++ {
				member2 := group.members[j]
				m[member.ID][member2.ID] = struct{}{}
			}
		}
	}

	cnt := 0
	for _, m2 := range m {
		cnt += len(m2)
	}
	return cnt
}

func parseGroupFile(filePath string) ([]Groups, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse csv from %s: %w", filePath, err)
	}
	return parseGroupLines(lines)
}

func parseGroupLines(lines [][]string) ([]Groups, error) {
	if err := validateMemberFile(lines); err != nil {
		return nil, fmt.Errorf("failed to parse group file: %w", err)
	}
	headers := lines[0]
	nameIndex, ok := findNameIndex(headers)
	if !ok {
		return nil, fmt.Errorf("failed to find NAME column")
	}

	groupsList := newGroupsList(len(headers)-1)
	for _, line := range lines[1:] {
		name, groupIDList, err := parseGroupLine(line, nameIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to parse group line from %s: %w", line, err)
		}

		for i, id := range groupIDList {
			groupsList[i].addGroup(&Member{Name: name}, id)
		}
	}
	return groupsList, nil
}

func newGroupsList(length int) []Groups {
	groupsList := make([]Groups, length)
	for i := 0; i < length; i++ {
		groupsList[i] = NewGroups()
	}
	return groupsList
}

func parseGroupLine(line []string, nameIndex int) (name string, groupIDList []GroupID, err error) {
	for i, groupIDOrName := range line {
		if i == nameIndex {
			name = groupIDOrName
		} else {
			groupID, err := strconv.Atoi(groupIDOrName)
			if err != nil {
				return "", nil, fmt.Errorf("failed to convert group ID to int from %s: %w", groupIDOrName, err)
			}
			groupIDList = append(groupIDList, GroupID(groupID))
		}
	}
	return
}

func parseMemberFile(filePath string) ([]*Member, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse csv from %s: %w", filePath, err)
	}
	return parseMemberLines(lines)
}

func parseMemberLines(lines [][]string) ([]*Member, error) {
	if err := validateMemberFile(lines); err != nil {
		return nil, fmt.Errorf("failed to parse group file: %w", err)
	}
	headers := lines[0]
	data := lines[1:]
	idIndex, ok := findIdIndex(headers)
	if !ok {
		return nil, fmt.Errorf("failed to find ID column")
	}
	nameIndex, ok := findNameIndex(headers)
	if !ok {
		return nil, fmt.Errorf("failed to find NAME column")
	}

	var members []*Member
	for _, d := range data {
		member, err := parseMemberLine(d, idIndex, nameIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to parse member line: %w", err)
		}
		members = append(members, member)
	}
	return members, nil
}

func parseMemberLine(line []string, idIndex, nameIndex int) (*Member, error) {
	idStr := line[idIndex]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse id column number(%s): %w", idStr, err)
	}
	return &Member{
		ID:   MemberID(id),
		Name: line[nameIndex],
	}, nil
}

func validateMemberFile(lines [][]string) error {
	if len(lines) == 0 {
		return errors.New("zero lines")
	}
	colNum := len(lines[0])
	for _, line := range lines {
		if len(line) != colNum {
			return fmt.Errorf("inbalid row found. expected col num: %d, actual: %d", colNum, len(line))
		}
	}
	return nil
}

func findColumnIndex(headers []string, name string) (int, bool) {
	for i, header := range headers {
		if header == name {
			return i, true
		}
	}
	return 0, false
}

func findNameIndex(headers []string) (int, bool) {
	return findColumnIndex(headers, "NAME")
}

func findIdIndex(headers []string) (int, bool) {
	return findColumnIndex(headers, "ID")
}

func registerEvalCommandFlags(cmd *cobra.Command) error {
	flags := []option.Flag{
		&option.BoolFlag{
			BaseFlag: &option.BaseFlag{
				Name:  "file",
				Usage: "file",
			},
			Value: false,
		},
	}
	return option.RegisterFlags(cmd, flags)
}

func init() {
	cmdGenerators = append(cmdGenerators, newEvalCmd)
}
