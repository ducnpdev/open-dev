package busi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/dtm-labs/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// PopulateDB populate example mysql data
func PopulateDB(skipDrop bool) {
	ResetXaData()
	file := fmt.Sprintf("%s/busi.%s.sql", dtmutil.GetSQLDir(), BusiConf.Driver)
	dtmutil.RunSQLScript(BusiConf, file, skipDrop)
	file = fmt.Sprintf("%s/dtmcli.barrier.%s.sql", dtmutil.GetSQLDir(), BusiConf.Driver)
	dtmutil.RunSQLScript(BusiConf, file, skipDrop)
	file = fmt.Sprintf("%s/dtmsvr.storage.%s.sql", dtmutil.GetSQLDir(), BusiConf.Driver)
	dtmutil.RunSQLScript(BusiConf, file, skipDrop)
	_, err := RedisGet().FlushAll(context.Background()).Result() // redis barrier need clear
	dtmimp.E2P(err)
	SetRedisBothAccount(10000, 10000)
	SetupMongoBarrierAndBusi()
}

// TransOutUID 1
const TransOutUID = 1

// TransInUID 2
const TransInUID = 2

// Redis 1
const Redis = "redis"

// Mongo 1
const Mongo = "mongo"

func handleGrpcBusiness(in *ReqGrpc, result1 string, result2 string, busi string) error {
	res := dtmimp.OrString(result1, result2, dtmcli.ResultSuccess)
	logger.Debugf("grpc busi %s %v %s %s result: %s", busi, in, result1, result2, res)
	if res == dtmcli.ResultSuccess {
		return nil
	} else if res == dtmcli.ResultFailure {
		return status.New(codes.Aborted, fmt.Sprintf("reason:%s", MainSwitch.FailureReason.Fetch())).Err()
	} else if res == dtmcli.ResultOngoing {
		return status.New(codes.FailedPrecondition, dtmcli.ResultOngoing).Err()
	}
	return status.New(codes.Internal, fmt.Sprintf("unknow result %s", res)).Err()
}

func handleGeneralBusiness(c *gin.Context, result1 string, result2 string, busi string) interface{} {
	info := infoFromContext(c)
	res := dtmimp.OrString(result1, result2, dtmcli.ResultSuccess)
	logger.Debugf("%s %s result: %s", busi, info.String(), res)
	if res == "ERROR" {
		return errors.New("ERROR from user")
	}
	if res == dtmimp.ResultFailure {
		return fmt.Errorf("reason:%s. %w", MainSwitch.FailureReason.Fetch(), dtmimp.ErrFailure)
	}
	return string2DtmError(res)
}

// old business handler. for compatible usage
func handleGeneralBusinessCompatible(c *gin.Context, result1 string, result2 string, busi string) (interface{}, error) {
	info := infoFromContext(c)
	res := dtmimp.OrString(result1, result2, dtmcli.ResultSuccess)
	logger.Debugf("%s %s result: %s", busi, info.String(), res)
	if res == "ERROR" {
		return nil, errors.New("ERROR from user")
	}
	return map[string]interface{}{"dtm_result": res}, nil
}

func sagaGrpcAdjustBalance(db dtmcli.DB, uid int, amount int64, result string) error {
	if result == dtmcli.ResultFailure {
		return status.New(codes.Aborted, dtmcli.ResultFailure).Err()
	}
	_, err := dtmimp.DBExec(BusiConf.Driver, db, "update dtm_busi.user_account set balance = balance + ? where user_id = ?", amount, uid)
	return err
}

// SagaAdjustBalance 1
func SagaAdjustBalance(db dtmcli.DB, uid int, amount int, result string) error {
	if strings.Contains(result, dtmcli.ResultFailure) {
		return dtmcli.ErrFailure
	}
	_, err := dtmimp.DBExec(BusiConf.Driver, db, "update dtm_busi.user_account set balance = balance + ? where user_id = ?", amount, uid)
	return err
}

// SagaMongoAdjustBalance 1
func SagaMongoAdjustBalance(ctx context.Context, mc *mongo.Client, uid int, amount int, result string) error {
	if strings.Contains(result, dtmcli.ResultFailure) {
		return dtmcli.ErrFailure
	}
	_, err := mc.Database("dtm_busi").Collection("user_account").UpdateOne(ctx,
		bson.D{{Key: "user_id", Value: uid}},
		bson.D{{Key: "$inc", Value: bson.D{{Key: "balance", Value: amount}}}})
	logger.Debugf("dtm_busi.user_account $inc balance of %d by %d err: %v", uid, amount, err)
	if err != nil {
		return err
	}
	var res bson.M
	err = mc.Database("dtm_busi").Collection("user_account").FindOne(ctx,
		bson.D{{Key: "user_id", Value: uid}}).Decode(&res)
	if err != nil {
		return err
	}
	balance := res["balance"].(float64)
	if balance < 0 {
		return fmt.Errorf("balance not enough %w", dtmcli.ErrFailure)
	}
	return nil
}

func tccAdjustTrading(db dtmcli.DB, uid int, amount int) error {
	affected, err := dtmimp.DBExec(BusiConf.Driver, db, `update dtm_busi.user_account
		set trading_balance=trading_balance+?
		where user_id=? and trading_balance + ? + balance >= 0`, amount, uid, amount)
	if err == nil && affected == 0 {
		return fmt.Errorf("update error, maybe balance not enough")
	}
	return err
}

func tccAdjustBalance(db dtmcli.DB, uid int, amount int) error {
	affected, err := dtmimp.DBExec(BusiConf.Driver, db, `update dtm_busi.user_account
		set trading_balance=trading_balance-?,
		balance=balance+? where user_id=?`, amount, amount, uid)
	if err == nil && affected == 0 {
		return fmt.Errorf("update user_account 0 rows")
	}
	return err
}
