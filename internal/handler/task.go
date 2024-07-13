package handler

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func (h *Handler) RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "task",
		Short: "Task manager",
	}

	return rootCmd
}

func (h *Handler) AddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Run: func(cmd *cobra.Command, args []string) {
			prompt := promptui.Prompt{
				Label: "Enter Task",
			}

			result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			//call service to create task

			fmt.Printf("Added task: %s\n", result)
		},
	}

	return addCmd
}

func (h *Handler) GetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a task by ID",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Please provide a task ID")
				return
			}

			//transform id to int
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Error occurred")
				log.Fatal(err)
				return
			}

			//call service to get task
			task, err := h.service.GetTask(id)
			if err != nil {
				fmt.Println("Error fetching task")
				log.Fatal(err)
				return
			}

			fmt.Println("Task: ", task)
		},
	}

	return getCmd
}

func (h *Handler) ListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := h.service.GetTasks()
			if err != nil {
				fmt.Println("Error fetching tasks")
				log.Fatal(err)
				return
			}

			fmt.Println("Tasks:")
			for i, task := range tasks {
				fmt.Printf("%d. %s\n", i+1, task.Title)
			}
		},
	}

	return listCmd
}

func parseTimeStamp(arg string) (*time.Time, *time.Time, error) {
	lexer := NewLexer(arg)
	tokens := lexer.Scan()

	parser := NewParser(tokens)
	start, end, err := parser.Parse()

	if err != nil {
		return start, end, err
	}

	return start, end, err
}