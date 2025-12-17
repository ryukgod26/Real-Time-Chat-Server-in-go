package main

import (
	"log"
	"fmt"
	"encoding/json"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type model struct{
	viewport viewport.Model
	textinput textinput.Model
	conn *websocket.Conm
	err error
}

func newModel(conn *websocket.Conn) model{
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = "New Message..."
	ti.CharLimit = 196
	ti.Width = 20

	vp := viewport.New(30, 5)
	vp.SetContent("Welcome to the Messenger\n")
	
	return model{
		viewport: vp,
		textinput: ti,
		conn: conn,
	}
	
}

func (m model) Init() tea.Cmd{
	tea.Batch(textinput.blink, waitForIncomingMessage(m.conn))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var(
		vpCmd tea.Cmd
		tiCmd tea.Cmd
	)

	switch msg := msg.(type){
	case tea.KeyMsg:
		switch msg.Type{
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			content := m.textInput.Value()
			if content != ""{
				outMsg := Message{Username: "Test_User", Content: content, Time: time.Now()}
				finalMsg, err := json.Marshal(outMsg)
				if err != nil{
					log.Println("Tui Marshalling Error")
					m.err = err
					return m, tea.Quit
				}
				err = m.conn.WriteMessage(websocket.TextMessage, finalMsg)
				if err != nil{
					log.Println("Writing Message tui Error")
					m.err = err
					return m, tea.Quit
				}
				m.textInput.Reset()

			}
		}
	case Message:
		formattedMsg := fmt.Sprintf("%s, %s, %s", msg.Username, msg.Content, msg.Time)
		m.viewport.SetContent(m.viewport.View() + "\n" + formatted)
		m.viewport.GotoBottom()
		return m, waitForIncomingMessage(m.conn)
	case err:
		m.err = err
		return m, nil
	}
}

func (m model) View() string{
	if m.err != nil{
		return fmt.Sprintf("Error: %v", m.err)
	}
	return fmt.Sprintf("%s\n%s\n(Press Esc to Exit)", m.viewport.View(), m.textinput.View())
}

func waitForImcomingMessage(conn *websocket.Conn) tea.Cmd{
	return func() {
		_, bytes, err := conn.ReadMessage()
		if err != nil{
			return err
		}
		var msg Message
		json.Unmarshal(bytes, &msg)
		return msg
	}
}

func main() {
	u := url.Url(Schene: "ws", Host: "localhost:8800", Path: "/ws")
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial Error:", err)
	}
	defer conn.Close()

	p := tea.NewProgram(intialModel(conn))
	if _, err := p.Run(); err != nil{
		fmt.Printf("Error Occurred: %v", err)
		os.Exit(1)
	}
}
