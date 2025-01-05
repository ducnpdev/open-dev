/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package busi

import (
	"context"
	"database/sql"

	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func init() {
	setupFuncs["BarrierSetup"] = func(app *gin.Engine) {
		app.POST(BusiAPI+"/SagaBTransIn", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := MustBarrierFromGin(c)
			return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransInUID, reqFrom(c).Amount, reqFrom(c).TransInResult)
			})
		}))
		app.POST(BusiAPI+"/SagaBTransInCom", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := MustBarrierFromGin(c)
			return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransInUID, -reqFrom(c).Amount, "")
			})
		}))
		app.POST(BusiAPI+"/SagaB2TransIn", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := MustBarrierFromGin(c)
			err := barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransInUID, reqFrom(c).Amount/2, reqFrom(c).TransInResult)
			})
			if err != nil {
				return err
			}
			return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransInUID, reqFrom(c).Amount/2, reqFrom(c).TransInResult)
			})
		}))
		app.POST(BusiAPI+"/SagaB2TransInCom", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := MustBarrierFromGin(c)
			err := barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransInUID, -reqFrom(c).Amount/2, "")
			})
			if err != nil {
				return err
			}
			return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransInUID, -reqFrom(c).Amount, "")
			})
		}))
		app.POST(BusiAPI+"/SagaBTransOut", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := MustBarrierFromGin(c)
			return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransOutUID, -reqFrom(c).Amount, reqFrom(c).TransOutResult)
			})
		}))
		app.POST(BusiAPI+"/SagaBTransOutCom", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := MustBarrierFromGin(c)
			return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransOutUID, reqFrom(c).Amount, "")
			})
		}))
		app.POST(BusiAPI+"/SagaBTransOutGorm", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			req := reqFrom(c)
			barrier := MustBarrierFromGin(c)
			tx := dbGet().DB.Begin()
			return barrier.Call(tx.Statement.ConnPool.(*sql.Tx), func(tx1 *sql.Tx) error {
				return tx.Exec("update dtm_busi.user_account set balance = balance + ? where user_id = ?", -req.Amount, TransOutUID).Error
			})
		}))

		app.POST(BusiAPI+"/TccBTransInTry", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			req := reqFrom(c)
			if req.TransInResult != "" {
				return string2DtmError(req.TransInResult)
			}
			return MustBarrierFromGin(c).CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return tccAdjustTrading(tx, TransInUID, req.Amount)
			})
		}))
		app.POST(BusiAPI+"/TccBTransInConfirm", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			return MustBarrierFromGin(c).CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return tccAdjustBalance(tx, TransInUID, reqFrom(c).Amount)
			})
		}))
		app.POST(BusiAPI+"/TccBTransInCancel", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			return MustBarrierFromGin(c).CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return tccAdjustTrading(tx, TransInUID, -reqFrom(c).Amount)
			})
		}))
		app.POST(BusiAPI+"/SagaMultiSource", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := MustBarrierFromGin(c)
			transOutSource := pdbGet()
			err := barrier.CallWithDB(transOutSource, func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransOutUID, -reqFrom(c).Amount, reqFrom(c).TransOutResult)
			})
			if err != nil {
				return err
			}
			transInSource := pdbGet()
			return barrier.CallWithDB(transInSource, func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransInUID, reqFrom(c).Amount, reqFrom(c).TransInResult)
			})
		}))
		app.POST(BusiAPI+"/SagaMultiSourceRevert", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := MustBarrierFromGin(c)
			transOutSource := pdbGet()
			err := barrier.CallWithDB(transOutSource, func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransOutUID, +reqFrom(c).Amount, "")
			})
			if err != nil {
				return err
			}
			transInSource := pdbGet()
			return barrier.CallWithDB(transInSource, func(tx *sql.Tx) error {
				return SagaAdjustBalance(tx, TransInUID, -reqFrom(c).Amount, "")
			})
		}))
		app.POST(BusiAPI+"/SagaRedisTransIn", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			return MustBarrierFromGin(c).RedisCheckAdjustAmount(RedisGet(), GetRedisAccountKey(TransInUID), reqFrom(c).Amount, 7*86400)
		}))
		app.POST(BusiAPI+"/SagaRedisTransInCom", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			return MustBarrierFromGin(c).RedisCheckAdjustAmount(RedisGet(), GetRedisAccountKey(TransInUID), -reqFrom(c).Amount, 7*86400)
		}))
		app.POST(BusiAPI+"/SagaRedisTransOut", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			return MustBarrierFromGin(c).RedisCheckAdjustAmount(RedisGet(), GetRedisAccountKey(TransOutUID), -reqFrom(c).Amount, 7*86400)
		}))
		app.POST(BusiAPI+"/SagaRedisTransOutCom", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			return MustBarrierFromGin(c).RedisCheckAdjustAmount(RedisGet(), GetRedisAccountKey(TransOutUID), reqFrom(c).Amount, 7*86400)
		}))
		app.POST(BusiAPI+"/SagaMongoTransIn", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			return MustBarrierFromGin(c).MongoCall(MongoGet(), func(sc mongo.SessionContext) error {
				return SagaMongoAdjustBalance(sc, sc.Client(), TransInUID, reqFrom(c).Amount, reqFrom(c).TransInResult)
			})
		}))
		app.POST(BusiAPI+"/SagaMongoTransInCom", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			return MustBarrierFromGin(c).MongoCall(MongoGet(), func(sc mongo.SessionContext) error {
				return SagaMongoAdjustBalance(sc, sc.Client(), TransInUID, -reqFrom(c).Amount, "")
			})
		}))
		app.POST(BusiAPI+"/SagaMongoTransOut", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			return MustBarrierFromGin(c).MongoCall(MongoGet(), func(sc mongo.SessionContext) error {
				return SagaMongoAdjustBalance(sc, sc.Client(), TransOutUID, -reqFrom(c).Amount, reqFrom(c).TransOutResult)
			})
		}))
		app.POST(BusiAPI+"/SagaMongoTransOutCom", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			return MustBarrierFromGin(c).MongoCall(MongoGet(), func(sc mongo.SessionContext) error {
				return SagaMongoAdjustBalance(sc, sc.Client(), TransOutUID, reqFrom(c).Amount, "")
			})
		}))
		app.POST(BusiAPI+"/TccBTransOutTry", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			req := reqFrom(c)
			if req.TransOutResult != "" {
				return string2DtmError(req.TransOutResult)
			}
			bb := MustBarrierFromGin(c)
			if req.Store == Redis {
				return bb.RedisCheckAdjustAmount(RedisGet(), GetRedisAccountKey(TransOutUID), req.Amount, 7*86400)
			} else if req.Store == Mongo {
				return bb.MongoCall(MongoGet(), func(sc mongo.SessionContext) error {
					return SagaMongoAdjustBalance(sc, sc.Client(), TransOutUID, -req.Amount, "")
				})
			}

			return bb.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return tccAdjustTrading(tx, TransOutUID, -req.Amount)
			})
		}))
		app.POST(BusiAPI+"/TccBTransOutConfirm", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			if reqFrom(c).Store == Redis || reqFrom(c).Store == Mongo {
				return nil
			}
			return MustBarrierFromGin(c).CallWithDB(pdbGet(), func(tx *sql.Tx) error {
				return tccAdjustBalance(tx, TransOutUID, -reqFrom(c).Amount)
			})
		}))
		app.POST(BusiAPI+"/TccBTransOutCancel", dtmutil.WrapHandler(TccBarrierTransOutCancel))
	}
}

