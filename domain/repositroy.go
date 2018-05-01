package domain

//UserRepository ...
type UserRepository interface {
	StoreUser(user NewUser) (string, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByToken(token string) (*User, error)
	UpdateUser(user *User) (*User, error)
	CreateJWTToken(email string, password string) (*JwtToken, error)
}
