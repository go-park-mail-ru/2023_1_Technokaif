package models

type Sex string

const (
	Male   Sex = "M"
	Female Sex = "F"
	Other  Sex = "O"
)

// User implements information about app's users
type User struct {
	ID        uint   `json:"-"`
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Sex       Sex    `json:"sex" binding:"required"`
}

func (u User) Validate() bool {
	if u.Username == "" ||
		u.FirstName == "" ||
		u.LastName == "" ||
		u.Password == "" ||
		u.Email == "" ||
		(u.Sex != Male && u.Sex != Female && u.Sex != Other) {
		return false
	}

	return true
}
