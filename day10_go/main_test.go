package main

import "testing"

func Test_parseLamps(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		s    string
		n    int
		want bits
	}{
		{
			name: "test1",
			s:    ".#.#",
			n:    4,
			want: bits(0x0a),
		},
		{
			name: "test2",
			s:    "#.#..#.#",
			n:    8,
			want: bits(0xa5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotn, gots := parseLamps(tt.s)
			if gotn != tt.n || gots != tt.want {
				t.Errorf("parseLamps() = (%v, %v), want (%v, %v)", gotn, gots, tt.n, tt.want)
			}
		})
	}
}

func Test_parseSwitch(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s    string
		want bits
	}{
		{
			name: "test1",
			s:    "0,2",
			want: bits(5),
		},
		{
			name: "test2",
			s:    "0,1,2,7",
			want: bits(135),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSwitch(tt.s)
			if got != tt.want {
				t.Errorf("parseSwitch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseMachine(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		line string
		want machine
	}{
		{
			name: "test1",
			line: "[.#.#](0,2){1,2}",
			want: machine{
				lamps:    bits(0xA),
				switches: []bits{bits(5)},
				joltages: []int{1, 2},
			},
		},
		{
			name: "test2",
			line: "[#.###] (0,1,3) (0,1,4) (0,2,3,4) (1,2) {20,29,13,6,16}",
			want: machine{
				lamps:    bits(0xA5),
				switches: []bits{bits(135)},
				joltages: []int{1, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseMachine(tt.line)
			// TODO: update the condition below to compare got with tt.want.
			if got.lamps != tt.want.lamps || len(got.switches) != len(tt.want.switches) || len(got.joltages) != len(tt.want.joltages) {
				t.Errorf("parseMachine() = %v, want %v", got, tt.want)
				got.print()
			}
		})
	}
}

func Test_part1(t *testing.T) {
	line := "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}"
	m := parseMachine(line)
	if d := m.search(1); d != -1 {
		t.Logf("solved in %d steps: %s", d, m)
	} else {
		t.Errorf("failed to solve machine: %s", m)
		m.print()
	}
}
