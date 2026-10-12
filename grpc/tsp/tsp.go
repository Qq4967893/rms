package tsp

import (
	"backend-rms/grpc"
	pb "backend-rms/grpc/et6api"
	"backend-rms/repo"
	"backend-rms/util"
	"context"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetAccountInfo(ctx context.Context, requestID string, login int64) (*pb.GetAccountInfoResponse_ResponseData, error) {
	res, err := grpc.TSPCli.GetAccountInfo(ctx, &pb.LoginRequest{
		Login: int32(login),
		CReq:  buildCReq(requestID),
	})

	if err != nil {
		return nil, err
	}

	return res.GetData(), err
}

func GetUserEquity(ctx context.Context, requestID string, logins []int64) ([]*pb.EquityInfo, error) {
	allData := make([]*pb.EquityInfo, 0)
	//每次访问tsp不能超过50个login
	sliceLogins, err := util.SliceChunk(logins, 50)
	if err != nil {
		return nil, err
	}
	for _, vLogins := range sliceLogins {
		ids := make([]string, len(logins))
		for i, v := range vLogins {
			ids[i] = strconv.FormatInt(v, 10)
		}

		res, err := grpc.TSPCli.GetUserEquity(ctx, &pb.LoginsRequest{
			Logins: strings.Join(ids, ","),
			CReq:   buildCReq(requestID),
		})

		if err != nil {
			return nil, err
		}

		for _, data := range res.GetData() {
			allData = append(allData, data)
		}
	}
	return allData, nil
}

func GetTradeList(ctx context.Context, requestID string, login int64, tickets []int64) (*pb.GetTradeResponse_ResponseData, error) {
	tks := make([]string, len(tickets))
	for i, v := range tickets {
		tks[i] = strconv.FormatInt(v, 10)
	}

	res, err := grpc.TSPCli.GetTrade(ctx, &pb.GetTradeRequest{
		Login:   int32(login),
		Tickets: strings.Join(tks, ","),
		CReq:    buildCReq(requestID),
	})

	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func GetSymbolInfos(ctx context.Context, requestID string, symbol string, originalOnly bool) ([]*pb.SymbolInfo, error) {
	var t pb.GetSymbolRequest_SymbolLevelType
	if originalOnly {
		t = pb.GetSymbolRequest_AllNoLevel
	} else {
		t = pb.GetSymbolRequest_All
	}

	res, err := grpc.TSPCli.GetSymbol(ctx, &pb.GetSymbolRequest{
		Symbol: symbol,
		Type:   t,
		CReq:   buildCReq(requestID),
	})

	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func GetSpreadInfo(ctx context.Context, requestID string) (*pb.GetSpreadInfoResponse_ResponseData, error) {
	res, err := grpc.TSPCli.GetSpreadInfo(ctx, &pb.Null{
		CReq: buildCReq(requestID),
	})

	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func GetSymbolQuoteTime(ctx context.Context, requestID, symbol string) (map[string]*pb.SessionInfo, error) {
	var t = pb.SessionRequest_Quote
	res, err := grpc.TSPCli.GetSession(ctx, &pb.SessionRequest{
		Symbol: symbol,
		Type:   t,
		CReq:   buildCReq(requestID),
	})

	if err != nil {
		return nil, err
	}
	return res.GetData(), nil
}

func ModifyGroups(ctx context.Context, requestID string, login int64, oldGroupName, newGroupName string) error {
	userGroup := make(map[int32]*pb.ModifyGroupsData)
	userGroup[int32(login)] = &pb.ModifyGroupsData{
		OldGroupName: oldGroupName,
		NewGroupName: newGroupName,
	}
	_, err := grpc.TSPCli.ModifyGroups(ctx, &pb.ModifyGroupsRequest{
		UserGroup: userGroup,
		CReq:      buildCReq(requestID),
	})
	if err != nil {
		return err
	}
	return nil
}

func OpenTrade(ctx context.Context, requestID string, login int64, symbol string, cmd int8, volume int64) (*pb.Order, error) {
	symbolInfo, err := repo.SymbolLeverageRepo.GetInfoBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	res, err := grpc.TSPCli.TradeOpen(ctx, &pb.TradeOpenRequest{
		Login:    int32(login),
		Source:   symbolInfo.Original,
		Leverage: symbolInfo.Leverage,
		Cmd:      pb.Cmd(cmd),
		Volume:   int32(volume),
		Client:   pb.TradeClient_Web,
		CReq:     buildCReq(requestID),
	})
	if err != nil {
		return nil, err
	}
	return res.GetData(), nil
}

func buildCReq(requestID string) *pb.CommonRequest {
	hostname, _ := os.Hostname()
	return &pb.CommonRequest{
		PHPHost:   hostname,
		RequestID: requestID,
		PHPTime:   time.Now().UnixMilli(),
	}
}
