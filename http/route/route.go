package route

import (
	"backend-rms/http/controller"
	"backend-rms/http/controller/report"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	g := r.Group("api")

	authGroup := g.Group("/auth")
	authController := new(controller.AuthController)
	{
		authGroup.GET("/code", authController.GetCode)
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/refresh_token", authController.RefreshToken)
		authGroup.POST("/logout", authController.Logout)
	}

	managerGroup := g.Group("/managers")
	//managerGroup.Use(middleware.RbacHandler())
	managerController := new(controller.ManagerController)
	{
		managerGroup.GET("/logged_details", managerController.GetLoggedDetails)
		managerGroup.GET("", managerController.GetList)
		managerGroup.POST("", managerController.Add)
		managerGroup.PUT("/:id", managerController.Update)
		managerGroup.DELETE("/:id", managerController.Delete)
		managerGroup.PUT("/:id/password", managerController.ResetPassword)
		managerGroup.PUT("/password", managerController.ModifyPassword)
	}

	hedgeAccountGroup := g.Group("/hedge_accounts")
	hedgeAccountController := new(controller.HedgeAccountController)
	{
		hedgeAccountGroup.GET("", hedgeAccountController.GetList)
		hedgeAccountGroup.PUT(":id", hedgeAccountController.Update)
	}

	roleGroup := g.Group("/roles")
	roleController := new(controller.RoleController)
	{
		roleGroup.GET("/:id/display_menus", roleController.GetDisplayMenus)
		roleGroup.GET("", roleController.GetList)
		roleGroup.POST("", roleController.Add)
		roleGroup.DELETE("/:id", roleController.Delete)
	}

	permissionGroup := g.Group("/permissions")
	permissionController := new(controller.PermissionController)
	{
		permissionGroup.GET("/:id", permissionController.GetInfo)
		permissionGroup.PUT("/:id", permissionController.Update)
	}

	MenuGroup := g.Group("/menus")
	menuController := new(controller.MenuController)
	{
		MenuGroup.GET("", menuController.GetList)
		MenuGroup.POST("", menuController.Add)
		MenuGroup.PUT("/:id", menuController.Update)
		MenuGroup.DELETE("/:id", menuController.Delete)
	}

	passportGroup := g.Group("/passports")
	passportController := new(controller.PassportController)
	{
		passportGroup.GET("/accounts", passportController.GetAccounts)
	}

	accountHedgingGroup := g.Group("/account_hedging")
	accountHedgingController := new(controller.AccountHedgingController)
	{
		accountHedgingGroup.GET("/confirm", accountHedgingController.GetConfirm)
		accountHedgingGroup.GET("/preview", accountHedgingController.GetPreview)
		accountHedgingGroup.POST("", accountHedgingController.Hedge)
	}

	passportDetailGroup := g.Group("passport_details")
	passportDetailController := new(controller.PassportDetailController)
	{
		passportDetailGroup.GET("/info", passportDetailController.GetInfo)
		passportDetailGroup.GET("/accounts", passportDetailController.GetAccountList)
		passportDetailGroup.GET("/trades", passportDetailController.GetTradeList)
		passportDetailGroup.GET("/symbols", passportDetailController.GetSymbolList)
	}

	changeGroupRecordGroup := g.Group("/change_group_records")
	changeGroupRecordController := new(controller.ChangeGroupRecordController)
	{
		changeGroupRecordGroup.GET("", changeGroupRecordController.GetList)
		changeGroupRecordGroup.GET("/operators", changeGroupRecordController.GetOperators)
	}

	hedgeTradeGroup := g.Group("/hedge_trades")
	hedgeTradeController := new(controller.HedgeTradeController)
	{
		hedgeTradeGroup.GET("", hedgeTradeController.GetList)
	}

	groupGroup := g.Group("/groups")
	groupController := new(controller.GroupController)
	{
		groupGroup.GET("", groupController.GetAll)
		groupGroup.POST("", groupController.Add)
		groupGroup.PUT("/:id", groupController.Update)
	}

	groupPairGroup := g.Group("/group_pairs")
	groupPairController := new(controller.GroupPairController)
	{
		groupPairGroup.GET("", groupPairController.GetAll)
		groupPairGroup.POST("", groupPairController.Add)
		groupPairGroup.PUT("/:id", groupPairController.Update)
		groupPairGroup.DELETE("/:id", groupPairController.Delete)
	}

	symbolGroup := g.Group("/symbols")
	symbolController := new(controller.SymbolController)
	{
		symbolGroup.GET("/originals", symbolController.GetOriginalList)
		symbolGroup.GET("/holding_accounts", symbolController.GetHoldingAccountList)
		symbolGroup.GET("/trades", symbolController.GetTradeList)
		symbolGroup.GET("/pending_trades", symbolController.GetPendingTradeList)
	}

	reportGroup := g.Group("/report")
	reportTradeController := new(report.TradeController)
	{
		reportGroup.GET("/trades", reportTradeController.GetTradeList)
	}

	reportSymbolController := new(report.SymbolController)
	{
		reportGroup.GET("/symbols", reportSymbolController.GetList)
		reportGroup.GET("/symbols/accounts", reportSymbolController.GetSymbolAccountList)
	}

	reportAccountController := new(report.AccountController)
	{
		reportGroup.GET("/accounts", reportAccountController.GetAccountList)
	}

	masterGroup := g.Group("/masters")
	masterController := new(controller.MasterController)
	{
		masterGroup.GET("", masterController.GetList)
		masterGroup.GET("/download", masterController.Download)

		masterHistoryController := new(controller.MasterHistoryController)
		{
			masterGroup.GET("/:masterId/histories", masterHistoryController.GetList)
		}

		followerController := new(controller.FollowerController)
		{
			masterGroup.GET("/:masterId/followers", followerController.GetList)
			masterGroup.GET("/:masterId/followers/:followerId/histories", followerController.GetHistoryList)
		}
	}

	reportIncomeController := new(report.IncomeController)
	{
		reportGroup.GET("/income/symbols", reportIncomeController.GetSymbolList)
		reportGroup.GET("/income/accounts", reportIncomeController.GetAccountList)
		reportGroup.GET("/income/trades", reportIncomeController.GetTradeList)
	}

	restrictedPassportsGroup := g.Group("/restricted_passports")
	restrictedPassportController := new(controller.RestrictedPassportController)
	{
		restrictedPassportsGroup.GET("", restrictedPassportController.GetList)
		restrictedPassportsGroup.GET("/records", restrictedPassportController.GetRecordList)
		restrictedPassportsGroup.PUT("/:id", restrictedPassportController.Remove)
	}
}
