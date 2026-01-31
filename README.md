# Real-Time Chat Server in Go

A full-featured, real-time WebSocket chat server built with Go, offering multiple client interfaces including a web-based UI and terminal-based UI (TUI) options. I also tried to add redis in this Project but that feature broke the app and there is still some bugs in the app due to that. But Currently this app is working as expected. Also Just a Small Part of the Readme is written with the help of AI (because I suck at writing readmes).

## Features

- Real-time WebSocket Communication: Instant message delivery using WebSocket protocol
- Multiple Client Interfaces:
  - Modern web-based chat interface with responsive design
  - Terminal UI using `tview` library (interactive TUI mode)
  - Alternative terminal UI using `bubbletea` framework
- Server/Client Architecture: Option to run as server or client (You can use rasberry pi as a server)
- Automatic Reconnection: Web client automatically reconnects on connection loss (You Don't have to refresh the Page to check server status)
- User Presence: Track active connections and user activity (I am not stealing your data btw)
- Message Broadcasting: Messages are broadcast to all connected clients (All the people on thge same servercan read your messages )
- Timestamped Messages: All messages include precise timestamps (So You Know How much Time has it been Since You Suck at coding)
- JSON Message Format: Structured message handling with JSON encoding (Atleast messages have structure unlike your life)

## Architecture

The application follows a hub-and-spoke architecture pattern:

1. Hub: Central message broker that manages all client connections (The Machine which will handle message heandling for all the clients)
2. Clients: Individual WebSocket connections with read/write pumps 
3. Message Queue: Channel-based message broadcasting system
4. Registration/Unregistration: Dynamic client management

## Usage
- Running the Interface to select server/client

```bash
./go_chat
```

or directly Build the Project with Go:

```bash
go run .
```

**Interactive Menu:**
1. Run as Server: Starts the WebSocket server on port 8800
2. Run as Client: Connect to an existing server via TUI client
3. Quit: Exit the application

## Message Format

I am using Json Format for sending and receiving messages:

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

- username: Required, non-empty string
- content: Required, non-empty string
- time: Automatically set by server to current time

## Client Interfaces

### 1. Web Client Features

- Modern UI: Gradient backgrounds and smooth animations
- Message Types:
  - User messages (red gradient, right-aligned)
  - Other users (white background, left-aligned)
  - System messages (yellow background, center-aligned)
- Auto-reconnect: Reconnects automatically after 3 seconds
- Responsive: Mobile-friendly design (Not Fully)
- Notifications: System messages for connection status like it will show you when you get disconnected from the server

### 2. TView Terminal Client Features

- Interactive Menu: Server/Client mode selection
- Server Mode: Real-time log viewer (to fing bugs)
- Client Mode: Full chat interface (You can chat with web interface in this way)
- Keyboard Navigation: Esc, Ctrl+C shortcuts (Will add more shortcuts in future)
- Color Support**: ANSI color formatting

### 3. Bubbletea Terminal Client Features

- Minimalist Design: Clean, distraction-free interface (to help you lock in )
- Viewport: Scrollable message history (So You can check why they think you sucks)
- Text Input: Character limit (196 chars) (So You don't just Write a essay in one message)
- Real-time Updates: Immediate message display (To Give You Reports about bugs asap)

## Technical Details

### Concurrency Model

- Goroutines per Client: 2 (readPump + writePump)
- Hub Goroutine: 1 (main event loop)
- Channel-based Communication: Non-blocking message passing (Using goroutines for concurrency)

### Performance Features

- Buffered Channels: 256 message buffer per client
- Ping/Pong Heartbeat: Connection health monitoring
- Non-blocking Broadcast: Failed sends don't block others
- Message Batching: Multiple queued messages sent together

## Development

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
or Terminal: Enter Your Server address and Then Press Connect