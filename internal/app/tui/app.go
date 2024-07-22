package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/EvoSched/gotask/internal/config"
	"github.com/EvoSched/gotask/internal/repository"
	"github.com/EvoSched/gotask/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ConfigDir = "configs"
)

type model struct {
	tasks []string
}

func Run() {
	//repository<-service<-handler<-tui

	//init config
	cfg, err := config.NewConfig(ConfigDir)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	//init sqlite
	db, err := repository.NewSQLite(&cfg.SQLite)

	//init repository
	r := repository.NewRepository(db)

	//init service
	s := service.NewService(r)

	fmt.Println(s)

	//open tui
	p := tea.NewProgram(model{
		tasks: []string{
			"Task 1: Write report",
			"Task 2: Clean the house",
			"Task 3: Read a book",
			"Task 4: Exercise",
		},
	})
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Your Tasks:\n\n"
	for i, task := range m.tasks {
		s += fmt.Sprintf("%d. %s\n", i+1, task)
	}
	s += "\nPress q to quit.\n"
	return s
}
