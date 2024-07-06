package cli

import (
	"log"

	"github.com/EvoSched/gotask/internal/config"
	"github.com/EvoSched/gotask/internal/handler"
	"github.com/EvoSched/gotask/internal/repository"
	"github.com/EvoSched/gotask/internal/service"

	"github.com/joho/godotenv"
)

const (
	ConfigDir  = "configs"
	ConfigFile = "main"
)

func Run() {
	//repository<-service<-handler<-cli

	//load env
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	//init config
	cfg, err := config.NewConfig(ConfigDir, ConfigFile)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	//init sqlite
	db, err := repository.NewSQLite(&cfg.SQLite)

	//init repository
	r := repository.NewRepository(db)

	//init service
	s := service.NewService(r)

	//init handler
	h := handler.NewHandler(s)

	//execute command
	h.Execute()
}