// TccBarrierTransOutCancel will be use in test
func TccBarrierTransOutCancel(c *gin.Context) interface{} {
	req := reqFrom(c)
	bb := MustBarrierFromGin(c)
	if req.Store == Redis {
		return bb.RedisCheckAdjustAmount(RedisGet(), GetRedisAccountKey(TransOutUID), -req.Amount, 7*86400)
	}
	if req.Store == Mongo {
		return bb.MongoCall(MongoGet(), func(sc mongo.SessionContext) error {
			return SagaMongoAdjustBalance(sc, sc.Client(), TransOutUID, reqFrom(c).Amount, "")
		})
	}
	return bb.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
		return tccAdjustTrading(tx, TransOutUID, reqFrom(c).Amount)
	})
}

func (s *busiServer) TransInBSaga(ctx context.Context, in *ReqGrpc) (*emptypb.Empty, error) {
	barrier := MustBarrierFromGrpc(ctx)
	return &emptypb.Empty{}, barrier.Call(txGet(), func(tx *sql.Tx) error {
		return sagaGrpcAdjustBalance(tx, TransInUID, in.Amount, in.TransInResult)
	})
}

func (s *busiServer) TransOutBSaga(ctx context.Context, in *ReqGrpc) (*emptypb.Empty, error) {
	barrier := MustBarrierFromGrpc(ctx)
	return &emptypb.Empty{}, barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
		return sagaGrpcAdjustBalance(tx, TransOutUID, -in.Amount, in.TransOutResult)
	})
}

