package apiv1

import (
	"api-service/models"
	"api-service/services"
	"api-service/utils"

	"github.com/gin-gonic/gin"
)

func GetAccounts(ctx *gin.Context) {
	logger := utils.SugarLogger()
	svc := services.NewAccountService(ctx)
	accounts, err := svc.GetAccounts()
	if err != nil {
		logger.Errorf("get accounts error: %v", err)
		ctx.JSON(400, gin.H{"error": "get accounts error"})
		return
	}
	ctx.JSON(200, gin.H{"data": accounts})
}

func GetAccountByUsername(ctx *gin.Context) {
	logger := utils.SugarLogger()
	username := ctx.Param("username")
	svc := services.NewAccountService(ctx)
	account, err := svc.GetAccountByUsername(username)
	if err == services.ErrorAccountNotFound {
		logger.Errorf("get account error: %v", err)
		ctx.JSON(404, gin.H{"error": "account not found"})
		return
	}
	if err != nil {
		logger.Errorf("get account error: %v", err)
		ctx.JSON(400, gin.H{"error": "get account error"})
		return
	}
	ctx.JSON(200, gin.H{"data": account})
}

func CreateAccount(ctx *gin.Context) {
	logger := utils.SugarLogger()
	var account models.Account
	if err := ctx.BindJSON(&account); err != nil {
		logger.Errorf("create account error: %v", err)
		ctx.JSON(400, gin.H{"error": "BindJSON error"})
		return
	}
	svc := services.NewAccountService(ctx)
	if err := svc.CreateAccount(&account); err != nil {
		logger.Errorf("create account error: %v", err)
		ctx.JSON(401, gin.H{"error": "create account error"})
		return
	}
	ctx.JSON(200, gin.H{"data": account})
}

func UpdateAccount(ctx *gin.Context) {
	logger := utils.SugarLogger()
	var account models.Account
	if err := ctx.BindJSON(&account); err != nil {
		logger.Errorf("update account error: %v", err)
		ctx.JSON(400, gin.H{"error": "BindJSON error"})
		return
	}
	svc := services.NewAccountService(ctx)
	if err := svc.UpdateAccount(&account); err != nil {
		logger.Errorf("update account error: %v", err)
		ctx.JSON(401, gin.H{"error": "update account error"})
		return
	}
	ctx.JSON(200, gin.H{"data": account})
}

func DeleteAccount(ctx *gin.Context) {
	logger := utils.SugarLogger()
	var account models.Account
	if err := ctx.BindJSON(&account); err != nil {
		logger.Errorf("delete account error: %v", err)
		ctx.JSON(400, gin.H{"error": "BindJSON error"})
		return
	}
	svc := services.NewAccountService(ctx)
	if err := svc.DeleteAccount(&account); err != nil {
		logger.Errorf("delete account error: %v", err)
		ctx.JSON(401, gin.H{"error": "delete account error"})
		return
	}
	ctx.JSON(200, gin.H{"data": account})
}
