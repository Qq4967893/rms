package middleware

import (
	myCasbin "backend-rms/core/casbin"
	"backend-rms/http/response"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

var allowedAPI = []map[string]string{
	{
		"api":    "^/api/auth/code$",
		"method": "GET",
	},
	{
		"api":    "^/api/auth/login$",
		"method": "POST",
	},
	{
		"api":    "^/api/managers/logged_details$",
		"method": "GET",
	},
	{
		"api":    "^/api/managers/password$",
		"method": "PUT",
	},
	{
		"api":    "^/api/auth/logout$",
		"method": "POST",
	},
	{
		"api":    "^/api/auth/refresh_token$",
		"method": "POST",
	},
	{
		"api":    "^/api/symbols/originals$",
		"method": "GET",
	},
	{
		"api":    "^/api/groups$",
		"method": "GET",
	},
	{
		"api":    "^/api/roles/[1-9]\\d*/display_menus$",
		"method": "GET",
	},
	{
		"api":    "^/api/permissions/[1-9]\\d*$",
		"method": "GET",
	},
}

func RbacHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.URL.Path
		act := c.Request.Method

		if allowed(uri, act) {
			return
		}

		loginManagerID := c.GetInt64("login_manager_id")
		success, _ := myCasbin.EnforceByManager(myCasbin.CASBIN_G_PREFIX+strconv.Itoa(int(loginManagerID)), uri, act)
		if success {
			c.Next()
		} else {
			response.Abort(c, response.StatusForbidden, response.CodeForbidden)
		}
	}
}

func allowed(uri, act string) bool {
	for _, v := range allowedAPI {
		matchBool, _ := regexp.MatchString(v["api"], uri)
		if matchBool && v["method"] == act {
			return true
		}
	}
	return false
}
