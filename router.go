package figo

type Router struct {
	routerMapFunc map[string]func() *Controller
}

func NewRouter() *Router {
	re := &Router{}

	re.routerMapFunc = make(map[string]func() *Controller, 10)

	return re
}

func (this *Router) AddController(path string, createFunc func() *Controller) {

	this.routerMapFunc[path] = createFunc

}

func (this *Router) GetController(path string) *Controller {

	if createFunc, exists := this.routerMapFunc[path]; exists {
		return createFunc()
	}

	return nil
}
