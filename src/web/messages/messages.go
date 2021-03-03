package messages

type EditGroup struct {
	Name string
}

type Group struct {
	Name    string
	General string
}

// RegisterMessage contains messages correlating to different parts of the register form
type RegisterMessage struct {
	Username string
	Email    string
	Password string
}
