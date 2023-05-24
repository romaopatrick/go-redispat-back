package router

import (
	"io"
	"log"
	"net/http"
	"redispat/red"
	"redispat/repo"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func UseRoutes(g *gin.Engine) *gin.Engine {
	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	v1 := g.Group("v1")
	{
		v1.GET(":conn/keys", GETKeys)
		v1.GET(":conn/keys/:k", GETValueByKey)
		v1.POST(":conn/keys/:k", POSTValueToKey)
		v1.DELETE(":conn/keys/:k", DELETEKey)
	}

	return g
}

func DELETEKey(ctx *gin.Context) {
	k := ctx.Param("k")
	cn := ctx.Param("conn")
	if k == "" || cn == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	c := repo.GetConnection(cn)
	if c == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("error: %v", r)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, r)
		}
	}()

	red.DeleteKey(ctx, c, k)
	ctx.Status(http.StatusNoContent)
}

func POSTValueToKey(ctx *gin.Context) {
	k := ctx.Param("k")
	cn := ctx.Param("conn")
	if k == "" || cn == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	c := repo.GetConnection(cn)
	if c == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("error: %v", r)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, r)
		}
	}()

	bdy, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}

	red.Set(ctx, c, k, string(bdy))
	ctx.JSON(http.StatusOK, string(bdy))
}

func GETValueByKey(ctx *gin.Context) {
	k := ctx.Param("k")
	cn := ctx.Param("conn")
	if k == "" || cn == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	c := repo.GetConnection(cn)
	if c == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("error: %v", r)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, r)
		}
	}()

	var val interface{}
	red.Get(ctx, c, k, &val)

	if val == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, val)
}

func GETKeys(ctx *gin.Context) {
	cn := ctx.Param("conn")
	if cn == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	contains := "*" + ctx.Query("contains") + "*"
	defer func() {
		if r := recover(); r != nil {
			log.Printf("error: %v", r)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, r)
		}
	}()

	c := repo.GetConnection(cn)
	if c == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ks := red.ListKeys(ctx, c, contains)
	if len(ks) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, ks)
}
