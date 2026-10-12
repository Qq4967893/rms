package download

import (
	"backend-rms/core/logger"
	"backend-rms/util"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

func HandleHttpDownload(c *gin.Context, data [][]string, fileName string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logger.Errorf("close file failed: %v", err)
		}
	}()

	// 创建一个工作表（1个sheet已满足该项目的需求）
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}

	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)

	for key, value := range data {
		cell, _ := excelize.CoordinatesToCellName(1, key+1)
		err := f.SetSheetRow("Sheet1", cell, &value)
		if err != nil {
			return errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName+".xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	err = f.Write(c.Writer)
	if err != nil {
		return errors.WithMessage(err, util.ErrWrapMessage())
	}
	return nil
}
