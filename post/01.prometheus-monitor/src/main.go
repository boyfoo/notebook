package src

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
)

// 计数器向量
var prodsVisit = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "jtthink_prods_visit",
	},
	[]string{"prod_id"},
)

func main() {
	r := gin.New()
	r.GET("/prods/visit", func(ctx *gin.Context) {
		pidStr := ctx.Query("pid")
		_, err := strconv.Atoi(pidStr)
		if err != nil {
			ctx.JSON(400, gin.H{"message": "error pid"})
		}
		prodsVisit.With(prometheus.Labels{
			"prod_id": pidStr,
		}).Inc()
		ctx.JSON(200, gin.H{"message": "OK"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.Run(":8080")
}
