package model


import "time"


type User struct {
ID int `json:"id"`
Username string `json:"username"`
Password string `json:"-"`
Role string `json:"role"`
CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
}


type LoginRequest struct {
Username string `json:"username"`
Password string `json:"password"`
}


type LoginResponse struct {
Token string `json:"token"`
User *User `json:"user"`
}