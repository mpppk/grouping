package domain

import (
	"reflect"
	"testing"
)

func Test_parseGroupLines(t *testing.T) {
	type args struct {
		lines [][]string
	}
	tests := []struct {
		name    string
		args    args
		want    []Groups
		wantErr bool
	}{
		{
			args: args{
				lines: [][]string{
					{"NAME", "1st", "2nd"},
					{"alice", "1", "2"},
					{"bob", "1", "2"},
					{"carol", "2", "1"},
					{"dave", "2", "1"},
				},
			},
			want: []Groups{
				{
					1: &Group{
						ID: 1,
						members: []*Member{
							{Name: "alice"},
							{Name: "bob"},
						},
					},
					2: &Group{
						ID: 2,
						members: []*Member{
							{Name: "carol"},
							{Name: "dave"},
						},
					},
				},
				{
					1: &Group{
						ID: 1,
						members: []*Member{
							{Name: "carol"},
							{Name: "dave"},
						},
					},
					2: &Group{
						ID: 2,
						members: []*Member{
							{Name: "alice"},
							{Name: "bob"},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGroupLines(tt.args.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGroupLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseGroupLines() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountDupMemberPairs(t *testing.T) {
	type args struct {
		groupsList []Groups
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			args: args{
				groupsList: []Groups{
					{
						1: &Group{
							members: []*Member{
								{Name: "alice"},
								{Name: "bob"},
							},
						},
						2: &Group{
							members: []*Member{
								{Name: "carol"},
								{Name: "dave"},
							},
						},
					},
					{
						1: &Group{
							members: []*Member{
								{Name: "carol"},
								{Name: "dave"},
							},
						},
						2: &Group{
							members: []*Member{
								{Name: "alice"},
								{Name: "bob"},
							},
						},
					},
				},
			},
			want:    2,
			wantErr: false,
		},
		{
			args: args{
				groupsList: []Groups{
					{
						1: &Group{
							members: []*Member{
								{Name: "alice"},
								{Name: "bob"},
							},
						},
						2: &Group{
							members: []*Member{
								{Name: "carol"},
								{Name: "dave"},
							},
						},
					},
					{
						1: &Group{
							members: []*Member{
								{Name: "alice"},
								{Name: "carol"},
							},
						},
						2: &Group{
							members: []*Member{
								{Name: "bob"},
								{Name: "dave"},
							},
						},
					},
				},
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CountDupMemberPairs(tt.args.groupsList)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Eval() got = %v, want %v", got, tt.want)
			}
		})
	}
}
