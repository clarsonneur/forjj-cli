package clier

// Define interface against kingpin or kingpinMock (for go test)

type Applicationer interface {
	Flag(string, string) FlagClauser
	Arg(string, string) ArgClauser
	Command(string, string) CmdClauser
}

type FlagClauser interface {
	String() *string
	Bool() *bool
	Required() FlagClauser
	Short(rune) FlagClauser
	Hidden() FlagClauser
	Default(...string) FlagClauser
	Envar(string) FlagClauser
	SetValue(Valuer) FlagClauser
}

type ArgClauser interface {
	String() *string
	Bool() *bool
	Required() ArgClauser
	Default(...string) ArgClauser
	SetValue(Valuer) ArgClauser
	Envar(string) ArgClauser
}

type CmdClauser interface {
	Command(string, string) CmdClauser
	Flag(string, string) FlagClauser
	Arg(string, string) ArgClauser
}

type ParseContexter interface {
	GetFlagValue(FlagClauser) (string, bool)
	GetArgValue(ArgClauser) (string, bool)
	SelectedCommand() CmdClauser
}

type Valuer interface {
	Set(string) error
	String() string
}
