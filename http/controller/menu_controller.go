package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"

	"github.com/gin-gonic/gin"
)

type MenuController struct{}

func (mec *MenuController) GetList(c *gin.Context) {
	result, err := service.MenuService.GetMenuList()
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", result)
}

func (mec *MenuController) Add(c *gin.Context) {
	req := request.AddMenuRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	menu := &model.Menu{
		Name:    req.Name,
		Pid:     req.Pid,
		Path:    req.Path,
		Display: req.Display,
		Sort:    req.Sort,
	}
	err := service.MenuService.Add(menu)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}

func (mec *MenuController) Update(c *gin.Context) {
	req := request.UpdateMenuRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	isMenuIdExisted, isMenuIdExistedErr := service.MenuService.IsIdExisted(req.Id)
	if isMenuIdExistedErr != nil {
		response.Error(c, isMenuIdExistedErr)
		return
	}

	if !isMenuIdExisted {
		response.FailWithMessage(c, "The Menu ID not exists")
		return
	}

	menu := map[string]interface{}{
		"name":    req.Name,
		"pid":     req.Pid,
		"path":    req.Path,
		"display": req.Display,
		"sort":    req.Sort,
	}

	err := service.MenuService.Update(req.Id, menu)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}

func (mec *MenuController) Delete(c *gin.Context) {
	req := request.DeleteMenuRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	err := service.MenuService.Delete(req.Id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Succeed(c)
}
