package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port: %v", err)
	}

	// run関数を実行してHTTPサーバーを起動
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})

	// リクエストを送信・レスポンスを受信
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try request to %q", url)
	rsp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()

	// レスポンスボディを取得
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	// レスポンスを検証
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q:", want, got)
	}

	// サーバーを終了
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
