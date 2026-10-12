package service

import (
	"backend-rms/model"
	"backend-rms/repo"
	"backend-rms/util"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type copyRecordService struct{}

type Fund struct {
	Amount            decimal.Decimal
	OtherMasterAmount decimal.Decimal
}

type CopyPLInfo struct {
	PeriodPl    decimal.Decimal
	ClosedPl    decimal.Decimal
	FloatingPl  decimal.Decimal
	MarginLevel decimal.Decimal
}

var CopyRecordService = &copyRecordService{}

func (s *copyRecordService) GetAll(followerId, masterId int64) (copiesRecords []model.CopyRecord, err error) {
	return repo.CopyRecordRepo.GetAll(followerId, masterId)
}

func (s *copyRecordService) GetCopies(masterId int64, groups []string) (copies []model.CopiesRecord, err error) {
	copies, err = repo.CopyRecordRepo.GetCopies(masterId, groups)
	return copies, errors.WithMessage(err, util.ErrWrapMessage())
}

// 获取时间点跟随者们跟随高手的数量
func (s *copyRecordService) GetPointFollowMasterNums(dateTime time.Time, followers []int64, groups []string) (map[int64]int, error) {
	copiesMasters, err := repo.CopyRecordRepo.GetCopiesMasters(dateTime, followers, groups)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	followMasterNums := make(map[int64]int, 0)
	for _, v := range copiesMasters {
		if _, ok := followMasterNums[v.Mt4To]; ok {
			followMasterNums[v.Mt4To] = followMasterNums[v.Mt4To] + 1
			continue
		}
		followMasterNums[v.Mt4To] = 1
	}
	return followMasterNums, nil
}

func (s *copyRecordService) GetPointCopiedAmounts(dateTime time.Time, followers []int64, groups []string, masterId int64) (map[int64]Fund, error) {
	copiesMasters, err := repo.CopyRecordRepo.GetCopiesMasters(dateTime, followers, groups)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	copyAmounts := make(map[int64]Fund)
	for _, v := range copiesMasters {
		if _, ok := copyAmounts[v.Mt4To]; !ok {
			copyAmounts[v.Mt4To] = Fund{
				Amount:            decimal.Zero,
				OtherMasterAmount: decimal.Zero,
			}
		}

		if v.Mt4From == masterId {
			copyAmounts[v.Mt4To] = Fund{
				Amount:            copyAmounts[v.Mt4To].Amount.Add(v.Amount),
				OtherMasterAmount: copyAmounts[v.Mt4To].OtherMasterAmount,
			}
		} else {
			copyAmounts[v.Mt4To] = Fund{
				Amount:            copyAmounts[v.Mt4To].Amount,
				OtherMasterAmount: copyAmounts[v.Mt4To].OtherMasterAmount.Add(v.Amount),
			}
		}
	}
	return copyAmounts, nil
}

func (s *copyRecordService) GetFollowingLogins(masterId int64, groups []string) ([]int64, error) {
	return repo.CopyRecordRepo.GetFollowingLogins(masterId, groups)
}

func (s *copyRecordService) GetCopiedPLInfos(groups []string, followers []int64, startTime, endTime time.Time, requestID string) (map[int64]CopyPLInfo, error) {
	endTimePLs := make(map[int64]decimal.Decimal, 0)
	startTimePLs := make(map[int64]decimal.Decimal, 0)
	var err error
	beforeSixMonth := time.Now().AddDate(0, -6, 0)
	if startTime.After(beforeSixMonth) && endTime.After(beforeSixMonth) {
		startTimePLs, err = FloatingPLService.GetPointAccountPLMap(groups, followers, startTime, requestID)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
		endTimePLs, err = FloatingPLService.GetPointAccountPLMap(groups, followers, endTime, requestID)
		if err != nil {
			return nil, errors.WithMessage(err, util.ErrWrapMessage())
		}
	}

	closedPLs, err := TradeService.GetAccountClosedPL(followers, groups, startTime, endTime, "")
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	floatingPls, err := CacheService.GetAllAccountFloatingPLMap()
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	marginLevels, err := MarginService.GetMarginLevelBatch(requestID, followers)
	if err != nil {
		return nil, errors.WithMessage(err, util.ErrWrapMessage())
	}

	copyPLInfos := make(map[int64]CopyPLInfo, 0)
	for _, v := range followers {
		startTimePL := decimal.Zero
		if _, ok := startTimePLs[v]; ok {
			startTimePL = startTimePLs[v]
		}

		endTimePL := decimal.Zero
		if _, ok := endTimePLs[v]; ok {
			endTimePL = endTimePLs[v]
		}

		closedPL := decimal.Zero
		if _, ok := closedPLs[v]; ok {
			closedPL = closedPLs[v]
		}
		periodPl := endTimePL.Sub(startTimePL).Add(closedPL)

		floatingPl := decimal.Zero
		if _, ok := floatingPls[v]; ok {
			floatingPl = floatingPls[v]
		}

		marginLevel := decimal.Zero
		if _, ok := marginLevels[v]; ok {
			marginLevel = marginLevels[v]
		}

		copyPLInfos[v] = CopyPLInfo{
			PeriodPl:    periodPl,
			ClosedPl:    closedPL,
			FloatingPl:  floatingPl,
			MarginLevel: marginLevel,
		}
	}
	return copyPLInfos, nil
}
