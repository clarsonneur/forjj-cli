package kingpinMock

import (
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
)

type ParseContext struct {
	cmds     []*CmdClause
	app      *Application
	Elements []interface{}
}

type ParseContextTester interface {
	GetContext() *ParseContext
}

// Following functions are implemented by clier.ParseContexter

func (p *ParseContext) GetFlagValue(f clier.FlagClauser) (string, bool) {
	var flag *FlagClause

	if v, ok := f.(*FlagClause); !ok {
		return "", false
	} else {
		flag = v
	}

	for _, element := range p.app.context.Elements {
		if f, ok := element.(*FlagClause); ok && f == flag {
			return f.context, true
		}
	}
	return "", false
}

func (p *ParseContext) GetArgValue(a clier.ArgClauser) (string, bool) {
	var arg *ArgClause

	if v, ok := a.(*ArgClause); !ok {
		return "", false
	} else {
		arg = v
	}

	for _, element := range p.app.context.Elements {
		if v, ok := element.(*ArgClause); ok && a == arg {
			return v.context, true
		}
	}
	return "", false
}

func (p *ParseContext) SelectedCommands() (res []clier.CmdClauser) {
	res = make([]clier.CmdClauser, 0, len(p.cmds))
	for _, cmd := range p.cmds {
		res = append(res, cmd)
	}
	return
}

// Following functions are specific to the Mock

func (a *Application) NewContext() *ParseContext {
	a.context = new(ParseContext)
	a.context.app = a
	return a.context
}

func (p *ParseContext) SetContext(p1 ...string) *ParseContext {
	p.cmds = make([]*CmdClause, 0, len(p1))
	if len(p1) == 0 {
		return p
	}

	var cmd *CmdClause
	if v, found := p.app.cmds[p1[0]]; !found {
		gotrace.Trace("Unable to find %s Command from Application layer.", p1[0])
		return nil
	} else {
		cmd = v
	}
	p.cmds = append(p.cmds, cmd)

	if len(p1) == 1 {
		return p
	}

	for _, cmd_name := range p1[1:] {
		if v, found := cmd.cmds[cmd_name]; !found {
			gotrace.Trace("Unable to find %s Command from Application layer.", cmd)
			return nil
		} else {
			p.cmds = append(p.cmds, v)
			p.Elements = append(p.Elements, v)
		}
	}
	return p
}

func (p *ParseContext) SetContextValue(name string, value string) *ParseContext {
	return p.setValue(true, false, name, value)
}

func (p *ParseContext) SetCliValue(name string, value string) *ParseContext {
	return p.setValue(false, true, name, value)
}

func (p *ParseContext) SetValue(name string, value string) *ParseContext {
	return p.setValue(true, true, name, value)
}

func (p *ParseContext) setValue(context, cli bool, name string, value string) *ParseContext {
	if p == nil {
		return nil
	}

	// App
	if len(p.cmds) == 0 {
		if v, found := p.app.flags[name]; found {
			if context {
				v.SetContextValue(value)
			}
			if cli {
				switch v.value.(type) {
				case *string:
					*v.value.(*string) = value
				case *bool:
					if value == "true" {
						*v.value.(*bool) = true
					} else {
						*v.value.(*bool) = false
					}
				}
			}
			p.Elements = append(p.Elements, v)
			return p
		}

		// Args
		if v, found := p.app.args[name]; found {
			if context {
				v.SetContextValue(value)
			}
			if cli {
				switch v.value.(type) {
				case *string:
					*v.value.(*string) = value
				case *bool:
					if value == "true" {
						*v.value.(*bool) = true
					} else {
						*v.value.(*bool) = false
					}
				}
			}
			p.Elements = append(p.Elements, v)
			return p
		}

		return nil
	}

	cmd := p.cmds[len(p.cmds)-1]

	// Flags
	if v, found := cmd.flags[name]; found {
		if context {
			v.SetContextValue(value)
		}
		if cli {
			switch v.value.(type) {
			case *string:
				*v.value.(*string) = value
			case *bool:
				if value == "true" {
					*v.value.(*bool) = true
				} else {
					*v.value.(*bool) = false
				}
			}
		}
		p.Elements = append(p.Elements, v)
		return p
	}

	// Args
	if v, found := cmd.args[name]; found {
		if context {
			v.SetContextValue(value)
		}
		if cli {
			switch v.value.(type) {
			case *string:
				*v.value.(*string) = value
			case *bool:
				if value == "true" {
					*v.value.(*bool) = true
				} else {
					*v.value.(*bool) = false
				}
			}
		}
		p.Elements = append(p.Elements, v)
	}
	return p
}

func (p *ParseContext) SetContextAppValue(name string, value string) *ParseContext {
	if p == nil {
		return nil
	}

	if v, found := p.app.flags[name]; found {
		switch v.value.(type) {
		case *string:
			*v.value.(*string) = value
		case *bool:
			if value == "true" {
				*v.value.(*bool) = true
			} else {
				*v.value.(*bool) = false
			}
		}
	}
	return p
}

func (p *ParseContext) GetContext() *ParseContext {
	return p
}
