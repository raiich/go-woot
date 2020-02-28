package main

import (
	"./simulate"
	"time"
)

func show(x string) {
	time.Sleep(time.Second)
	print("\033[H\033[2J")
	println(x)
}

//func main1() {
//	site1 := w.NewSite("site1", "a")
//	former := clone(site1.Raw())
//	print("\033[H\033[2J")
//	println(site1.Value())
//
//	site1.GenerateIns(0, 'A')
//	latter := clone(site1.Raw())
//	show(diff(former, latter))
//	show(site1.Value())
//
//	site1.GenerateIns(2, 'B')
//	former, latter = latter, clone(site1.Raw())
//	show(diff(former, latter))
//	show(site1.Value())
//
//	site1.GenerateDel(1)
//	former, latter = latter, clone(site1.Raw())
//	show(diff(former, latter))
//	show(site1.Value())
//}

type Move struct {
	row int
	col int
}

func diff(base string, edited string) string {
	return string(diff2([]rune(base), []rune(edited)))
}

func diff2(base []rune, edited []rune) []rune {
	dp := make([][]int, len(base) + 1)
	move := make([][]Move, len(base) + 1)

	for i := range dp {
		dp[i] = make([]int, len(edited) + 1)
		move[i] = make([]Move, len(edited) + 1)
	}
	dp[0][0] = 0
	move[0][0] = Move{0, 0}

	for j := range edited {
		dp[0][j + 1] = dp[0][j]
		move[0][j + 1] = Move {0, 1}
	}

	for i := range base {
		dp[i + 1][0] = dp[i][0]
		move[i + 1][0] = Move {1, 0}

		for j := range edited {
			plus := 0
			if base[i] == edited[j] {
				plus = 1
			}

			a, b, c := dp[i + 1][j], dp[i][j + 1], dp[i][j] + plus
			x, y, z := Move {0, 1}, Move{1, 0}, Move{1, 1}
			if a <= c && b <= c {
				dp[i + 1][j + 1] = c
				move[i + 1][j + 1] = z
			} else if a <= b && c <= b {
				dp[i + 1][j + 1] = b
				move[i + 1][j + 1] = y
			} else if b <= a && c <= a {
				dp[i + 1][j + 1] = a
				move[i + 1][j + 1] = x
			}
		}
	}

	ret := make([]rune, 0)

	for i, j := len(base), len(edited); i != 0 && j != 0; {
		m := move[i][j]
		c := ' '
		if m.col == 1 && m.row == 1 {
			if base[i - 1] == edited[j - 1] {
				c = base[i - 1]
			} else {
				c = '$'
			}
			i -= 1
			j -= 1
		} else if m.col == 1 {
			c = '+'
			j -= 1
		} else if m.row == 1 {
			c = '-'
			i -= 1
		} else {
			panic("over")
		}
		ret = append([]rune{c}, ret...)
	}
	return ret
}

func main() {
	println(diff("abcdef", "axcdef"))
}

func main2() {
	trial := 3
	var slots = [][]rune{
		[]rune("abcdefghij"),
		[]rune("1234567890"),
		[]rune("^*-/_=+?:."),
	}

	peers := simulate.NewPeers(slots)
	out := make(chan simulate.Info)
	go simulate.Run(peers, out, trial)

	for line := range out {
		println(peers[line.Peer].NumSite, line.Value)
	}
}
