package commander

// CommandGroup is a collection of commands and middlewares that only apply to this path.
type CommandGroup struct {
	subgroups    []*CommandGroup
	commands     map[string]*command
	middleware   []CommandMiddleware
	Preprocessor Preprocessor
}

// CommandMiddleware is a simple middleware for commands.
type CommandMiddleware func(Handler) Handler

// New returns a new initialized Command Group
func New() *CommandGroup {
	return &CommandGroup{
		commands: make(map[string]*command),
	}
}

// Group creates a new subgroup on this command group path. This is useful when using middlewares, as they
// will only apply to this path.
func (g *CommandGroup) Group() *CommandGroup {
	ng := New()

	g.subgroups = append(g.subgroups, ng)

	return ng
}

// Use adds a new middleware to this command group path.
func (g *CommandGroup) Use(mw CommandMiddleware) {
	g.middleware = append(g.middleware, mw)
}
