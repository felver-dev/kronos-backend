package dto

// LoginRequest représente la requête de connexion
// Les tags binding sont utilisés par Gin pour valider automatiquement les données
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // Nom d'utilisateur (obligatoire)
	Password string `json:"password" binding:"required"` // Mot de passe (obligatoire)
}

// LoginResponse représente la réponse après une connexion réussie
type LoginResponse struct {
	Token        string  `json:"token"`                   // Token JWT d'accès
	RefreshToken string  `json:"refresh_token,omitempty"` // Token de rafraîchissement (optionnel)
	User         UserDTO `json:"user"`                    // Informations de l'utilisateur connecté
}

// RefreshTokenRequest représente la requête pour renouveler un token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"` // Token de rafraîchissement (obligatoire)
}
