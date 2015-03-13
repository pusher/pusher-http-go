package pusher

type Users struct {
	List []User `json:"users"`
}

type User struct {
	Id string `json:"id"`
}
