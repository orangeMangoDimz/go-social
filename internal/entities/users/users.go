package usersEntity

import "golang.org/x/crypto/bcrypt"

// User represents a user in the system
//
//	@Description	User account information
type User struct {
	ID        int64    `json:"id" example:"1"`                           // User ID
	Username  string   `json:"username" example:"johndoe"`               // Username
	Email     string   `json:"email" example:"johndoe@example.com"`      // Email address
	Password  Password `json:"-"`                                        // Password (never returned in responses)
	CreatedAt string   `json:"created_at" example:"2024-01-01 12:00:00"` // Account creation timestamp
	IsActive  bool     `json:"is_active"`
	RoleID    int64    `json:"role_id"`
	Role      Role     `json:"role"`
}

type Password struct {
	text *string
	Hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.Hash = hash

	return nil
}

func (p *Password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(text))
}
