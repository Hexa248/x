package models

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleStaff Role = "staff"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Role     Role
}
