package bootstrap

import (
	"context"
	"fmt"
	Config "gin-laravel/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func initServer() {
	fmt.Printf("运行gin 模式%s http://%s \r\n", Config.App.GetString("server.GinMode"), Config.App.GetString("server.Address"))

	//设置gin模式
	gin.SetMode(Config.App.GetString("server.GinMode"))
	gin.DefaultWriter = ioutil.Discard

	srv := &http.Server{
		Addr:    Config.App.GetString("server.Address"),
		Handler: Engine,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
