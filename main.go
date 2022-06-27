package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func corsMiddleware(w http.ResponseWriter, r *http.Request, allowOrigin string, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
	if r.Method == http.MethodOptions {
		return
	}
	next(w, r)
}

func main() {
	listenPtr := flag.String("listen", ":80", "HTTP listen [host]:port")
	engineBinPtr := flag.String("engineBin", "./stockfish/stockfish_15_linux_x64", "UCI engine binary path")
	allowOriginPtr := flag.String("allowOrigin", "*", "Value for HTTP header Access-Control-Allow-Origin")
	maxTimePtr := flag.Int("maxTime", 0, "Max time limit (ns)")
	helpPtr := flag.Bool("help", false, "Show usage info")
	
	flag.Parse()
	if *helpPtr {
		flag.Usage()
		return
	}
	
	log.Println("listen", *listenPtr)
	log.Println("engineBin", *engineBinPtr)
	log.Println("allowOrigin", *allowOriginPtr)
	log.Println("maxTime", *maxTimePtr)
	
	uciHttpHandlers := NewUciHttpHandlers(*engineBinPtr, time.Duration(*maxTimePtr))
	
	http.HandleFunc("/go", func (w http.ResponseWriter, r *http.Request) {
		corsMiddleware(w, r, *allowOriginPtr, uciHttpHandlers.HandleGo)
	})
	
	log.Fatal(http.ListenAndServe(*listenPtr, nil))
}
