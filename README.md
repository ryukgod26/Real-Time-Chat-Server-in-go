# Real-Time Chat Server in Go

A full-featured, real-time WebSocket chat server built with Go, offering multiple client interfaces including a web-based UI and terminal-based UI (TUI) options. I also tried to add redis in this Project but that feature broke the app and there is still some bugs in the app due to that. But Currently this app is working as expected. Also Just a Small Part of the Readme is written with the help of AI (because I suck at writing readmes).

## 🚀 Features

- **Real-time WebSocket Communication**: Instant message delivery using WebSocket protocol
- **Multiple Client Interfaces**:
  - Modern web-based chat interface with responsive design
  - Terminal UI using `tview` library (interactive TUI mode)
  - Alternative terminal UI using `bubbletea` framework
- **Server/Client Architecture**: Choose to run as server or client from a unified interface
- **Automatic Reconnection**: Web client automatically reconnects on connection loss
- **User Presence**: Track active connections and user activity
- **Message Broadcasting**: Messages are broadcast to all connected clients
- **Timestamped Messages**: All messages include precise timestamps
- **JSON Message Format**: Structured message handling with JSON encoding

## 🏗️ Architecture

The application follows a hub-and-spoke architecture pattern:

1. **Hub**: Central message broker that manages all client connections
2. **Clients**: Individual WebSocket connections with read/write pumps
3. **Message Queue**: Channel-based message broadcasting system
4. **Registration/Unregistration**: Dynamic client management

```
┌─────────────────────────────────────────────┐
│                    Hub                      │
│  ┌─────────────────────────────────────┐   │
│  │ • Broadcast Channel                 │   │
│  │ • Client Registry Map               │   │
│  │ • Register/Unregister Channels      │   │
│  └─────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
                    │
        ┌───────────┼───────────┐
        ▼           ▼           ▼
    ┌────────┐  ┌────────┐  ┌────────┐
    │Client 1│  │Client 2│  │Client 3│
    └────────┘  └────────┘  └────────┘
```

## 📁 Project Structure

```
Real-Time-Chat-Server-in-go/
├── main.go              # Entry point with TUI menu
├── hub.go               # Hub implementation for client management
├── client.go            # WebSocket client handler with read/write pumps
├── message.go           # Message data structure
├── textviewer.go        # TView-based TUI implementation
├── test.html            # Web-based chat client interface
├── go.mod               # Go module dependencies
├── cmd/
│   └── tui/
│       └── main.go      # Bubbletea-based alternative TUI client
└── releases/            # Binary releases directory
```

## 🎮 Usage

### Running the Main Application

```bash
./go_chat
```

or directly with Go:

```bash
go run .
```

**Interactive Menu:**
1. **Run as Server**: Starts the WebSocket server on port 8800
2. **Run as Client**: Connect to an existing server via TUI client
3. **Quit**: Exit the application

## 📨 Message Format

Messages are exchanged in JSON format:

### Outgoing Message (Client → Server)

```json
{
  "username": "alice",
  "content": "Hello, world!",
  "time": "2026-01-27T10:30:00Z"
}
```

### Incoming Message (Server → Client)

```json
{
  "username": "bob",
  "content": "Hi there!",
  "time": "2026-01-27T10:30:15Z"
}
```

### Field Validation

- **username**: Required, non-empty string
- **content**: Required, non-empty string
- **time**: Automatically set by server to current time

## 🖥️ Client Interfaces

### 1. Web Client Features

- **Modern UI**: Gradient backgrounds and smooth animations
- **Message Types**:
  - User messages (red gradient, right-aligned)
  - Other users (white background, left-aligned)
  - System messages (yellow background, center-aligned)
- **Auto-reconnect**: Reconnects automatically after 3 seconds
- **Responsive**: Mobile-friendly design
- **Notifications**: System messages for connection status

### 2. TView Terminal Client Features

- **Interactive Menu**: Server/Client mode selection
- **Server Mode**: Real-time log viewer
- **Client Mode**: Full chat interface
- **Keyboard Navigation**: Esc, Ctrl+C shortcuts
- **Color Support**: ANSI color formatting

### 3. Bubbletea Terminal Client Features

- **Minimalist Design**: Clean, distraction-free interface
- **Viewport**: Scrollable message history
- **Text Input**: Character limit (196 chars)
- **Real-time Updates**: Immediate message display

## 🔍 Technical Details

### Concurrency Model

- **Goroutines per Client**: 2 (readPump + writePump)
- **Hub Goroutine**: 1 (main event loop)
- **Channel-based Communication**: Non-blocking message passing

### Connection Management

1. **Client Connection**:
   - HTTP upgrade to WebSocket
   - Client registration in Hub
   - Goroutines launched for read/write

2. **Message Flow**:
   - Client readPump receives message
   - Message validated and unmarshaled
   - Sent to Hub broadcast channel
   - Hub distributes to all client writePumps
   - writePumps send to respective connections

3. **Disconnection**:
   - Connection error detected
   - Unregister sent to Hub
   - Client removed from registry
   - Resources cleaned up

### Error Handling

- **Invalid JSON**: Logged and skipped, connection remains open
- **Empty username/content**: Logged and skipped
- **Connection errors**: Graceful cleanup and unregistration
- **Write timeouts**: Connection closed after writeWait period
- **Read timeouts**: Connection closed after pongWait period

### Security Considerations

- **Message size limit**: 512 bytes maximum
- **Input validation**: Username and content required
- **Connection timeouts**: Prevents resource exhaustion
- **WebSocket origin**: Currently allows all origins (upgrader default)

### Performance Features

- **Buffered Channels**: 256 message buffer per client
- **Ping/Pong Heartbeat**: Connection health monitoring
- **Non-blocking Broadcast**: Failed sends don't block others
- **Message Batching**: Multiple queued messages sent together

## 🛠️ Development

### Testing with Multiple Clients

1. Terminal 1: Start server
```bash
go run . 
# Select "Run as Server"
```

2. Terminal 2: Start TUI client
```bash
go run .
# Select "Run as Client"
```

3. Browser: Open `http://localhost:8800`
