package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/zsjinwei/http-transponder/config"
	"github.com/zsjinwei/http-transponder/handler"
)

func main() {
	configPath := flag.String("config", "config.yaml", "config file path")
	flag.Parse()
	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatalf("load config failed: %v", err)
	}
	http.HandleFunc(config.GlobalConfig.ReceivePath, handler.ForwardHandler())
	log.Printf("Listening on %s, receiving at %s", config.GlobalConfig.ListenAddr, config.GlobalConfig.ReceivePath)
	if err := http.ListenAndServe(config.GlobalConfig.ListenAddr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
