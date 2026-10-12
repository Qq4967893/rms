package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"

	"github.com/gin-gonic/gin"
)

type GroupController struct{}

func (gc *GroupController) GetAll(c *gin.Context) {
	groups, err := service.GroupService.GetAllList()
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", groups)
}

func (gc *GroupController) Add(c *gin.Context) {
	req := request.AddGroupRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	group := &model.Group{
		Name: req.Name,
		Book: req.Book,
	}
	if err := service.GroupService.Add(group); err != nil {
		response.Error(c, err)
		return
	}

	response.Succeed(c)
}

func (gc *GroupController) Update(c *gin.Context) {
	req := request.UpdateGroupRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	group, err := service.GroupService.GetInfoById(req.Id)
	if err != nil {
		response.Error(c, err)
		return
	}

	if group == nil {
		response.FailWithMessage(c, "Group doesn't exist")
		return
	}

	group.Name = req.Name
	group.Book = req.Book
	if err := service.GroupService.Update(group); err != nil {
		response.Error(c, err)
		return
	}

	response.Succeed(c)
}
