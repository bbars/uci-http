package engine

import (
)

type FindMoveReq struct {
	Fen string `json:"fen"`
	Depth int `json:"depth"`
	Time int `json:"time"`
}

type FindMoveRes struct {
	Bestmove string `json:"bestmove"`
	Score int `json:"score"`
	Pv []string `json:"pv"`
}

type MoveFinder interface {
	FindMove(req *FindMoveReq) (*FindMoveRes, error)
}
