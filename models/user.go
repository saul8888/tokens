package models

//database
type User struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty"` //password is empty don't show
	Role     string `json:"role"`
}
