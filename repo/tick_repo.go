package repo

import (
	"backend-rms/core/database"
	"backend-rms/model"
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type tickRepo struct{}

var TickRepo = &tickRepo{}

func (r *tickRepo) GetHistoryTick(ctx context.Context, collName string, ctm int64) (tick *model.Tick, err error) {
	var res bson.M
	filter := bson.D{{"ctm", bson.D{{"$lte", ctm}}}}
	opts := options.FindOne().
		SetProjection(bson.D{{"ctm", 1}, {"ask", 1}, {"bid", 1}}).
		SetSort(bson.D{{"ctm", -1}})
	err = database.TickMongoDB.Collection(collName).
		FindOne(ctx, filter, opts).
		Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		err = errors.Wrap(err, "mongo: error in getting history tick")
		return
	}
	tick = &model.Tick{
		Ctm: res["ctm"].(int64),
		Ask: decimal.NewFromFloat(res["ask"].(float64)),
		Bid: decimal.NewFromFloat(res["bid"].(float64)),
	}
	return
}
