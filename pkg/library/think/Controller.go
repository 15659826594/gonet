package think

import (
	"gonet/pkg/library/traits/controller"
)

type Controller struct {
	*controller.Jump
	*View
}

//func (t *Controller) Context(c *gin.Context) *Controller {
//	if value, exists := c.Get("context"); exists {
//		if v, ok := value.(*Controller); ok {
//			return v
//		}
//	}
//	context := &Controller{}
//	c.Set("context", context)
//	return context
//}

//func (t *Controller) View(c *gin.Context) *View {
//	if value, exists := c.Get("view"); exists {
//		if v, ok := value.(*View); ok {
//			return v
//		}
//	}
//	view := &View{Context: c}
//	c.Set("view", view)
//	return view
//}
