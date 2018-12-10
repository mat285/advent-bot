package advent

import (
	"fmt"
	"sort"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/request"
)

const (
	leaderboardURLFormat = "https://adventofcode.com/2018/leaderboard/private/view/%s.json"

	EnvVarLeaderBoardID = "ADVENT_LEADERBOARD_ID"
)

type Response struct {
	Members map[string]*Member `json:"members"`
}

type Member struct {
	Name       string `json:"name"`
	LocalScore int    `json:"local_score"`
}

type Board struct {
	Members []*Member
}

func GetLeaderBoard() (*Board, error) {
	req := request.Get(fmt.Sprintf(leaderboardURLFormat, env.Env().String(EnvVarLeaderBoardID)))
	resp := &Response{}
	_, err := req.JSONError(resp)
	if err != nil {
		return nil, err
	}
	b := &Board{}
	for _, m := range resp.Members {
		b.Members = append(b.Members, m)
	}
	sort.Sort(b)
	return b, nil
}

func (b *Board) Len() int {
	return len(b.Members)
}

func (b *Board) Less(i, j int) bool {
	return b.Members[j].LocalScore < b.Members[i].LocalScore
}

func (b *Board) Swap(i, j int) {
	t := b.Members[i]
	b.Members[i] = b.Members[j]
	b.Members[j] = t
}
