package commands

type PingCommand struct {
	Message *string
}

type EchoCommand struct {
	Message string
}

type HelloCommand struct {
	Proto string
}

type SetInfoLibNameCommand struct {
	LibName string
}

type SetInfoLibVersionCommand struct {
	LibVersion string
}

// redis client会发这个命令
type CommandCommand struct {
}

type UnknownCommand struct {
}
