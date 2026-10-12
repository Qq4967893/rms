package market

import (
	"backend-rms/service"
	"backend-rms/util"
	"backend-rms/util/formatter"
	"backend-rms/ws"
	"backend-rms/ws/controller"
	"backend-rms/ws/request"
	"backend-rms/ws/request/market"
	"backend-rms/ws/response"
	"time"

	"github.com/pkg/errors"
)

type ExposureController struct{}

func (ec *ExposureController) GetExposureList(c *ws.Context) {
	controller.TickerHandle(c, 2*time.Second, func(c *ws.Context) (data any, err error) {
		req := market.GetExposureListRequest{}
		if ok := request.Validate(c, &req); !ok {
			return
		}

		return ec.buildExposureList(req.Params.Books, req.Params.Sort, req.Params.Order)
	})
}

func (ec *ExposureController) buildExposureList(books []string, sort, order string) (map[string][]response.SummaryExposure, error) {
	exposures, err := service.CacheService.GetBookExposures(books)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	//删除stock类型
	delete(exposures, "stock")

	notStockExposures := make(map[string][]response.SummaryExposure)
	for k, v := range exposures {
		err = util.DataOrderBy(v, sort, order)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}

		typeExposures := make([]response.SummaryExposure, 0)
		for _, vv := range v {
			currencyExposure := response.SummaryExposure{
				Currency:    vv.Currency,
				Exposure:    formatter.Amount(vv.Exposure),
				ExposureUSD: formatter.Amount(vv.ExposureUSD),
			}
			typeExposures = append(typeExposures, currencyExposure)
		}
		notStockExposures[k] = typeExposures
	}
	return notStockExposures, nil
}
