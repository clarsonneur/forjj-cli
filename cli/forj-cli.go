package cli

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
	"regexp"
	"strings"
)

// ForjCli is the Core cli for forjj command.
type ForjCli struct {
	App         clier.Applicationer        // *kingpin.Application       // Kingpin Application object
	flags       map[string]*ForjFlag       // Collection of Objects at Application level
	objects     map[string]*ForjObject     // Collection of Objects that forjj will manage.
	actions     map[string]*ForjAction     // Collection recognized actions
	list        map[string]*ForjObjectList // Collection of object list
	context     ForjCliContext             // Context from cli parsing
	values      map[string]*ForjRecords    // Collection of Object Values.
	filters     map[string]string          // List of field data identification from a list.
	sel_actions map[string]*ForjAction     // Selected actions
}

type ForjListParam interface {
	IsFound() bool
	GetAll() []ForjData
	IsList() bool
}

type ForjParamCopier interface {
	CopyToFlag(clier.CmdClauser) *ForjFlag
	CopyToArg(clier.CmdClauser) *ForjArg
}

type ForjParam interface {
	String() string
	IsFound() bool
	GetBoolValue() bool
	GetStringValue() string
	GetValue() interface{}
	Default(string) ForjParam
	set_cmd(clier.CmdClauser, string, string, string, *ForjOpts)
	loadFrom(clier.ParseContexter)
	IsList() bool
	CopyToFlag(clier.CmdClauser) *ForjFlag
	CopyToArg(clier.CmdClauser) *ForjArg
}

type forjParam interface {
	GetFlag() *ForjFlag
	GetArg() *ForjArg
}

// ForjParams type
const (
	// Arg : To set a ForjParam as Argument.
	Arg = "arg"
	// Flag : To set a ForjParam as Flag.
	Flag = "flag"
)

// List of ForjParams internal types
const (
	// String - Define the Param data type to string.
	String = "string"
	// Bool - Define the Param data type to bool.
	Bool = "bool"
	// List - Define a ForjObjectList data type.
	List = "list"
)

// NewForjCli : Initialize a new ForjCli object
//
// panic if app is nil.
func NewForjCli(app clier.Applicationer) (c *ForjCli) {
	if app.IsNil() {
		panic("kingpin.Application cannot be nil.")
	}
	c = new(ForjCli)
	c.objects = make(map[string]*ForjObject)
	c.actions = make(map[string]*ForjAction)
	c.flags = make(map[string]*ForjFlag)
	c.values = make(map[string]*ForjRecords)
	c.list = make(map[string]*ForjObjectList)
	c.filters = make(map[string]string)
	c.sel_actions = make(map[string]*ForjAction)
	c.App = app
	return
}

// AddFieldListCapture Add a Field list capture
func (c *ForjCli) AddFieldListCapture(key, capture string) error {
	if _, found := c.filters[key]; found {
		return fmt.Errorf("Key '%s' already exist.", key)
	}

	if _, err := regexp.Compile(capture); err != nil {
		return fmt.Errorf("Capture '%s' not created: Regexp error found: %s", capture, err)
	} else {
		parentheses_reg, _ := regexp.Compile(`\(`)
		if len(parentheses_reg.FindAllString(strings.Replace(capture, `\(`, "", -1), -1)) == 0 {
			capture = "(" + capture + ")"
		}
	}

	c.filters[key] = capture
	return nil
}

// AddAppFlag create a flag object at the application layer.
func (c *ForjCli) AddAppFlag(paramIntType, name, help string, options *ForjOpts) {
	f := new(ForjFlag)
	f.flag = c.App.Flag(name, help)
	f.set_options(options)

	switch paramIntType {
	case String:
		f.flagv = f.flag.String()
	case Bool:
		f.flagv = f.flag.Bool()
	}
	c.flags[name] = f
}

func (c *ForjCli) buildCapture(selector string) string {
	for key, capture := range c.filters {
		selector = strings.Replace(selector, "#"+key, capture, -1)
	}
	return strings.Replace(selector, "##", "#", -1)
}

// getValue : Core get value code for GetBoolValue and GetStringValue
func (c *ForjCli) getValue(object, key, param_name string) (interface{}, bool) {
	var value *ForjRecords

	if v, found := c.values[object]; !found {
		return nil, false
	} else {
		value = v
	}

	if v, found := value.Get(key, param_name); found {
		return v, true
	}
	return nil, false
}

// newParam create a ForjFlag or ForjArg defined by `paramType`
func (c *ForjCli) newParam(paramType, name string) ForjParam {
	switch paramType {
	case Arg:
		return new(ForjArg)
	case Flag:
		return new(ForjFlag)
	case List:
		l := new(ForjFlagList)
		if v, found := c.list[name]; found {
			l.obj = v
		} else {
			gotrace.Trace("Unable to find '%s' list.", name)
		}
		return l
	}
	return nil
}

// Create the ForjAction object to attach to the object parent.
func newForjObjectAction(ar *ForjAction, name, desc string) *ForjObjectAction {
	a := new(ForjObjectAction)
	a.action = ar
	a.name = ar.name + "_" + name
	a.cmd = ar.cmd.Command(name, fmt.Sprintf(ar.help, desc))
	a.params = make(map[string]ForjParam)
	a.plugins = make([]string, 0, 5)
	return a
}

func (c *ForjCli) getObject(obj_name string) (*ForjObject, error) {
	if v, found := c.objects[obj_name]; found {
		return v, nil
	}
	return nil, fmt.Errorf("Unable to find object '%s'", obj_name)
}

func (c *ForjCli) getObjectAction(obj_name, action string) (o *ForjObject, a *ForjObjectAction, err error) {
	err = nil
	if v, err := c.getObject(obj_name); err != nil {
		return nil, nil, err
	} else {
		o = v
	}

	if v, found := o.actions[action]; !found {
		return nil, nil, fmt.Errorf("Unable to find action '%s' from object '%s'", action, obj_name)
	} else {
		a = v
	}
	return
}

func (c *ForjCli) getObjectListAction(list_name, action string) (o *ForjObject, l *ForjObjectList, a *ForjObjectAction, err error) {
	err = nil
	if v, found := c.list[list_name]; !found {
		return nil, nil, nil, fmt.Errorf("Unable to find object list '%s'", list_name)
	} else {
		l = v
		o = l.obj
	}

	if v, found := o.actions[action]; !found {
		return nil, nil, nil, fmt.Errorf("Unable to find action '%s' from object '%s'", action, list_name)
	} else {
		a = v
	}
	return
}

func (c *ForjCli) getAction(action string) (a *ForjAction, err error) {
	err = nil
	if v, found := c.actions[action]; !found {
		return nil, fmt.Errorf("Unable to find action '%s'", action)
	} else {
		a = v
	}
	return
}
