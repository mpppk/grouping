package domain

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type GroupID int
type Group struct {
	ID      GroupID
	members []*Member
}

func (g *Group) getMemberPairs() (pairs [][]*Member) {
	for i, member := range g.members {
		for j := i + 1; j < len(g.members); j++ {
			members := []*Member{member, g.members[j]}
			sort.Slice(members, func(i, j int) bool {
				return members[i].Name > members[j].Name
			})
			pairs = append(pairs, members)
		}
	}
	return
}

type Groups map[GroupID]*Group

func NewGroups() Groups {
	return map[GroupID]*Group{}
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

type PairMap map[string]map[string]int

func (p PairMap) AddPairs(pairs [][]*Member) error {
	for _, pair := range pairs {
		if err := p.AddPair(pair); err != nil {
			return err
		}
	}
	return nil
}

func (p PairMap) AddPair(pair []*Member) error {
	if len(pair) != 2 {
		return fmt.Errorf("invalid pair because length is not 2. actual %d", len(pair))
	}
	name0, name1 := pair[0].Name, pair[1].Name
	if _, ok := p[name0]; !ok {
		p[name0] = map[string]int{}
	}

	p[name0][name1]++
	return nil
}

func (p PairMap) CountDup() (cnt int) {
	for _, m := range p {
		for _, c := range m {
			cnt += c - 1
		}
	}
	return
}

func CountDupMemberPairs(groupsList []Groups) (int, error) {
	pairMap := PairMap{}
	for _, groups := range groupsList {
		for _, group := range groups {
			if err := pairMap.AddPairs(group.getMemberPairs()); err != nil {
				return 0, fmt.Errorf("failed to add pairs to PairMap: %w", err)
			}
		}
	}

	return pairMap.CountDup(), nil
}

func ParseGroupFile(filePath string) ([]Groups, error) {
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

	groupsList := newGroupsList(len(headers) - 1)
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

func findColumnIndex(headers []string, name string) (int, bool) {
	for i, header := range headers {
		if header == name {
			return i, true
		}
	}
	return 0, false
}
