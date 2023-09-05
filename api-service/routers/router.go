package routers

import (
	apiv1 "api-service/controllers/api/v1"

	"github.com/gin-gonic/gin"
)

func InitRoute(r gin.IRouter) {
	v1 := r.Group("/api/v1")
	{
		account := v1.Group("/accounts")
		{
			account.GET("", apiv1.GetAccounts)
			account.GET("/:username", apiv1.GetAccountByUsername)
			account.POST("", apiv1.CreateAccount)
			account.PUT("/:username", apiv1.UpdateAccount)
			account.DELETE("/:username", apiv1.DeleteAccount)
		}
	}
}
