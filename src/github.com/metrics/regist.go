package metrics

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"

	fastp "github.com/DroiTaipei/fasthttp-prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func FasthttpPrometheusRegist(router *fasthttprouter.Router, subsystem string) fasthttp.RequestHandler {
	p := fastp.NewPrometheus(subsystem)
	return p.WrapHandler(router)
}

func GinPrometheusRegist(router *gin.Engine) {
	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)
}
