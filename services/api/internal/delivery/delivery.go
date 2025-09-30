package delivery

import (
	"net/http"
	"research-apm/pkg/errors"
	"research-apm/pkg/errors/codes"
	"research-apm/pkg/ginx/response"
	"research-apm/pkg/tracer"
	"research-apm/services/api/internal/entity"
	"research-apm/services/api/internal/service"

	"github.com/gin-gonic/gin"
)

func NewDelivery(engine *gin.Engine, service service.Service) *http.Server {
	route := engine.Group("api/v1")
	route.GET("/user", GetUser(service))
	route.POST("/user", Create(service))
	route.GET("/message", GetMessage(service))
	route.GET("/client-do", GetClientDO(service))
	route.GET("/profil", GetProfil(service))
	return &http.Server{
		Handler: engine,
		Addr:    ":8080",
	}
}
func GetUser(service service.Service) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx, span := tracer.StartSpan(ginCtx.Request.Context(), "delivery.GetUser")
		defer span.End()
		result, err := service.GetUser(ctx)
		response.New(ginCtx, result, err)

	}
}
func Create(service service.Service) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx, span := tracer.StartSpan(ginCtx.Request.Context(), "delivery.Create")
		defer span.End()
		type Body struct {
			Name    string `json:"name" binding:"required"`
			Address string `json:"address" binding:"required"`
		}
		var body Body
		if err := ginCtx.ShouldBindJSON(&body); err != nil {
			response.New(ginCtx, nil, errors.New(codes.BadRequest, "payload tidak valid", err))
			return
		}
		result, err := service.CreateUser(ctx, entity.User{
			Name:    body.Name,
			Address: body.Address,
		})
		response.New(ginCtx, result, err)
	}
}

func GetMessage(service service.Service) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx, span := tracer.StartSpan(ginCtx.Request.Context(), "delivery.GetMessage")
		defer span.End()
		result, err := service.GetMessage(ctx)
		response.New(ginCtx, result, err)

	}
}

func GetClientDO(service service.Service) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx, span := tracer.StartSpan(ginCtx.Request.Context(), "delivery.GetClientDO")
		defer span.End()
		result, err := service.GetClientDO(ctx)
		response.New(ginCtx, result, err)

	}
}
func GetProfil(service service.Service) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx, span := tracer.StartSpan(ginCtx.Request.Context(), "delivery.GetProfil")
		defer span.End()
		result, err := service.GetProfil(ctx)
		response.New(ginCtx, result, err)

	}
}
