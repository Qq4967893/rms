package validation

import (
	"backend-rms/core/database"
	"backend-rms/core/logger"
	"reflect"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"

	"github.com/go-playground/validator/v10"
)

func ExistsValidation(fl validator.FieldLevel) bool {
	sli := strings.Split(fl.Param(), "^")
	dbName := sli[0]
	table := sli[1]
	tableField := sli[2]

	var exceptValue int64
	if len(sli) > 3 {
		exceptField := sli[3]
		exceptValue, _ = reflect.Indirect(fl.Top()).FieldByName(exceptField).Interface().(int64)
	}

	isExisted, err := isRecordExisted(dbName, table, tableField, fl.Field().String(), exceptValue)
	if err != nil {
		logger.ErrorIf(err)
	}

	return isExisted
}

func NotExistsValidation(fl validator.FieldLevel) bool {
	isExisted := ExistsValidation(fl)
	return !isExisted
}

func TickQueryableTime(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Struct && field.Type() == reflect.TypeOf(time.Time{}) {
		startTime := carbon.Now().SubMonths(6).StartOfDay().Carbon2Time()
		endTime := time.Now()
		t := field.Interface().(time.Time)
		return (t.After(startTime) || t.Equal(startTime)) && (t.Before(endTime) || t.Equal(endTime))
	}
	return true
}

func isRecordExisted(dbName, table, field string, value interface{}, exceptID int64) (bool, error) {
	db, err := database.GetDB(dbName)
	if err != nil {
		return false, err
	}

	tx := db.Table(table).Where(field+" = ?", value)
	if exceptID != 0 {
		tx.Where("id != ?", exceptID)
	}

	var count int64
	result := tx.Count(&count)

	return count > 0, result.Error
}
