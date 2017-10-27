package appcontext

var (
	a *Context
)

// New is a Context creation function
func New() *Context {
	c := &Context{}
	return c
}