func (s *busiServer) TransInRevertBSaga(ctx context.Context, in *ReqGrpc) (*emptypb.Empty, error) {
	barrier := MustBarrierFromGrpc(ctx)
	return &emptypb.Empty{}, barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
		return sagaGrpcAdjustBalance(tx, TransInUID, -in.Amount, "")
	})
}

func (s *busiServer) TransOutRevertBSaga(ctx context.Context, in *ReqGrpc) (*emptypb.Empty, error) {
	barrier := MustBarrierFromGrpc(ctx)
	return &emptypb.Empty{}, barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
		return sagaGrpcAdjustBalance(tx, TransOutUID, in.Amount, "")
	})
}

func (s *busiServer) TransInRedis(ctx context.Context, in *ReqGrpc) (*emptypb.Empty, error) {
	barrier := MustBarrierFromGrpc(ctx)
	return &emptypb.Empty{}, barrier.RedisCheckAdjustAmount(RedisGet(), GetRedisAccountKey(TransInUID), int(in.Amount), 86400)
}

func (s *busiServer) TransOutRedis(ctx context.Context, in *ReqGrpc) (*emptypb.Empty, error) {
	barrier := MustBarrierFromGrpc(ctx)
	return &emptypb.Empty{}, barrier.RedisCheckAdjustAmount(RedisGet(), GetRedisAccountKey(TransOutUID), int(-in.Amount), 86400)
}

func (s *busiServer) TransInRevertRedis(ctx context.Context, in *ReqGrpc) (*emptypb.Empty, error) {
	barrier := MustBarrierFromGrpc(ctx)
	return &emptypb.Empty{}, barrier.RedisCheckAdjustAmount(RedisGet(), GetRedisAccountKey(TransInUID), -int(in.Amount), 86400)
}

func (s *busiServer) TransOutRevertRedis(ctx context.Context, in *ReqGrpc) (*emptypb.Empty, error) {
	barrier := MustBarrierFromGrpc(ctx)
	return &emptypb.Empty{}, barrier.RedisCheckAdjustAmount(RedisGet(), GetRedisAccountKey(TransOutUID), int(in.Amount), 86400)
}

func (s *busiServer) QueryPreparedB(ctx context.Context, in *ReqGrpc) (*emptypb.Empty, error) {
	barrier := MustBarrierFromGrpc(ctx)
	err := barrier.QueryPrepared(dbGet().ToSQLDB())
	return &emptypb.Empty{}, dtmgrpc.DtmError2GrpcError(err)
}

func (s *busiServer) QueryPreparedRedis(ctx context.Context, in *ReqGrpc) (*emptypb.Empty, error) {
	barrier := MustBarrierFromGrpc(ctx)
	err := barrier.RedisQueryPrepared(RedisGet(), 86400)
	return &emptypb.Empty{}, dtmgrpc.DtmError2GrpcError(err)
}
