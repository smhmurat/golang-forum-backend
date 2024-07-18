package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type AuthProvider struct {
	ID         int    `json:"id"`
	UserID     string `json:"user_id"`
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`
}
