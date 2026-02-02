package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Temps d'attente pour écrire un message au client
	writeWait = 10 * time.Second

	// Temps d'attente pour lire le prochain message pong du client
	pongWait = 60 * time.Second

	// Envoi de pings au client avec cette période (doit être < pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Taille maximale des messages en bytes
	maxMessageSize = 512 * 1024
)

// readPump pompe les messages du client WebSocket vers le hub
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erreur WebSocket: %v", err)
			}
			break
		}
	}
}

// writePump pompe les messages du hub vers le client WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Le hub a fermé le canal
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Ajouter les messages en attente dans le canal au message actuel
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs gère les requêtes WebSocket depuis le client
func ServeWs(hub *Hub, conn *websocket.Conn, userID uint, username string) {
	client := &Client{
		ID:       userID, // Utiliser UserID comme ID temporaire
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Hub:      hub,
		UserID:   userID,
		Username: username,
	}

	client.Hub.register <- client

	// Lancer les goroutines pour lire et écrire
	go client.writePump()
	go client.readPump()
}
