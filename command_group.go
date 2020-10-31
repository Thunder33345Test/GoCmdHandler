package CommandHandler

type CommandGroup interface {
	Register(cmd Command)
	Unregister(cmd Command) bool
	Lookup(cmd string) (Command, bool)
	LookupName(cmd string) (Command, bool)
	LookupAlias(cmd string) (Command, bool)
}

type Describable interface {
	GroupName() string
	GroupDescription() string
	GroupVisible() bool
}

type Transparent interface {
	Commands() map[string]Command
	Aliases() map[string]string
}

type CommandGroupBase struct {
	name        string
	description string
	commands    map[string]Command
	aliases     map[string]string
}

func NewGroup(Name string, Description string) CommandGroupBase {
	return CommandGroupBase{
		name:        Name,
		description: Description,
		commands:    map[string]Command{},
		aliases:     map[string]string{},
	}
}

func (c *CommandGroupBase) GroupName() string {
	return c.name
}

func (c *CommandGroupBase) GroupDescription() string {
	return c.description
}

func (c *CommandGroupBase) GroupVisible() bool {
	return true
}

func (c *CommandGroupBase) Commands() map[string]Command {
	return c.commands
}

func (c *CommandGroupBase) Aliases() map[string]string {
	return c.aliases
}

func (c *CommandGroupBase) Register(cmd Command) {
	c.commands[cmd.Name()] = cmd
	for _, al := range cmd.Alias() {
		c.aliases[al] = cmd.Name()
	}
}

func (c *CommandGroupBase) Unregister(cmd Command) bool {
	found := false
	for i, cs := range c.commands {
		if cs == cmd {
			delete(c.commands, i)
			found = true
		}
	}
	if found {
		for i, a := range c.aliases {
			if a == cmd.Name() {
				delete(c.aliases, i)
			}
		}
	}
	return found
}

func (c *CommandGroupBase) Lookup(cmd string) (Command, bool) {
	l, res := c.LookupName(cmd)
	if res {
		return l, res
	}
	l, res = c.LookupAlias(cmd)
	if res {
		return l, res
	}
	return nil, false
}

func (c *CommandGroupBase) LookupName(cmd string) (Command, bool) {
	route, found := c.commands[cmd]
	if found {
		return route, true
	}
	return nil, false
}

func (c *CommandGroupBase) LookupAlias(cmd string) (Command, bool) {
	n, found := c.aliases[cmd]
	if found {
		return c.LookupName(n)
	}
	return nil, false
}
