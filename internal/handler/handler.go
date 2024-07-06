package handler

import (
	"fmt"
	"os"

	"github.com/EvoSched/gotask/internal/service"
)

// TODO: divide to cli and tui handlers
type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Execute() {
	rootCmd := h.RootCmd()

	rootCmd.AddCommand(h.AddCmd(), h.GetCmd(), h.ListCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
