package main

import (
	"fmt"
	"log"
	"os"
	"git.tablet.sh/tablet/boba/demo"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
	// types "tablet/obsidian-publish/types"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	wo, ho, _ := term.GetSize(int(os.Stdout.Fd()))
	log.Printf("term ix [%dx%d]", wo, ho)

	p := tea.NewProgram(demo.Setup(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
