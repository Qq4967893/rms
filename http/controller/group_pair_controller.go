package controller

import (
	"backend-rms/http/request"
	"backend-rms/http/response"
	"backend-rms/model"
	"backend-rms/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GroupPairController struct{}

func (gpc *GroupPairController) GetAll(c *gin.Context) {
	groupPairs, err := service.GroupPairService.GetAll()
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SucceedWithData(c, "", groupPairs)
}

func (gpc *GroupPairController) Add(c *gin.Context) {
	req := request.AddGroupPairRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	isAnyGroupExisted, isAnyGroupExistedErr := service.GroupPairService.IsAnyGroupExisted(req.GroupA, req.GroupB, 0)
	if isAnyGroupExisted {
		response.FailWithMessage(c, "One of groups exists already")
		return
	}
	if isAnyGroupExistedErr != nil {
		response.Error(c, isAnyGroupExistedErr)
		return
	}

	groupPair := &model.GroupPair{
		GroupA: req.GroupA,
		GroupB: req.GroupB,
	}
	if err := service.GroupPairService.Add(groupPair); err != nil {
		response.Error(c, err)
		return
	}

	response.Succeed(c)
}

func (gpc *GroupPairController) Update(c *gin.Context) {
	req := request.UpdateGroupPairRequest{}
	if ok := request.Validate(c, &req); !ok {
		return
	}

	groupPair, err := service.GroupPairService.GetInfoById(req.Id)
	if err != nil {
		response.Error(c, err)
		return
	}

	if groupPair == nil {
		response.FailWithMessage(c, "GroupPair doesn't exist")
		return
	}

	isAnyGroupExisted, isAnyGroupExistedErr := service.GroupPairService.IsAnyGroupExisted(req.GroupA, req.GroupB, req.Id)
	if isAnyGroupExisted {
		response.FailWithMessage(c, "One of groups exists already")
		return
	}
	if isAnyGroupExistedErr != nil {
		response.Error(c, isAnyGroupExistedErr)
		return
	}

	groupPair.GroupA = req.GroupA
	groupPair.GroupB = req.GroupB
	if err := service.GroupPairService.Update(groupPair); err != nil {
		response.Error(c, err)
		return
	}

	response.Succeed(c)
}

func (gpc *GroupPairController) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	rows, err := service.GroupPairService.Delete(id)
	if err != nil {
		response.Error(c, err)
		return
	}

	if rows < 1 {
		response.FailWithMessage(c, "GroupPair doesn't exist")
		return
	}

	response.Succeed(c)
}
