package repo

type (
	Connection struct {
		Name      string
		Password  string
		Addresses []string
	}
)

func NewConnection(
	name string,
	pass string,
	addrs ...string) *Connection {
	return &Connection{
		name, pass, addrs,
	}
}
