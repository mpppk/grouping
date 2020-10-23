package domain

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type MemberID int
type Member struct {
	ID   MemberID
	Name string
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
func findNameIndex(headers []string) (int, bool) {
	return findColumnIndex(headers, "NAME")
}

func findIdIndex(headers []string) (int, bool) {
	return findColumnIndex(headers, "ID")
}
