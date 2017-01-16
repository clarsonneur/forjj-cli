package kingpinMock

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
	"reflect"
)

type ArgClause struct {
	vtype     int // Value type requested.
	name      string
	help      string
	required  bool
	vdefault  []string
	envar     string
	set_value bool
}

func NewArg(name, help string) (f *ArgClause) {
	f = new(ArgClause)
	f.name = name
	f.help = help
	f.vdefault = make([]string, 0, 1)
	gotrace.Trace("Arg created : (%p)%#v", f, f)
	return f
}

func (a *ArgClause) Stringer() string {
	ret := fmt.Sprintf("Arg (%p):\n", a)
	ret += fmt.Sprintf("  name: '%s'\n", a.name)
	ret += fmt.Sprintf("  help: '%s'\n", a.help)
	ret += fmt.Sprintf("  vtype: '%d'\n", a.vtype)
	ret += fmt.Sprintf("  required: '%s'\n", a.required)
	ret += fmt.Sprintf("  vdefault: '%s'\n", a.vdefault)
	ret += fmt.Sprintf("  envar: '%s'\n", a.envar)
	ret += fmt.Sprintf("  set_value: '%s'", a.set_value)
	return ret
}

func (f *ArgClause) IsHelp(help string) bool {
	return (f.help == help)
}

func (a *ArgClause) String() *string {
	a.vtype = StringType
	return new(string)
}

func (a *ArgClause) GetType() string {
	switch {
	case a.vtype == BoolType:
		return "bool"
	case a.vtype == StringType:
		return "string"
	}
	return "any"
}

func (f *ArgClause) Bool() *bool {
	f.vtype = BoolType
	return new(bool)
}

func (f *ArgClause) IsBool() bool {
	return (f.vtype == BoolType)
}

func (f *ArgClause) Required() clier.ArgClauser {
	f.required = true
	return f
}

func (f *ArgClause) IsRequired() bool {
	return (f.required == true)
}

func (f *ArgClause) Default(p1 ...string) clier.ArgClauser {
	f.vdefault = p1
	return f
}
func (f *ArgClause) IsDefault(p1 ...string) bool {
	return reflect.DeepEqual(f.vdefault, p1)
}

func (f *ArgClause) Envar(p1 string) clier.ArgClauser {
	f.envar = p1
	return f
}

func (f *ArgClause) IsEnvar(p1 string) bool {
	return (f.envar == p1)
}

func (f *ArgClause) SetValue(_ clier.Valuer) clier.ArgClauser {
	f.set_value = true
	return f
}

func (f *ArgClause) IsSetValue(_ clier.Valuer) bool {
	return f.set_value
}
