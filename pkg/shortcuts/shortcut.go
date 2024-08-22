package shortcuts

type Shortcut interface {
	GetName() string
	GetCommand() string
	GetText() string
	ParseText(string)
}
