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
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main(){
	nsa:=api.NewNameServiceApi()
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	}
	level.Info(logger).Log("starting","on")


	appLogger := log.With(logger,"app")
	var g run.Group
	{
		server:=http.Server{
			Addr:"127.0.0.1:25000",
			Handler:nsa,
		}

		g.Add(func() error {
			level.Info(appLogger).Log("server","127.0.0.1:25000")
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
