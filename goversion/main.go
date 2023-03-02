package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func pprint(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func initProgram() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func main() {
	client := NewClient("sk-xxxxxxxxxxxxxxxxxxxxxxxxxxx")

	reader := bufio.NewReader(os.Stdin)

	var messages []Message

	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		messages = append(messages, Message{
			Role:    "user",
			Content: line,
		})

		completionReq := ChatCompletionReq{
			Model:    "gpt-3.5-turbo",
			Messages: messages,
			Stream:   true,
		}

		var chatCompletion string
		client.ChatCompletion(completionReq, func(completion string, err error) {
			if err != nil {
				panic(err)
			}

			chatCompletion += completion
			fmt.Print(completion)
		})
		fmt.Println()

		messages = append(messages, Message{
			Role:    "system",
			Content: chatCompletion,
		})
	}
}
