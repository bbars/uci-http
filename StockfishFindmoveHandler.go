package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"stockfish-http/engine"
	"strconv"
	"log"
)

type StockfishFindmoveHandler struct {
	stockfishBin string
	defaultReq *engine.FindMoveReq
}

func (this *StockfishFindmoveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := this.processRequest(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err.Error())
	} else {
		json.NewEncoder(w).Encode(res)
	}
}

func (this *StockfishFindmoveHandler) processRequest(r *http.Request) (interface{}, error) {
	if r.Method != "GET" {
		return nil, fmt.Errorf("Invalid HTTP method: %s", r.Method)
	}
	params := r.URL.Query()
	
	req := &engine.FindMoveReq{
		Fen: params.Get("fen"),
		Depth: this.defaultReq.Depth,
		Time: this.defaultReq.Time,
	}
	
	if params.Has("depth") {
		depth, err := strconv.Atoi(params.Get("depth"))
		if err != nil {
			return nil, err
		}
		req.Depth = depth
	}
	if params.Has("time") {
		time, err := strconv.Atoi(params.Get("time"))
		if err != nil {
			return nil, err
		}
		req.Time = time
	}
	
	return this.FindMove(req)
}

func (this *StockfishFindmoveHandler) FindMove(req *engine.FindMoveReq) (res *engine.FindMoveRes, err error) {
	cmd := exec.Command(this.stockfishBin)
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	if req.Fen != "" {
		_, err = fmt.Fprintf(stdin, "position fen %s\n", req.Fen)
		if err != nil {
			return nil, err
		}
	}
	
	_, err = fmt.Fprintf(stdin, "go depth %d movetime %d\n", req.Depth, req.Time)
	
	res = &engine.FindMoveRes{}
	
	var line string
	scores := make(map[string]int)
	pvs := make(map[string][]string)
	scanner := bufio.NewScanner(stdout)
loop:
	for scanner.Scan() {
		line = scanner.Text()
		// print(line + "\n")
		words := regexp.MustCompile("[\\s\\n]+").Split(line, -1)
		switch words[0] {
		case "bestmove":
			res.Bestmove = words[1]
			break loop
		case "info":
			var pv0 string
			for i := len(words) - 1; i >= 0; i-- {
				switch words[i] {
				case "pv":
					pv0 = words[i + 1]
					pvs[pv0] = words[i + 1:len(words)]
				break
				case "score":
					if words[i + 1] == "cp" {
						i, err := strconv.Atoi(words[i + 2])
						if err == nil {
							scores[pv0] = i
						}
					}
				break
				}
			}
		}
	}
	res.Score = scores[res.Bestmove]
	res.Pv = pvs[res.Bestmove]
	if err = scanner.Err(); err != nil {
		return res, err
	}
	return res, nil
}
