package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/yano-kentaro/go_todo_app/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}

	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()

	// HTTPサーバーを起動
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	// シグナルを受け取るためのコンテキストを作成
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// シグナルを受け取ったらサーバーを終了
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// サーバーの終了を待つ
	return eg.Wait()
}
