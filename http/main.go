package main

import (
	"github.com/oklog/run"
	"net/http"
	"github.com/francis-tang/tutorials/http/api"
	"context"
	"time"
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

func main(){
	nsa:=api.NewNameServiceApi()

	var g run.Group
	{
		server:=http.Server{
			Addr:"127.0.0.1:25000",
			Handler:nsa,
		}

		g.Add(func() error {
			return server.ListenAndServe()
		},func (error){
			ctx,cancel := context.WithTimeout(context.Background(),time.Second)
			defer cancel()
			server.Shutdown(ctx)
		})
	}
	{
		ctx,cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			c:=make(chan os.Signal,1)
			signal.Notify(c,syscall.SIGINT,syscall.SIGTERM)

			select {
			case sig:= <-c:
				return fmt.Errorf("received %s",sig)
			case <-ctx.Done():
				return ctx.Err()

			}
		},func (error){
			cancel()
		})
	}

	g.Run()

}
