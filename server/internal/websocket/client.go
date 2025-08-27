package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking for production
		return true
	},
}

// Client represents a single, generic WebSocket client connection.
// It is decoupled from any specific hub.
type Client struct {
	conn         *websocket.Conn
	send         chan []byte
	userID       uuid.UUID
	ctx          context.Context
	cancel       context.CancelFunc
	mu           sync.RWMutex
	connectionID *uuid.UUID // For chat/typing connections
	isActive     bool
	lastActivity time.Time
}

// NewClient creates a new WebSocket client without any hub reference.
func NewClient(conn *websocket.Conn, userID uuid.UUID) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		conn:         conn,
		send:         make(chan []byte, 256),
		userID:       userID,
		ctx:          ctx,
		cancel:       cancel,
		isActive:     true,
		lastActivity: time.Now(),
	}
}

// SetConnectionID sets the connection ID for chat or typing clients.
func (c *Client) SetConnectionID(connectionID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connectionID = &connectionID
}

// GetConnectionID gets the connection ID for chat or typing clients.
func (c *Client) GetConnectionID() *uuid.UUID {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connectionID
}

// UpdateActivity updates the client's last activity time.
func (c *Client) UpdateActivity() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastActivity = time.Now()
}

// IsStale checks if the client connection is stale (inactive for too long).
func (c *Client) IsStale() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return time.Since(c.lastActivity) > pongWait*2
}

// readPump pumps messages from the websocket connection.
// It accepts the specific hub's unregister channel and a reference to the typing hub
// (which will be nil for non-typing connections).
func (c *Client) readPump(unregister chan<- *Client, typingHub *TypingHub) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🚨 PANIC recovered in readPump for user %s: %v", c.userID, r)
			debug.PrintStack()
		}
		// Unregister from the specific hub that started this pump.
		unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.UpdateActivity()
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			return
		}
		c.UpdateActivity()
		if err := c.handleMessage(message, typingHub); err != nil {
			log.Printf("error handling message: %v", err)
		}
	}
}

// writePump pumps messages to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// handleMessage processes incoming WebSocket messages.
func (c *Client) handleMessage(message []byte, typingHub *TypingHub) error {
	var wsMessage WebSocketMessage
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		return err
	}

	switch wsMessage.Type {
	case EventMessageTyping:
		// Only handle typing events if a typingHub was provided.
		if typingHub != nil {
			c.handleTypingEvent(wsMessage.Data, typingHub)
		}
	// Other cases like EventPing can be added here if needed.
	default:
		log.Printf("Unhandled message type: %s", wsMessage.Type)
	}
	return nil
}

// handleTypingEvent handles typing indicator events.
func (c *Client) handleTypingEvent(data interface{}, typingHub *TypingHub) {
	var typingEvent TypingEvent
	jsonData, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonData, &typingEvent); err != nil {
		log.Printf("Error parsing typing event: %v", err)
		return
	}

	if connID := c.GetConnectionID(); connID != nil {
		typingHub.BroadcastTypingIndicator(*connID, typingEvent, c.userID)
	}
}

// SendMessage sends a WebSocket message to the client.
func (c *Client) SendMessage(eventType EventType, data interface{}) {
	c.mu.RLock()
	if !c.isActive {
		c.mu.RUnlock()
		return
	}
	c.mu.RUnlock()

	message := NewWebSocketMessage(eventType, data)
	messageBytes, err := message.ToJSON()
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}

	select {
	case c.send <- messageBytes:
		log.Printf("📨 Successfully queued message %s for user %s", eventType, c.userID)
	default:
		log.Printf("⚠️ Client send buffer full for user %s, closing connection", c.userID)
		c.Close()
	}
}

// IsActive returns whether the client is active.
func (c *Client) IsActive() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isActive
}

// Close closes the client connection.
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if !c.isActive {
		return
	}
	
	c.isActive = false
	c.cancel()
	close(c.send)
}

// --- HUB-SPECIFIC SERVE FUNCTIONS ---

// ServeStatusWS handles a generic WebSocket connection for the StatusHub.
func ServeStatusWS(hub *StatusHub, c *gin.Context, userID uuid.UUID) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	client := NewClient(conn, userID)
	hub.register <- client

	go client.writePump()
	// The readPump for a status client doesn't need to handle typing events.
	go client.readPump(hub.unregister, nil)
}

// ServeChatWS handles a chat-specific WebSocket connection.
func ServeChatWS(hub *ChatHub, c *gin.Context, userID uuid.UUID, connectionID uuid.UUID) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	client := NewClient(conn, userID)
	client.SetConnectionID(connectionID)
	hub.register <- client

	go client.writePump()
	// The readPump for a chat client doesn't need to handle typing events.
	go client.readPump(hub.unregister, nil)
}

// ServeTypingWS handles a typing-specific WebSocket connection.
func ServeTypingWS(hub *TypingHub, c *gin.Context, userID uuid.UUID, connectionID uuid.UUID) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	client := NewClient(conn, userID)
	client.SetConnectionID(connectionID)
	hub.register <- client

	go client.writePump()
	// The readPump for a typing client MUST be able to handle incoming typing events.
	go client.readPump(hub.unregister, hub)
}
