package models

import (
	"time"
)

// UserSession représente une session utilisateur (token JWT)
// Table: user_sessions
type UserSession struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	TokenHash        string    `gorm:"type:varchar(255);not null;index" json:"-"`    // Hash du token JWT
	RefreshTokenHash string    `gorm:"type:varchar(255)" json:"-"`                   // Hash du refresh token (optionnel)
	ExpiresAt        time.Time `gorm:"not null;index" json:"expires_at"`             // Date d'expiration du token
	IPAddress        string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"` // Adresse IP (IPv4 ou IPv6)
	UserAgent        string    `gorm:"type:text" json:"user_agent,omitempty"`        // User-Agent du navigateur
	CreatedAt        time.Time `json:"created_at"`
	LastActivity     time.Time `json:"last_activity"` // Dernière activité

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName spécifie le nom de la table
func (UserSession) TableName() string {
	return "user_sessions"
}
