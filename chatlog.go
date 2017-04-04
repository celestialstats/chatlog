package chatlog

type ChatLog struct {
	Protocol string
	Server string
}

func NewChatLog(Protocol, Server string) *ChatLog {
	cl := new(ChatLog)
	cl.Protocol = Protocol
	cl.Server = Server
	return cl
}

func (chatLog *ChatLog) OpenLog() {
	
}

func (chatLog *ChatLog) AddEntry(Initiator, LineType, Content string) {
	
}

func (chatlog
