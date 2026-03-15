package controller

import (
	"gonet/internal/common/controller"
	"gonet/pkg/app/route"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	route.Register(&Index{
		NoNeedLogin: []string{"*"},
		NoNeedRight: []string{"*"},
	})
}

type Index struct {
	controller.Frontend
	NoNeedLogin []string
	NoNeedRight []string
}

func (t *Index) Index() (gin.HandlerFunc, []string, []string) {
	return func(c *gin.Context) {
		//siteViper, _ := Config.Get("site")
		//t.View.Fetch()
		//c.HTML(http.StatusOK, "index/view/index/index.html", gin.H{
		//	"site": siteViper.AllSettings(),
		//	"U":    c.GetString("url"),
		//})
		t.View.Fetch(c)
	}, []string{"index", "/"}, []string{http.MethodGet}
}
