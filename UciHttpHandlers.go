package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	
	// "github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

type UciHttpHandlers interface {
	HandleGo(w http.ResponseWriter, r *http.Request)
}

type ResponseRes struct {
	Res interface{} `json:"res"`
}

type ResponseErr struct {
	Err string `json:"err"`
}

type CmdPositionAndGo struct {
	Position uci.CmdPosition `json:"position"`
	Go uci.CmdGo `json:"go"`
}


type requestProcessor func (r *http.Request) (interface{}, error)

type uciHttpHandlers struct {
	engineBin string
	maxTime time.Duration
}

func NewUciHttpHandlers(engineBin string, maxTime time.Duration) UciHttpHandlers {
	return &uciHttpHandlers{
		engineBin: engineBin,
		maxTime: maxTime,
	}
}

func (this *uciHttpHandlers) newEngine() (*uci.Engine, error) {
	return uci.New(this.engineBin)
}

func (this *uciHttpHandlers) processRequest(w http.ResponseWriter, r *http.Request, fn requestProcessor) {
	res, err := fn(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ResponseErr{
			Err: err.Error(),
		})
	} else {
		json.NewEncoder(w).Encode(ResponseRes{
			Res: res,
		})
	}
}

func (this *uciHttpHandlers) HandleGo(w http.ResponseWriter, r *http.Request) {
	this.processRequest(w, r, func (r *http.Request) (interface{}, error) {
		var cmd *CmdPositionAndGo
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			return nil, err
		}
		engine, err:= this.newEngine()
		if err != nil {
			return nil, err
		}
		defer engine.Close()
		if err := engine.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
			return nil, err
		}
		if this.maxTime > 0 && (cmd.Go.MoveTime == 0 || cmd.Go.MoveTime > this.maxTime) {
			cmd.Go.MoveTime = this.maxTime
		}
		if err := engine.Run(cmd.Position, cmd.Go); err != nil {
			return nil, err
		}
		return engine.SearchResults(), nil
	})
}
