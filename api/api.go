package api

type Users struct {
	ID        string
	Username  string
	Email     string
	Firstname string
	Lastname  string
	Password  string
	Birthday  string
	ImageURL  string
	Role      string
	CreatedAt string
	UpdatedAt string
}

type Error struct {
	Code    int
	Message string
}
