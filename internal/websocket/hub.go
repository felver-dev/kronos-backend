package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client représente un client WebSocket connecté
type Client struct {
	ID       uint
	Conn     *websocket.Conn
	Send     chan []byte
	Hub      *Hub
	UserID   uint
	Username string
}

// Hub maintient l'ensemble des clients actifs et les messages de diffusion
type Hub struct {
	// Clients enregistrés
	clients map[*Client]bool

	// Canal pour enregistrer les clients
	register chan *Client

	// Canal pour désenregistrer les clients
	unregister chan *Client

	// Canal pour diffuser les messages à tous les clients
	broadcast chan []byte

	// Mutex pour la sécurité des threads
	mu sync.RWMutex
}

// NewHub crée une nouvelle instance de Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

// Run démarre le hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client WebSocket connecté: UserID=%d, Username=%s, Total clients=%d", client.UserID, client.Username, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("Client WebSocket déconnecté: UserID=%d, Username=%s, Total clients=%d", client.UserID, client.Username, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastNotification envoie une notification à tous les clients
func (h *Hub) BroadcastNotification(notification interface{}) {
	message, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Erreur lors de la sérialisation de la notification: %v", err)
		return
	}
	h.broadcast <- message
}

// SendToUser envoie un message à un utilisateur spécifique
func (h *Hub) SendToUser(userID uint, notification interface{}) {
	message, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Erreur lors de la sérialisation de la notification: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

// SendToUsers envoie un message à plusieurs utilisateurs
func (h *Hub) SendToUsers(userIDs []uint, notification interface{}) {
	message, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Erreur lors de la sérialisation de la notification: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	userMap := make(map[uint]bool)
	for _, id := range userIDs {
		userMap[id] = true
	}

	for client := range h.clients {
		if userMap[client.UserID] {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

// GetClientCount retourne le nombre de clients connectés
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
