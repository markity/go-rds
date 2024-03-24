package commands

type PingCommand struct {
	Message *string
}

type HelloCommand struct {
	Proto int
}

type SetInfoLibNameCommand struct {
	LibName string
}

type SetInfoLibVersionCommand struct {
	LibVersion string
}

type CommandCommand struct {
}
