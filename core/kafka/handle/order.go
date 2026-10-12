package handle

import (
	"backend-rms/core/logger"
	"backend-rms/model"
	"backend-rms/service"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/shopspring/decimal"
)

type OrderEvent struct {
	Ticket     int64           `json:"ticket"`
	Login      int64           `json:"login"`
	Volume     int64           `json:"volume"`
	Symbol     string          `json:"symbol"`
	Profit     decimal.Decimal `json:"profit"`
	AcType     int             `json:"actype"`  //开平仓 0开仓，1平仓
	OrdType    int             `json:"ordtype"` //订单类型 2部分平仓
	Cmd        int8            `json:"cmd"`
	OpenPrice  decimal.Decimal `json:"open_price"`
	ClosePrice decimal.Decimal `json:"close_price"`
	OpenTime   int64           `json:"open_time"`
	CloseTime  int64           `json:"close_time"`
	Comment    string          `json:"comment"`
	Sl         decimal.Decimal `json:"sl"`
	Tp         decimal.Decimal `json:"tp"`
	MarginRate decimal.Decimal `json:"margin_rate"`
}

func HandleOrderEvent(msg *sarama.ConsumerMessage) {
	var event OrderEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		logger.Errorf("Failed to parse order message: %v, raw: %s", err, string(msg.Value))
		return
	}

	// 判断是否是开平仓事件
	if event.AcType > 1 {
		return
	}
	logger.Infof("open or close event data:%v", event)

	// 判断订单用户是否在初始化的用户内，如果不在，查询所在组是否是A，B组，如果在，证明是新增用户，加入到初始化的用户内。如果不在直接return。
	if _, ok := service.InitializedData.Accounts[event.Login]; !ok {
		loginInfo, err := service.AccountService.GetInfoByLogin(event.Login)
		if err != nil {
			logger.Errorf("Order event get account info is err:%v", err)
			return
		}
		_, err = service.GroupService.GetBookByGroup(loginInfo.Group)
		if err != nil {
			logger.Infof("Order event login is not AB book:%v", err)
			return
		}

		service.InitializedData.Accounts[event.Login] = &model.Account{
			Login:    loginInfo.Login,
			Group:    loginInfo.Group,
			Currency: loginInfo.Currency,
			Balance:  loginInfo.Balance,
			Regdate:  loginInfo.Regdate,
		}
	}

	// 处理部分平仓的情况（例子：原订单是一个100volume的订单，部分平仓10volume；
	// 会发送2个订单事件：1发送10volume的平仓事件，发送一个90volume新开仓事件；而不会去发送原订单的平仓事件，这样就会导致根据订单事件获取品种手数出现问题；
	// 所以把发送10volume的平仓事件的volume改为100volume也就是原订单的总volume值）
	if event.OrdType == 2 {
		if _, ok := service.InitializedData.Trades[event.Ticket]; !ok {
			logger.Errorf("Order event 部分平仓原平仓订单不存在: %v", event.Ticket)
			return
		}
		event.Volume = service.InitializedData.Trades[event.Ticket].Volume
	}
	trade := &model.Trade{
		Ticket:     event.Ticket,
		Login:      event.Login,
		Cmd:        event.Cmd,
		Symbol:     event.Symbol,
		Volume:     event.Volume,
		OpenPrice:  event.OpenPrice,
		ClosePrice: event.ClosePrice,
	}
	switch event.AcType {
	case 0:
		service.SummaryService.HandleForOrder(trade, "open")
		service.ExposureService.HandleForOrder(trade, "open")

		// 订单事件的订单添加到初始化订单内
		service.InitializedData.Trades[event.Ticket] = trade
		break
	case 1:
		service.SummaryService.HandleForOrder(trade, "close")
		service.ExposureService.HandleForOrder(trade, "close")

		// 订单事件的订单从初始化订单内删除
		delete(service.InitializedData.Trades, trade.Ticket)
		break
	}
	// 把summary变化更新到redis
	book := service.InitializedData.Groups[service.InitializedData.Accounts[trade.Login].Group].Book
	service.CacheService.SetSummaryBook(book)
}
