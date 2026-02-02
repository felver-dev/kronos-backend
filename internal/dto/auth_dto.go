package dto

// LoginRequest représente la requête de connexion
// Les tags binding sont utilisés par Gin pour valider automatiquement les données
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"` // Email (obligatoire, format email)
	Password string `json:"password" binding:"required"`   // Mot de passe (obligatoire)
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

// RegisterRequest représente la requête d'inscription d'un utilisateur
type RegisterRequest struct {
	Username     string `json:"username" binding:"required,min=3"`       // Nom d'utilisateur (obligatoire, min 3 caractères)
	Email        string `json:"email" binding:"required,email"`          // Email (obligatoire, format email)
	Password     string `json:"password" binding:"required,min=6"`       // Mot de passe (obligatoire, min 6 caractères)
	FirstName    string `json:"first_name,omitempty"`                   // Prénom (optionnel)
	LastName     string `json:"last_name,omitempty"`                     // Nom (optionnel)
	FilialeID    *uint  `json:"filiale_id" binding:"required"`            // ID de la filiale (obligatoire)
}

// RegisterResponse représente la réponse après une inscription réussie
type RegisterResponse struct {
	Token        string  `json:"token"`                   // Token JWT d'accès
	RefreshToken string  `json:"refresh_token,omitempty"` // Token de rafraîchissement (optionnel)
	User         UserDTO `json:"user"`                    // Informations de l'utilisateur créé
}
