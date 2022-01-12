package entity

type Error struct {
	Err string `json:"err"`
}

var (
	ErrInvalidInput     = Error{Err: "Invalid Input"}
	ErrBadSession       = Error{Err: "Bad Session"}
	ErrUserCreateFailed = Error{Err: "Could not create user"}
	ErrUserExists       = Error{Err: "A user with that name already exists"}
	ErrUserDoesNotExist = Error{Err: "User does not exist, you must register first"}
	ErrDatabaseError    = Error{Err: "Database Error"}
	ErrForbidden        = Error{Err: "Forbidden"}
	ErrUnauthorised     = Error{Err: "Unauthorised"}
	ErrInvalidLobby     = Error{Err: "Invalid Lobby ID"}
	ErrLobbyNotFound    = Error{Err: "Lobby Not Found"}
	ErrLoggedIn         = Error{Err: "You must be logged out to use this endpoint"}
)
