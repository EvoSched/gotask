package cobra

import (
	"fmt"
	"os"

	"github.com/EvoSched/gotask/internal/service"
)

// TODO: divide to cli and tui handlers
type Cmd struct {
	repo *service.TaskRepo
}

func NewCmd(repo *service.TaskRepo) *Cmd {
	return &Cmd{repo}
}

func (c *Cmd) Execute() {
	rootCmd := c.RootCmd()

	rootCmd.AddCommand(c.AddCmd(), c.ModCmd(), c.DeleteCmd(), c.GetCmd(), c.ListCmd(), c.DueCmd(), c.ArchivedCmd(),
		c.DoneCmd(), c.UndoCmd(), c.NoteCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
