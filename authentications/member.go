package authentications

type Member interface {
	UserData() (string, error)
}
