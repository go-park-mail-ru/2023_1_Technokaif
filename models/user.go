package models

type Sex uint8

const (
	Male Sex = iota
	Female
	Other
)

// User implements information about app's users
type User struct {
	ID        uint   `json:"-"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Sex       Sex    `json:"sex"`
}
