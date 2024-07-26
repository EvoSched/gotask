package ui

import (
	"fmt"
	"github.com/EvoSched/gotask/internal/models"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
)

func (m model) View() string {
	switch m.state {
	case Main:
		// todo we need to include a banner 'Main'
		return renderTable(m.main)
	case Pending:
		// todo we need to include a banner 'Pending'
		return renderTable(m.pend)
	case Archived:
		// todo we need to include a banner 'Archive'
		return renderTable(m.arch)
	case Detail:
		return renderDetail(m.task)
	case List:
		return renderList(m.list)
	case Help:
		return renderHelp()
	}
	return ""
}

func renderList(list list.Model) string {
	return listStyle.Render(list.View())
}

func renderDetail(t models.Task) string {
	return ""
}

func renderTable(tasks table.Model) string {
	return tasks.View() + "\n"
}

func renderHelp() string {
	return fmt.Sprintf(`
PAGE COMMANDS:
tab    -   Switches page
1      -   Switches to table page
2      -   Switches to list page
h      -   Switches to help page

TABLE COMMANDS:
t      -   Displays all tasks both new and archived
p      -   Displays all pending and overdue tasks
f      -   Displays all archived tasks
a      -   Adds a new task with all attributes available
d      -   Marks currently selected task complete
u      -   Marks currently selected task incomplete
m      -   Modifies currently selected task
r      -   Removes currently selected task
n      -   Adds note to currently selected task
enter  -   Switches over to task detail from selected task
s      -   Sorts table either ascending or descending order

TASK DETAIL COMMANDS:
m      -   Switches view mode to edit
tab    -   Switches entries in edit mode
q      -   Exits out of task detail page

LIST COMMANDS:
enter -   Selects available command from list
`)
}
