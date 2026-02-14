package template

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Template struct {
	*gin.Context
	mu   sync.RWMutex
	keys map[any]any
}

func NewTemplate(c *gin.Context) *Template {
	tpl := &Template{
		Context: c,
	}
	tpl.Assign("U", c.GetString("url"))
	return tpl
}

func (t *Template) Assign(key string, value any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.keys == nil {
		t.keys = make(map[any]any)
	}
	t.keys[key] = value
}

func (t *Template) Display(name string) {
	cloned := make(map[any]any, len(t.keys)+1)
	for k, v := range t.keys {
		cloned[k] = v
	}
	cloned["Think"] = t.keys
	t.HTML(http.StatusOK, name, cloned)
}
