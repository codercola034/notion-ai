package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/codercola034/notion-ai/notion"
	"github.com/codercola034/notion-ai/tui"
)

func init() {
	if err := notion.CheckToken(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	tuiOn := flag.Bool("tui", false, "run in terminal user interface mode")
	customPrompt := flag.String("prompt", "", "run with custom prompt")
	http := flag.Bool("http", false, "run in http mode")
	port := flag.Int("port", 8080, "port to run http server on")
	if *customPrompt != "" {
		notion.DefaultTranscripts.SetPrompt(*customPrompt)
	}
	flag.Parse()

	if *http {
		if err := httpServer(*port); err != nil {
			fmt.Println("Error running http server:", err)
			os.Exit(1)
		}
		return
	}

	if *tuiOn {
		p := tea.NewProgram(tui.NewModel())
		_, err := p.Run()
		if err != nil {
			fmt.Println("Error on running tui:", err)
			os.Exit(1)
		}
		return
	}

	// default prompt mode
	if err := prompt(); err != nil {
		fmt.Println("Error running default prompt:", err)
		os.Exit(1)
	}
	return
}
