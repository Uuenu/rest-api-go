package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"rest-api-go/internal/config"
	"rest-api-go/internal/user"
	"rest-api-go/internal/user/db"
	"rest-api-go/pkg/client/mongodb"
	"rest-api-go/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	cfgMongodb := cfg.Mongodb
	logger.Infof("cfgMongodb.Username: %s", cfgMongodb.Username)

	mongoDbClient, err := mongodb.NewClient(context.Background(), cfgMongodb.Host, cfgMongodb.Port, cfgMongodb.Username,
		cfgMongodb.Password, cfgMongodb.Database, cfgMongodb.AuthDb)
	if err != nil {
		panic(err)
	}

	storage := db.NewStorage(mongoDbClient, cfgMongodb.Collection, logger)

	users, err := storage.FindAll(context.Background())
	fmt.Println(users)
	if err != nil {
		panic(err)
	}

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)

}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start aplication")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect ocp path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("Create socket")
		socketPaht := path.Join(appDir, "app.sock")
		logger.Debugf("socket path %s", socketPaht)

		logger.Info("Listen unix socket")
		listener, listenErr = net.Listen("unix", socketPaht)
		logger.Infof("server is listening unix socket %s", socketPaht)

	} else {
		logger.Info("Listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s%s", cfg.Listen.BindIP, cfg.Listen.Port)

	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))

}
