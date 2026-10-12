package util

import (
	"crypto/rand"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func DataOrderBy(slice any, pagingSort string, pagingOrder string) error {
	if pagingOrder != "asc" && pagingOrder != "desc" {
		return errors.New("the order parameter must be desc or asc")
	}

	pagingSort = changeToHump(pagingSort)

	sv, err := sliceValue(slice)
	if err != nil {
		return err
	}
	t := sv.Type().Elem()

	//判断指针类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return errors.Errorf("data type %T not support, shuld be struct or pointer to struct", slice)
	}

	// 获取字段的详情信息
	sf, ok := t.FieldByName(pagingSort)
	if !ok {
		return errors.Errorf("pagingSort name %s not found", pagingSort)
	}
	// 根据字段的类型创建函数。
	var compare func(a, b reflect.Value) bool

	//判断类别
	switch sf.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if pagingOrder == "desc" {
			compare = func(a, b reflect.Value) bool { return a.Int() > b.Int() }
		} else {
			compare = func(a, b reflect.Value) bool { return a.Int() < b.Int() }
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if pagingOrder == "desc" {
			compare = func(a, b reflect.Value) bool { return a.Uint() > b.Uint() }
		} else {
			compare = func(a, b reflect.Value) bool { return a.Uint() < b.Uint() }
		}
	case reflect.Float32, reflect.Float64:
		if pagingOrder == "desc" {
			compare = func(a, b reflect.Value) bool { return a.Float() > b.Float() }
		} else {
			compare = func(a, b reflect.Value) bool { return a.Float() < b.Float() }
		}
	case reflect.String:
		if pagingOrder == "desc" {
			compare = func(a, b reflect.Value) bool { return a.String() > b.String() }
		} else {
			compare = func(a, b reflect.Value) bool { return a.String() < b.String() }
		}
	case reflect.Bool:
		if pagingOrder == "desc" {
			compare = func(a, b reflect.Value) bool { return a.Bool() && !b.Bool() }
		} else {
			compare = func(a, b reflect.Value) bool { return !a.Bool() && b.Bool() }
		}
	case reflect.Struct:
		if (sf.Type == reflect.TypeOf(decimal.Decimal{})) {
			if pagingOrder == "desc" {
				compare = func(a, b reflect.Value) bool {
					return a.Interface().(decimal.Decimal).GreaterThan(b.Interface().(decimal.Decimal))
				}
			} else {
				compare = func(a, b reflect.Value) bool {
					return a.Interface().(decimal.Decimal).LessThan(b.Interface().(decimal.Decimal))
				}
			}
		} else if (sf.Type == reflect.TypeOf(carbon.Carbon{})) {
			if pagingOrder == "desc" {
				compare = func(a, b reflect.Value) bool {
					return a.Interface().(carbon.Carbon).Gt(b.Interface().(carbon.Carbon))
				}
			} else {
				compare = func(a, b reflect.Value) bool {
					return a.Interface().(carbon.Carbon).Lt(b.Interface().(carbon.Carbon))
				}
			}
		} else {
			return errors.Errorf("pagingSort type struct %s not supported", sf.Type)
		}
	default:
		return errors.Errorf("pagingSort type %s not supported", sf.Type)
	}

	sort.Slice(slice, func(i, j int) bool {
		a := sv.Index(i)
		b := sv.Index(j)
		if t.Kind() == reflect.Ptr {
			a = a.Elem()
			b = b.Elem()
		}
		a = a.FieldByIndex(sf.Index)
		b = b.FieldByIndex(sf.Index)
		return compare(a, b)
	})

	return nil
}

func GetResultKey(offset, limit, total int) (resultKey map[string]int, err error) {
	resultKey = map[string]int{
		"startKey": 0,
		"endKey":   0,
	}

	if offset < 0 || limit <= 0 {
		return nil, errors.New("offset or limit must be greater than to 0")
	}

	if offset > total {
		return
	}

	if (offset + limit) > total {
		resultKey["endKey"] = total
		return
	}

	resultKey["startKey"] = offset
	resultKey["endKey"] = offset + limit
	return
}

func ErrWrapMessage() string {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return "Unknown file(0) unkonwn func"
	}
	return fmt.Sprintf("%s\n\t%s:%d", runtime.FuncForPC(pc).Name(), file, line)
}

func InSlice[T string | int | int64, S []T](value T, slice S) bool {
	for _, v := range slice {
		if value == v {
			return true
		}
	}
	return false
}

func SliceChunk[T int | int64 | string | float64](data []T, size int) ([][]T, error) {
	if size < 1 {
		return nil, errors.New("size: cannot be less than 1")
	}
	length := len(data)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var result [][]T
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		result = append(result, data[i*size:end])
		i++
	}
	return result, nil
}

func GetCurrentFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	funcInfo := runtime.FuncForPC(pc[0])
	return funcInfo.Name()
}

//生成随机字符串
func GenerateRandomString(length int) string {
	letters := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func sliceValue(slice any) (reflect.Value, error) {
	//获取反射值对象
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return reflect.Value{}, errors.Errorf("Invalid slice type, value of type %v", slice)
	}
	return v, nil
}

func changeToHump(field string) (humpField string) {
	//下划线转换为大驼峰方式
	//pl需要全部转换为大写
	humpField = strings.Replace(field, "_pl", "PL", -1)
	humpField = strings.Replace(humpField, "_", " ", -1)
	humpField = cases.Title(language.Und, cases.NoLower).String(humpField)
	humpField = strings.Replace(humpField, " ", "", -1)
	return
}
