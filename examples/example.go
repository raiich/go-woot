package main

import (
	"./simulate"
	"time"
)

type move struct {
	row int
	col int
}

func Diff(base string, edited string) string {
	return string(arrayDiff([]rune(base), []rune(edited)))
}

func arrayDiff(base []rune, edited []rune) []rune {
	dp := make([][]int, len(base) + 1)
	route := make([][]move, len(base) + 1)

	for i := range dp {
		dp[i] = make([]int, len(edited) + 1)
		route[i] = make([]move, len(edited) + 1)
	}
	dp[0][0] = 0
	route[0][0] = move{0, 0}

	for j := range edited {
		dp[0][j + 1] = dp[0][j]
		route[0][j + 1] = move{0, 1}
	}

	for i := range base {
		dp[i + 1][0] = dp[i][0]
		route[i + 1][0] = move{1, 0}

		for j := range edited {
			plus := 0
			if base[i] == edited[j] {
				plus = 1
			} else {
				plus=-1
			}

			a, b, c := dp[i + 1][j], dp[i][j + 1], dp[i][j] + plus
			x, y, z := move{0, 1}, move{1, 0}, move{1, 1}
			if a <= b && c <= b {
				dp[i + 1][j + 1] = b
				route[i + 1][j + 1] = y
			} else if a <= c && b <= c {
				dp[i + 1][j + 1] = c
				route[i + 1][j + 1] = z
			} else if b <= a && c <= a {
				dp[i + 1][j + 1] = a
				route[i + 1][j + 1] = x
			}
		}
	}

	ret := make([]rune, 0)

	for i, j := len(base), len(edited); i != 0 || j != 0; {
		m := route[i][j]
		c := ' '
		if m.col == 1 && m.row == 1 {
			if base[i - 1] == edited[j - 1] {
				c = base[i - 1]
			} else {
				c = '$'
			}
			i -= 1
			j -= 1
		} else if m.row == 1 {
			c = '*'
			i -= 1
		} else if m.col == 1 {
			c = '_'
			j -= 1
		} else {
			panic("over")
		}
		ret = append([]rune{c}, ret...)
	}
	return ret
}

func show(lines []string, changed int, op string, delim string) {
	time.Sleep(time.Second / 2)
	print("\033[H\033[2J")
	// println(delim)
	// print("\033[H\033[2J")
	for i, line := range lines {
		if i == changed {
			println(i, op, line)
		} else {
			println(i, "     ", line)
		}
	}
}

func main() {
	trial := 4
	var slots = [][]rune{
		[]rune("ABCDEFGHIJ"),
		[]rune("abcdefghij"),
		[]rune("1234567890"),
	}

	peers := simulate.NewPeers(slots)
	out := make(chan simulate.Info)
	go simulate.Run(peers, out, trial)

	display := make([]string, len(peers))
	for i, p := range peers {
		display[i] = p.Value()
	}

	for line := range out {
		before := display[line.Peer]

		diff := Diff(before, line.Value)
		if before != diff {
			display[line.Peer] = diff
			// println()
			// println("    :", before, " ; ", line.Value)
			show(display, line.Peer, line.Operation, "----")
			display[line.Peer] = line.Value
			show(display, line.Peer, line.Operation, "~")
		}
	}
}

func TestDiff() {
	if Diff("abcdef", "axcdef") != "a+-cdef" {
		panic("")
	}
	if Diff("abc", "12345") != "+++++---" {
		panic("")
	}
	if Diff("abc", "") != "---" {
		panic("")
	}
	if Diff("", "") != "" {
		panic("")
	}
}

