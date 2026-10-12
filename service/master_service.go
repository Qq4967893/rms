package service

import (
	pb "backend-rms/grpc/et6api"
	"backend-rms/grpc/tsp"
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"context"
	"math"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

type Master struct {
	Login               int64
	Group               string
	Passport            int64
	FloatingPL          decimal.Decimal
	PeriodClosedPL      decimal.Decimal
	StartTimeFloatingPL decimal.Decimal
	EndTimeFloatingPL   decimal.Decimal
	NetPeriodPL         decimal.Decimal
	IsLimited           bool
	Balance             decimal.Decimal
	FollowerCount       int64
	CopiedAmount        decimal.Decimal
	MarginLevel         decimal.Decimal
	IsPresent           bool
}

func (m *Master) GetNetPeriodPL() decimal.Decimal {
	if m.IsLimited {
		return decimal.Zero
	}
	return m.EndTimeFloatingPL.Sub(m.StartTimeFloatingPL).Add(m.PeriodClosedPL)
}

type MasterHistoricalRecord struct {
	Login         int64
	Commission    decimal.Decimal
	ClosedPL      decimal.Decimal
	FollowerCount int64
	FollowedTimes int64
	StartTime     time.Time
	EndTime       time.Time
}

type masterService struct{}

var MasterService = &masterService{}

func (ms *masterService) GetMasterList(requestID, isPresent string, groups []string, login int64, startTime, endTime time.Time) ([]*Master, error) {
	var err error
	var reportDailyTime *time.Time
	timeEnd := carbon.Time2Carbon(endTime)
	yesterday := timeEnd.Yesterday().Carbon2Time()
	reportDailyTime, err = ReportDailyService.GetReportDailyTime(yesterday)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	if reportDailyTime == nil {
		return nil, errors.Errorf("time %s does not exist in report daily table", yesterday)
	}

	today := carbon.Time2Carbon(*reportDailyTime).AddSeconds(1).ToDateTimeString()

	var logins []int64
	var masters []*Master
	loginMap := make(map[int64]struct{})
	timeStart := carbon.Time2Carbon(startTime)
	timeLimit := carbon.Now().SubMonths(6).StartOfMinute()
	{
		var masterList, currentMasterList []*model.Master
		g := new(errgroup.Group)
		g.Go(func() error {
			masterList, err = repo.MasterRepo.GetMasterList(groups, login, startTime, endTime)
			return errors.WithMessage(err, util.ErrWrapMessage())
		})
		g.Go(func() error {
			currentMasterList, err = repo.MasterRepo.GetCurrentMasterList(groups, login)
			return errors.WithMessage(err, util.ErrWrapMessage())
		})
		if err = g.Wait(); err != nil {
			return nil, errors.Wrap(err, "g.Wait is failed")
		}

		var currentMasterMap = make(map[int64]struct{}, len(currentMasterList))
		for _, c := range currentMasterList {
			currentMasterMap[c.Login] = struct{}{}
		}

		for _, m := range masterList {
			if _, ok := loginMap[m.Login]; ok {
				continue
			}
			master := &Master{
				Login:    m.Login,
				Group:    m.Group,
				Passport: m.Passport,
			}
			_, master.IsPresent = currentMasterMap[m.Login]
			if (isPresent == "yes" && !master.IsPresent) || (isPresent == "no" && master.IsPresent) {
				continue
			}
			logins = append(logins, m.Login)
			masters = append(masters, master)
			loginMap[m.Login] = struct{}{}
		}
	}

	if len(masters) == 0 {
		return nil, nil
	}

	equityMap := make(map[int64]*pb.EquityInfo, len(logins))
	closedPLMap := make(map[int64]decimal.Decimal, len(logins))
	startTimeFloatingPLMap := make(map[int64]decimal.Decimal, len(logins))
	endTimeFloatingPLMap := make(map[int64]decimal.Decimal, len(logins))
	yesterdayBalanceMap := make(map[int64]decimal.Decimal, len(logins))
	todayNetCashFlowMap := make(map[int64]decimal.Decimal, len(logins))
	todayClosedPLMap := make(map[int64]decimal.Decimal, len(logins))
	copiedAmountMap := make(map[int64]decimal.Decimal, len(logins))
	followerCountMap := make(map[int64]int64, len(logins))

	g := new(errgroup.Group)
	g.Go(func() error {
		var equityInfos []*pb.EquityInfo
		equityInfos, err = tsp.GetUserEquity(context.Background(), requestID, logins)
		if err != nil {
			return errors.WithMessage(err, util.ErrWrapMessage())
		}
		for _, equityInfo := range equityInfos {
			equityMap[int64(equityInfo.Login)] = equityInfo
		}
		return nil
	})

	g.Go(func() error {
		closedPLMap, err = TradeService.GetClosedPLMap(logins, timeStart.ToDateTimeString(), timeEnd.ToDateTimeString())
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	var isLimited bool
	if (timeStart.Gte(timeLimit)) && (timeEnd.Gte(timeLimit)) {
		g.Go(func() error {
			startTimeFloatingPLMap, err = FloatingPLService.GetPointAccountPLMap(groups, logins, startTime, requestID)
			return errors.WithMessage(err, util.ErrWrapMessage())
		})
		g.Go(func() error {
			endTimeFloatingPLMap, err = FloatingPLService.GetPointAccountPLMap(groups, logins, endTime, requestID)
			return errors.WithMessage(err, util.ErrWrapMessage())
		})
	} else {
		isLimited = true
	}

	g.Go(func() error {
		yesterdayBalanceMap, err = ReportDailyService.GetLoginBalanceMapByTime(yesterday, logins)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	g.Go(func() error {
		todayNetCashFlowMap, err = CashFlowService.GetNetCashFlowMap(logins, today, timeEnd.ToDateTimeString())
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	g.Go(func() error {
		todayClosedPLMap, err = TradeService.GetClosedPLMap(logins, today, timeEnd.ToDateTimeString())
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	g.Go(func() error {
		copiedAmountMap, err = repo.MasterRepo.GetCopiedAmountMap(groups, logins, endTime)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	g.Go(func() error {
		followerCountMap, err = repo.MasterRepo.GetFollowerCountMap(groups, logins)
		return errors.WithMessage(err, util.ErrWrapMessage())
	})

	if err = g.Wait(); err != nil {
		return nil, errors.Wrap(err, "g.Wait is failed")
	}

	for _, m := range masters {
		equityInfo := equityMap[m.Login]
		marginLevel := decimal.NewFromFloat(equityInfo.MarginLevel)
		floatingPL := decimal.NewFromFloat(equityInfo.Equity).Sub(decimal.NewFromFloat(equityInfo.Balance))
		balance := yesterdayBalanceMap[login].Add(todayNetCashFlowMap[login]).Add(todayClosedPLMap[login])
		m.FloatingPL = floatingPL
		m.PeriodClosedPL = closedPLMap[m.Login]
		m.StartTimeFloatingPL = startTimeFloatingPLMap[m.Login]
		m.EndTimeFloatingPL = endTimeFloatingPLMap[m.Login]
		m.NetPeriodPL = m.GetNetPeriodPL()
		m.IsLimited = isLimited
		m.Balance = balance
		m.FollowerCount = followerCountMap[m.Login]
		m.CopiedAmount = copiedAmountMap[m.Login]
		m.MarginLevel = marginLevel
	}

	return masters, nil
}

func (ms *masterService) GetMasterHistoricalList(login int64, paging *model.Paging) ([]*MasterHistoricalRecord, int, error) {
	groups, err := GroupService.GetGroupsByBook("ALL")
	if err != nil {
		return nil, 0, errors.WithMessage(err, util.ErrWrapMessage())
	}
	masterRecords, err := repo.MasterRepo.GetMasterHistoricalList(login)
	if err != nil {
		return nil, 0, errors.WithMessage(err, util.ErrWrapMessage())
	}
	if len(masterRecords) == 0 {
		return nil, 0, nil
	}
	historicalRecords, err := ms.buildMasterHistoricalRecords(masterRecords)
	if err != nil {
		return nil, 0, errors.WithMessage(err, util.ErrWrapMessage())
	}

	start := paging.Offset
	if start > len(historicalRecords) {
		start = len(historicalRecords)
	}

	end := paging.Offset + paging.Limit
	if end > len(historicalRecords) {
		end = len(historicalRecords)
	}

	totalCount := len(historicalRecords)
	historicalRecords = historicalRecords[start:end]

	g := new(errgroup.Group)
	for _, record := range historicalRecords {
		record := record
		startTime := record.StartTime
		endTime := record.EndTime
		if endTime.IsZero() {
			endTime = time.Now()
		}
		g.Go(func() error {
			commission, err := repo.MasterRepo.GetCommission(login, startTime, endTime)
			if err != nil {
				return errors.WithMessage(err, util.ErrWrapMessage())
			}
			record.Commission = commission
			return nil
		})
		g.Go(func() error {
			closedPL, err := repo.TradeRepo.GetClosedPL(login, startTime, endTime)
			if err != nil {
				return errors.WithMessage(err, util.ErrWrapMessage())
			}
			record.ClosedPL = closedPL
			return nil
		})
		g.Go(func() error {
			followerCount, err := repo.MasterRepo.GetFollowerCount(groups, login, startTime, endTime)
			if err != nil {
				return errors.WithMessage(err, util.ErrWrapMessage())
			}
			record.FollowerCount = followerCount
			return nil
		})
		g.Go(func() error {
			followedTimes, err := repo.MasterRepo.GetFollowedTimes(groups, login, startTime, endTime)
			if err != nil {
				return errors.WithMessage(err, util.ErrWrapMessage())
			}
			record.FollowedTimes = followedTimes
			return nil
		})
	}
	if err = g.Wait(); err != nil {
		return nil, 0, errors.Wrap(err, "g.Wait is failed")
	}

	return historicalRecords, totalCount, nil
}

func (ms *masterService) buildMasterHistoricalRecords(masterRecords []*model.MasterRecord) ([]*MasterHistoricalRecord, error) {
	capacity := int(math.Ceil(float64(len(masterRecords)) / 2))
	historicalRecords := make([]*MasterHistoricalRecord, 0, capacity)
	for i := 0; i < len(masterRecords); i += 2 {
		var startRecord, endRecord *model.MasterRecord
		startRecord = masterRecords[i]
		if i+1 < len(masterRecords) {
			endRecord = masterRecords[i+1]
		} else {
			endRecord = &model.MasterRecord{
				Type:      0,
				CreatedAt: time.Time{},
			}
		}
		if startRecord.Type != 1 || endRecord.Type != 0 {
			return nil, errors.Errorf("Incorrect master record order, startRecord: %+v, endRecord: %+v", *startRecord, *endRecord)
		}
		historicalRecords = append(historicalRecords, &MasterHistoricalRecord{
			Login:     startRecord.Login,
			StartTime: startRecord.CreatedAt,
			EndTime:   endRecord.CreatedAt,
		})
	}
	return historicalRecords, nil
}
