package main

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "strings"
    "time"
    "github.com/rivo/tview"
    "github.com/gorilla/websocket"
	"github.com/gdamore/tcell/v2"
)


type textViewWriter struct {
    app *tview.Application
    tv  *tview.TextView
}

func (w *textViewWriter) Write(p []byte) (int, error) {
    s := string(p)
    w.app.QueueUpdateDraw(func() {
        fmt.Fprint(w.tv, s)
    })
    return len(p), nil
}

func buildMenu(app *tview.Application) *tview.Modal {
    modal := tview.NewModal().
        SetText("Select mode").
        AddButtons([]string{"Run as Server", "Run as Client", "Quit"}).
        SetDoneFunc(func(buttonIndex int, buttonLabel string) {
            switch buttonLabel {
            case "Run as Server":
                root := buildServerUI(app)
                app.SetRoot(root, true)
            case "Run as Client":
                root := buildClientForm(app)
                app.SetRoot(root, true)
            default:
                app.Stop()
            }
        })

    return modal
}

func buildServerUI(app *tview.Application) tview.Primitive {
    logView := tview.NewTextView().
        SetDynamicColors(true).
        SetScrollable(true).
        SetChangedFunc(func() { app.Draw() })
    logView.SetBorder(true).SetTitle(" Server Logs ")

    footer := tview.NewTextView().SetText("Press Esc to go back, Ctrl+C to exit").
        SetTextColor(tview.Styles.SecondaryTextColor)

    flex := tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(logView, 0, 1, false).
        AddItem(footer, 1, 0, false)

    log.SetOutput(&textViewWriter{app: app, tv: logView})

    go func() {
        hub := createHub()
        go hub.run()

        http.HandleFunc("/", serveHome)
        http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
            log.Printf("[WS] upgrade request from %s, path=%s", r.RemoteAddr, r.URL.Path)
            serveWs(hub, w, r)
        })

        log.Printf("[SERVER] Starting on %s ...\n", *addr)
        if err := http.ListenAndServe(*addr, nil); err != nil {
            log.Printf("[SERVER] ListenAndServe error: %v\n", err)
        }
    }()

    flex.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
        switch e.Key() {
        case tcell.KeyEsc:
            app.SetRoot(buildMenu(app), true)
            return nil
        }
        return e
    })

    return flex
}

type chatClient struct {
    app      *tview.Application
    msgView  *tview.TextView
    input    *tview.InputField
    conn     *websocket.Conn
    username string
}

func buildClientForm(app *tview.Application) tview.Primitive {
    form := tview.NewForm()
    serverURL := "ws://localhost:8800/ws"
    username := ""
    password := ""

    form.AddInputField("Server URL", serverURL, 60, nil, func(text string) { serverURL = strings.TrimSpace(text) })
    form.AddInputField("Username", "", 32, nil, func(text string) { username = strings.TrimSpace(text) })
    form.AddPasswordField("Password", "", 64, '*', func(text string) { password = text })
    form.AddButton("Connect", func() {
        if serverURL == "" || username == "" {
            modal := tview.NewModal().SetText("Server URL and Username are required").
                AddButtons([]string{"OK"}).
                SetDoneFunc(func(i int, l string) { app.SetRoot(form, true) })
            app.SetRoot(modal, true)
            return
        }
        root, err := buildChatUI(app, serverURL, username, password)
        if err != nil {
            modal := tview.NewModal().
                SetText(fmt.Sprintf("Failed to connect:\n\n%v", err)).
                AddButtons([]string{"Back"}).
                SetDoneFunc(func(i int, l string) { app.SetRoot(form, true) })
            app.SetRoot(modal, true)
            return
        }
        app.SetRoot(root, true)
    })

    form.AddButton("Back", func() {
        app.SetRoot(buildMenu(app), true)
    })

    form.SetBorder(true).SetTitle(" Connect as Client ").SetTitleAlign(tview.AlignLeft)

    return tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(form, 0, 1, true).
        AddItem(
            tview.NewTextView().SetText("Tip: example ws://localhost:8800/ws").
                SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false)
}

func buildChatUI(app *tview.Application, serverURL, username, _ string) (tview.Primitive, error) {

    u, err := url.Parse(serverURL)
    if err != nil || (u.Scheme != "ws" && u.Scheme != "wss") || u.Host == "" {
        return nil, fmt.Errorf("invalid WebSocket URL: %s", serverURL)
    }
    dialer := websocket.Dialer{
        HandshakeTimeout: 10 * time.Second,
    }
    conn, _, err := dialer.Dial(u.String(), nil)
    if err != nil {
        return nil, err
    }

    msgView := tview.NewTextView().
        SetDynamicColors(true).
        SetScrollable(true).
        SetChangedFunc(func() { app.Draw() })
    msgView.SetBorder(true).SetTitle(fmt.Sprintf(" Chat - %s ", username))

    input := tview.NewInputField().SetLabel("Message: ")
    status := tview.NewTextView().
        SetText(fmt.Sprintf("Connected to %s as %s | Esc: back | Ctrl+C: quit", serverURL, username)).
        SetTextColor(tview.Styles.SecondaryTextColor)

    client := &chatClient{
        app:      app,
        msgView:  msgView,
        input:    input,
        conn:     conn,
        username: username,
    }

    go func() {
        defer conn.Close()
        for {
            _, data, err := conn.ReadMessage()
            if err != nil {
                client.appendLine("[red::b]Disconnected[/]: " + err.Error())
                return
            }
            client.appendLine(string(data))
        }
    }()

    input.SetDoneFunc(func(key tcell.Key) {
        if key != tcell.KeyEnter {
            return
        }
        text := strings.TrimSpace(input.GetText())
        if text == "" {
            return
        }
        msg := fmt.Sprintf("%s: %s", username, text)
        err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
        if err != nil {
            client.appendLine("[red]send error:[-] " + err.Error())
            return
        }
        input.SetText("")
    })

    flex := tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(msgView, 0, 1, false).
        AddItem(input, 3, 0, true).
        AddItem(status, 1, 0, false)

    flex.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
        switch e.Key() {
        case tcell.KeyEsc:
            _ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
            _ = conn.Close()
            app.SetRoot(buildMenu(app), true)
            return nil
        }
        return e
    })

    _ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[joined] %s", username)))

    return flex, nil
}

func (c *chatClient) appendLine(s string) {
    c.app.QueueUpdateDraw(func() {
        fmt.Fprintln(c.msgView, s)
    })
}

// func main() {
//     flag.Parse()

//     app := tview.NewApplication()
//     root := buildMenu(app)
//     if err := app.SetRoot(root, true).EnableMouse(true).Run(); err != nil {
//         log.Fatal(err)
//     }
// }
