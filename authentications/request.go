package authentications

type Request interface {
	StringToSign() (string, error)
	UserData() (string, error)
}
