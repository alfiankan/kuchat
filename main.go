package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gojek/courier-go"
	"github.com/google/uuid"
)

type ChatMessage struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	CreatedAt time.Time `json:"created_at"`
	Data      string    `json:"data"`
}

var messages []ChatMessage = []ChatMessage{}

func recvMsg(chatClient *courier.Client, topic string) {

	cb := func(ctx context.Context, ps courier.PubSub, m *courier.Message) {
		msg := new(ChatMessage)
		if err := m.DecodePayload(msg); err != nil {
			log.Println(err)
		}
		messages = append(messages, *msg)
	}

	err := chatClient.Subscribe(context.Background(), topic, cb, courier.QOSTwo)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(1000 * time.Hour)
}

func sendMsg(chatClient *courier.Client, sessionID string, msg ChatMessage) error {
	publishTopic := fmt.Sprintf("/chats/%s/%s", sessionID, msg.From)

	err := chatClient.Publish(context.Background(), publishTopic, msg, courier.QOSTwo)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func main() {

	if os.Args[1] == "init" {
		fmt.Println("Session ID", uuid.NewString())
		os.Exit(0)
	}

	chatClient, err := courier.NewClient(
		courier.WithAddress("127.0.0.1", 1883),
		courier.WithUsername(os.Args[2]),
	)

	if err != nil {
		panic(err)
	}

	if err := chatClient.Start(); err != nil {
		log.Fatal("Failed connect to broker", err)
	}

	sessionID := os.Args[1]
	senderEmail := os.Args[2]
	destinationEmail := os.Args[3]

	m := initialModel()
	m.chatClient = chatClient
	m.senderEmail = senderEmail
	m.destinationEmail = destinationEmail
	m.sessionID = sessionID

	subscribeTopic := fmt.Sprintf("/chats/%s/%s", sessionID, destinationEmail)

	go func() {

		recvMsg(chatClient, subscribeTopic)

	}()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {

		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	viewport         viewport.Model
	textarea         textarea.Model
	err              error
	chatClient       *courier.Client
	sessionID        string
	senderEmail      string
	destinationEmail string
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(lipgloss.NewStyle().GetMaxWidth(), 10)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea: ta,
		viewport: vp,
		err:      nil,
	}
}

func (m model) Init() tea.Cmd {

	return tea.Batch(textarea.Blink, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:

			newMsg := ChatMessage{
				From:      m.senderEmail,
				To:        m.destinationEmail,
				CreatedAt: time.Now(),
				Data:      m.textarea.Value(),
			}
			err := sendMsg(m.chatClient, m.sessionID, newMsg)

			if err == nil {

				messages = append(messages, newMsg)

				m.textarea.Reset()
				m.viewport.GotoBottom()
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	content := ""
	for _, c := range messages {

		if c.From == m.senderEmail {
			content += fmt.Sprintf("\033[31m[You] %s\033[0m \n", c.Data)
		} else {
			content += fmt.Sprintf("\033[32m[%s] %s\033[0m \n", c.From, c.Data)
		}

	}
	m.viewport.SetContent(content)
	m.viewport.GotoBottom()

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
