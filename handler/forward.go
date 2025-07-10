package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/zsjinwei/http-transponder/config"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        500,
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConnsPerHost: 100,
	},
}

// 转发器核心处理
func ForwardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read body failed", http.StatusBadRequest)
			return
		}
		var wg sync.WaitGroup
		errs := make(chan error, len(config.GlobalConfig.ForwardURLs))
		for _, url := range config.GlobalConfig.ForwardURLs {
			wg.Add(1)
			target := url // 防止闭包问题
			go func() {
				defer wg.Done()
				req, err := http.NewRequest("POST", target, io.NopCloser(
					bytes.NewReader(body)))
				if err != nil {
					errs <- err
					return
				}
				// 可转发特定头信息，如果需要
				req.Header = r.Header.Clone()
				resp, err := client.Do(req)
				if err != nil {
					errs <- err
					return
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				if resp.StatusCode >= 300 {
					errs <- fmt.Errorf("forward %s, bad status: %d", target, resp.StatusCode)
				}
			}()
		}
		wg.Wait()
		close(errs)
		if len(errs) > 0 {
			http.Error(w, "forward failed, see logs", http.StatusBadGateway)
			for e := range errs {
				log.Println("Forward error:", e)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
