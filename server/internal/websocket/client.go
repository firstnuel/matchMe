package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
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

// Client represents a WebSocket client connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID uuid.UUID
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex

	// Connection-specific data
	connectionID *uuid.UUID // For chat connections
	isActive     bool
	lastActivity time.Time
}

// NewClient creates a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn, userID uuid.UUID) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		hub:          hub,
		conn:         conn,
		send:         make(chan []byte, 256),
		userID:       userID,
		ctx:          ctx,
		cancel:       cancel,
		isActive:     true,
		lastActivity: time.Now(),
	}
}

// SetConnectionID sets the connection ID for chat clients
func (c *Client) SetConnectionID(connectionID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connectionID = &connectionID
}

// GetConnectionID gets the connection ID for chat clients
func (c *Client) GetConnectionID() *uuid.UUID {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connectionID
}

// UpdateActivity updates the client's last activity time
func (c *Client) UpdateActivity() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastActivity = time.Now()
}

// IsActive returns whether the client is active
func (c *Client) IsActive() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isActive
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
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
		select {
		case <-c.ctx.Done():
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("websocket error: %v", err)
				}
				return
			}

			c.UpdateActivity()

			// Handle incoming message
			if err := c.handleMessage(message); err != nil {
				log.Printf("error handling message: %v", err)
				c.sendError("Invalid message format")
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
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

			// Add queued messages to the current message.
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
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(message []byte) error {
	var wsMessage WebSocketMessage
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		return err
	}

	// Handle different message types
	switch wsMessage.Type {
	case EventPing:
		c.sendPong()
	case EventMessageTyping:
		c.handleTypingEvent(wsMessage.Data)
	default:
		log.Printf("Unhandled message type: %s", wsMessage.Type)
	}

	return nil
}

// handleTypingEvent handles typing indicator events
func (c *Client) handleTypingEvent(data interface{}) {
	// Parse typing data
	jsonData, _ := json.Marshal(data)
	var typingEvent TypingEvent
	if err := json.Unmarshal(jsonData, &typingEvent); err != nil {
		log.Printf("Error parsing typing event: %v", err)
		return
	}

	// Set the user ID from the client
	typingEvent.UserID = c.userID
	typingEvent.UpdatedAt = time.Now()

	// Broadcast to connection if client has one
	if c.GetConnectionID() != nil {
		c.hub.broadcastToConnection(*c.GetConnectionID(), EventMessageTyping, typingEvent, c.userID)
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(message string) {
	errorEvent := ErrorEvent{
		Code:    400,
		Message: message,
	}
	wsMessage := NewWebSocketMessage(EventError, errorEvent)
	if jsonData, err := wsMessage.ToJSON(); err == nil {
		select {
		case c.send <- jsonData:
		default:
			// Channel is full, skip
		}
	}
}

// sendPong sends a pong message to the client
func (c *Client) sendPong() {
	wsMessage := NewWebSocketMessage(EventPong, nil)
	if jsonData, err := wsMessage.ToJSON(); err == nil {
		select {
		case c.send <- jsonData:
		default:
			// Channel is full, skip
		}
	}
}

// SendMessage sends a WebSocket message to the client
func (c *Client) SendMessage(eventType EventType, data interface{}) {
	wsMessage := NewWebSocketMessage(eventType, data)
	if jsonData, err := wsMessage.ToJSON(); err == nil {
		select {
		case c.send <- jsonData:
		default:
			// Channel is full, close the client
			close(c.send)
		}
	}
}

// Close closes the client connection
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isActive {
		c.isActive = false
		c.cancel()
		close(c.send)
	}
}

// ServeWS handles websocket requests from the peer.
func ServeWS(hub *Hub, c *gin.Context, userID uuid.UUID) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := NewClient(hub, conn, userID)
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.writePump()
	go client.readPump()
}
