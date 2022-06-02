package main

import (
	"flag"
	"log"
	"net/http"
	"stockfish-http/engine"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	listenPtr := flag.String("listen", ":80", "HTTP listen [host]:port")
	stockfishBinPtr := flag.String("stockfishBin", "./stockfish/stockfish_15_linux_x64", "Stockfish binary")
	allowOriginPtr := flag.String("allowOrigin", "*", "Value for HTTP header Access-Control-Allow-Origin")
	defaultDepthPtr := flag.Int("defaultDepth", 0, "Default depth limit")
	defaultTimePtr := flag.Int("defaultTime", 0, "Default time limit (ms)")
	helpPtr := flag.Bool("help", false, "Show usage info")
	
	flag.Parse()
	if *helpPtr {
		flag.Usage()
		return
	}
	
	log.Println("listen", *listenPtr)
	log.Println("stockfishBin", *stockfishBinPtr)
	log.Println("allowOrigin", *allowOriginPtr)
	log.Println("defaultDepth", *defaultDepthPtr)
	log.Println("defaultTime", *defaultTimePtr)
	
	stockfishFindmoveHandler := &StockfishFindmoveHandler{
		stockfishBin: *stockfishBinPtr,
		defaultReq: &engine.FindMoveReq{
			Depth: *defaultDepthPtr,
			Time: *defaultTimePtr,
		},
	}
	http.HandleFunc("/findmove", func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", *allowOriginPtr)
		stockfishFindmoveHandler.ServeHTTP(w, r)
	})
	
	log.Fatal(http.ListenAndServe(*listenPtr, nil))
}
