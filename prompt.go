package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/codercola034/notion-ai/notion"
)

const helpMessage = `Commands:
  /? - show this help message
  /exit - exit the program
  /show - show all transcripts
`

func prompt() error {
	fmt.Println("Send a Message to Your Notion AI. Type /? for help.")
	for {
		fmt.Print("\033[1;34m" + ">>> " + "\033[0m")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		switch input {
		case "/?\n":
			fmt.Printf(helpMessage)
			continue
		case "/show\n":
			data, err := json.MarshalIndent(notion.DefaultTranscripts, "", "  ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", data)
			continue
		case "/exit\n":
			os.Exit(0)
		}

		res, err := notion.GetCompletion(input)
		if err != nil {
			return err
		}
		for out, ended, err := res.Output(); !ended; out, ended, err = res.Output() {
			if err != nil {
				return err
			}
			fmt.Printf(out)
		}
	}
}
