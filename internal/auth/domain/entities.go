// internal/auth/domain/user.go
package domain

import (
	"fmt"
	"net/mail"
)

type Role string

const (
	RoleEmployee  Role = "employee"
	RoleModerator Role = "moderator"
)

func (r Role) Validate() error {
	if r != RoleModerator && r != RoleEmployee {
		return fmt.Errorf("unsupported role: %s", r)
	}
	return nil
}

func (r Role) String() string {
	return string(r)
}

type Email string

func (e Email) String() string {
	return string(e)
}

func (e Email) Validate() error {
	_, err := mail.ParseAddress(string(e))
	if err != nil {
		return err
	}
	return nil
}

type User struct {
	ID       string
	Email    Email
	Password string // password hash
	Role     Role
}
