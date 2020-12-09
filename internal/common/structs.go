package common

// User represents a singe user.
// It contains his username and favourite number
type User struct {
	Username       string `json:"username"`
	FavoriteNumber int    `json:"favNumber"`
}

// UsersList contains multiple users and a count
type UsersList struct {
	Count int    `json:"count"`
	Users []User `json:"users"`
}

// Command is a message sent from the websocket client.
type Command struct {
	Action  string      `json:"action"`
	Message interface{} `json:"msg"`
}

// FavoriteNumber is used as the data field in Command in case a "setNumber" command is sent
type FavoriteNumber struct {
	Username string `json:"username"`
	Number   int    `json:"number"`
}
