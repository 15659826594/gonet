package controller

import (
	"gota/internal/common/controller"
	"gota/pkg/app/route"
	"gota/pkg/logger"
	"log/slog"
	"net/http"
	"time"

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
type Product struct {
	ID    int
	Name  string
	Price float64
}

func (t *Index) Index() (gin.HandlerFunc, []string, []string) {
	return func(c *gin.Context) {
		logger.RecordAttrs(c, slog.LevelInfo, "携带上下文",
			slog.String("method", "GET"),
			slog.String("path", "/api/users"),
			slog.Duration("duration", time.Since(time.Now())))

		product := Product{ID: 1, Name: "商品A", Price: 99.9}

		// 方式1: Any
		logger.RecordAttrs(c, slog.LevelInfo, "商品信息", slog.Any("product", product))
		slog.Info("111")
		// 方式2: Group
		logger.RecordAttrs(c, slog.LevelInfo, "商品信息",
			slog.Group("product",
				slog.Int("id", product.ID),
				slog.String("name", product.Name),
				slog.Float64("price", product.Price),
			),
		)
		t.View.Fetch(c)
	}, []string{"index", "/"}, []string{http.MethodGet}
}
