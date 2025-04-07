/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package test

import (
	"database/sql"
	"testing"

	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"
	"github.com/dtm-labs/dtm/client/dtmgrpc/dtmgimp"
	"github.com/dtm-labs/dtm/client/workflow"
	"github.com/dtm-labs/dtm/test/busi"
	"github.com/dtm-labs/logger"
	"github.com/stretchr/testify/assert"
)

var ongoingStep = 0

func fetchOngoingStep(dest int) bool {
	c := ongoingStep
	logger.Debugf("ongoing step is: %d", c)
	if c == dest {
		ongoingStep++
		return true
	}
	return false
}

func TestWorkflowSimpleResume(t *testing.T) {
	workflow.SetProtocolForTest(dtmimp.ProtocolHTTP)
	req := busi.GenReqHTTP(30, false, false)
	gid := dtmimp.GetFuncName()
	ongoingStep = 0

	workflow.Register(gid, func(wf *workflow.Workflow, data []byte) error {
		if fetchOngoingStep(0) {
			return dtmcli.ErrOngoing
		}
		var req busi.ReqHTTP
		dtmimp.MustUnmarshal(data, &req)
		_, err := wf.NewBranch().NewRequest().SetBody(req).Post(Busi + "/TransOut")
		return err
	})

	err := workflow.Execute(gid, gid, dtmimp.MustMarshal(req))
	assert.Error(t, err)
	cronTransOnceForwardNow(t, gid, 1000)
	assert.Equal(t, StatusSucceed, getTransStatus(gid))
}

func TestWorkflowGrpcRollbackResume(t *testing.T) {
	workflow.SetProtocolForTest(dtmimp.ProtocolGRPC)
	gid := dtmimp.GetFuncName()
	ongoingStep = 0
	workflow.Register(gid, func(wf *workflow.Workflow, data []byte) error {
		var req busi.ReqGrpc
		dtmgimp.MustProtoUnmarshal(data, &req)
		if fetchOngoingStep(0) {
			return dtmcli.ErrOngoing
		}
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if fetchOngoingStep(4) {
				return dtmcli.ErrOngoing
			}
			_, err := busi.BusiCli.TransOutRevertBSaga(wf.Context, &req)
			return err
		})
		_, err := busi.BusiCli.TransOutBSaga(wf.Context, &req)
		if fetchOngoingStep(1) {
			return dtmcli.ErrOngoing
		}
		if err != nil {
			return err
		}
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if fetchOngoingStep(3) {
				return dtmcli.ErrOngoing
			}
			_, err := busi.BusiCli.TransInRevertBSaga(wf.Context, &req)
			return err
		})
		_, err = busi.BusiCli.TransInBSaga(wf.Context, &req)
		if fetchOngoingStep(2) {
			return dtmcli.ErrOngoing
		}
		return err
	}, func(wf *workflow.Workflow) {
		wf.Options.CompensateErrorBranch = true
	})
	before := getBeforeBalances("mysql")
	req := &busi.ReqGrpc{Amount: 30, TransInResult: "FAILURE"}
	err := workflow.Execute(gid, gid, dtmgimp.MustProtoMarshal(req))
	assert.Error(t, err, dtmcli.ErrOngoing)
	assert.Equal(t, StatusPrepared, getTransStatus(gid))
	cronTransOnceForwardNow(t, gid, 1000)
	assert.Equal(t, StatusPrepared, getTransStatus(gid))
	cronTransOnceForwardNow(t, gid, 1000)
	assert.Equal(t, StatusPrepared, getTransStatus(gid))
	cronTransOnceForwardNow(t, gid, 1000)
	assert.Equal(t, StatusPrepared, getTransStatus(gid))
	cronTransOnceForwardNow(t, gid, 1000)
	assert.Equal(t, StatusPrepared, getTransStatus(gid))
	cronTransOnceForwardNow(t, gid, 1000)
	assert.Equal(t, StatusFailed, getTransStatus(gid))
	assertSameBalance(t, before, "mysql")
}

func TestWorkflowXaResume(t *testing.T) {
	workflow.SetProtocolForTest(dtmimp.ProtocolGRPC)
	ongoingStep = 0
	gid := dtmimp.GetFuncName()
	workflow.Register(gid, func(wf *workflow.Workflow, data []byte) error {
		_, err := wf.NewBranch().DoXa(busi.BusiConf, func(db *sql.DB) ([]byte, error) {
			if fetchOngoingStep(0) {
				return nil, dtmcli.ErrOngoing
			}
			return nil, busi.SagaAdjustBalance(db, busi.TransOutUID, -30, dtmcli.ResultSuccess)
		})
		if err != nil {
			return err
		}
		_, err = wf.NewBranch().DoXa(busi.BusiConf, func(db *sql.DB) ([]byte, error) {
			if fetchOngoingStep(1) {
				return nil, dtmcli.ErrOngoing
			}
			return nil, busi.SagaAdjustBalance(db, busi.TransInUID, 30, dtmcli.ResultSuccess)
		})
		if err != nil {
			return err
		}
		if fetchOngoingStep(2) {
			return dtmcli.ErrOngoing
		}

		return err
	})
	before := getBeforeBalances("mysql")
	err := workflow.Execute(gid, gid, nil)
	assert.Equal(t, dtmcli.ErrOngoing, err)

	cronTransOnceForwardNow(t, gid, 1000)
	assert.Equal(t, StatusPrepared, getTransStatus(gid))
	cronTransOnceForwardNow(t, gid, 1000)
	assert.Equal(t, StatusPrepared, getTransStatus(gid))
	cronTransOnceForwardNow(t, gid, 1000)
	assert.Equal(t, StatusSucceed, getTransStatus(gid))
	assertNotSameBalance(t, before, "mysql")
}
