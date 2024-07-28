package cli

import (
	"github.com/EvoSched/gotask/internal/sqlite"
	"log"

	"github.com/EvoSched/gotask/internal/cobra"
	"github.com/EvoSched/gotask/internal/config"
	"github.com/EvoSched/gotask/internal/service"
)

const (
	ConfigDir = "configs"
)

func Run() {

	//init config
	cfg, err := config.NewConfig(ConfigDir)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	//init sqlite
	db, err := sqlite.NewSQLite(&cfg.SQLite)

	//init service
	r := service.NewTaskRepo(db)

	//init cobra
	h := cobra.NewCmd(r)

	//execute command
	h.Execute()
}
