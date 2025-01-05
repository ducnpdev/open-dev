/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package dtmsvr

import (
	"fmt"
	"time"

	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"
	"github.com/dtm-labs/dtm/dtmsvr/storage"
	"github.com/dtm-labs/logger"
)

// Version store the passin version for dtm server
var Version = ""

func svcSubmit(t *TransGlobal) interface{} {
	if t.TransType == "workflow" {
		t.Status = dtmcli.StatusPrepared
		t.changeStatus(t.ReqExtra["status"], withRollbackReason(t.ReqExtra["rollback_reason"]), withResult(t.ReqExtra["result"]))
		return nil
	}
	t.Status = dtmcli.StatusSubmitted
	branches, err := t.saveNew()

	if err == storage.ErrUniqueConflict {
		dbt := GetTransGlobal(t.Gid)
		if dbt.Status == dtmcli.StatusPrepared {
			dbt.changeStatus(t.Status)
			branches = GetStore().FindBranches(t.Gid)
		} else if dbt.Status != dtmcli.StatusSubmitted {
			return fmt.Errorf("current status '%s', cannot sumbmit. %w", dbt.Status, dtmcli.ErrFailure)
		}
	}
	return t.Process(branches)
}

func svcPrepare(t *TransGlobal) interface{} {
	t.Status = dtmcli.StatusPrepared
	_, err := t.saveNew()
	if err == storage.ErrUniqueConflict {
		dbt := GetTransGlobal(t.Gid)
		if dbt.Status != dtmcli.StatusPrepared {
			return fmt.Errorf("current status '%s', cannot prepare. %w", dbt.Status, dtmcli.ErrFailure)
		}
		return nil
	}
	return err
}

func svcPrepareWorkflow(t *TransGlobal) (*storage.TransGlobalStore, []TransBranch, error) {
	t.Status = dtmcli.StatusPrepared
	_, err := t.saveNew()
	if err == storage.ErrUniqueConflict { // transaction exists, query the branches
		st := GetStore()
		return st.FindTransGlobalStore(t.Gid), st.FindBranches(t.Gid), nil
	}
	return &t.TransGlobalStore, []TransBranch{}, err
}

func svcAbort(t *TransGlobal) interface{} {
	dbt := GetTransGlobal(t.Gid)
	if dbt.TransType == "msg" && dbt.Status == dtmcli.StatusPrepared {
		dbt.changeStatus(dtmcli.StatusFailed)
		return nil
	}
	if t.TransType != "xa" && t.TransType != "tcc" || dbt.Status != dtmcli.StatusPrepared && dbt.Status != dtmcli.StatusAborting {
		return fmt.Errorf("trans type: '%s' current status '%s', cannot abort. %w", dbt.TransType, dbt.Status, dtmcli.ErrFailure)
	}
	dbt.changeStatus(dtmcli.StatusAborting, withRollbackReason(t.RollbackReason))
	branches := GetStore().FindBranches(t.Gid)
	return dbt.Process(branches)
}

func svcForceStop(t *TransGlobal) interface{} {
	dbt := GetTransGlobal(t.Gid)
	if dbt.Status == dtmcli.StatusSucceed || dbt.Status == dtmcli.StatusFailed {
		return fmt.Errorf("global transaction force stop error. status: %s. error: %w", dbt.Status, dtmcli.ErrFailure)
	}
	dbt.changeStatus(dtmcli.StatusFailed)
	return nil
}

func svcResetNextCronTime(t *TransGlobal) error {
	dbt := GetTransGlobal(t.Gid)
	return dbt.resetNextCronTime()
}

func svcRegisterBranch(transType string, branch *TransBranch, data map[string]string) error {
	branches := []TransBranch{*branch, *branch}
	if transType == "tcc" {
		for i, b := range []string{dtmimp.OpCancel, dtmimp.OpConfirm} {
			branches[i].Op = b
			branches[i].URL = data[b]
		}
	} else if transType == "xa" {
		branches[0].Op = dtmimp.OpRollback
		branches[0].URL = data["url"]
		branches[1].Op = dtmimp.OpCommit
		branches[1].URL = data["url"]
	} else if transType == "workflow" {
		if data["sync"] == "" && conf.UpdateBranchSync == 0 {
			now := time.Now()
			updateBranchAsyncChan <- branchStatus{
				gid:        branch.Gid,
				branchID:   branch.BranchID,
				op:         data["op"],
				status:     data["status"],
				finishTime: &now,
			}
			return nil
		}
		branches = []TransBranch{*branch}
		branches[0].Status = data["status"]
		branches[0].Op = data["op"]
	} else {
		return fmt.Errorf("unknow trans type: %s", transType)
	}

	err := dtmimp.CatchP(func() {
		GetStore().LockGlobalSaveBranches(branch.Gid, dtmcli.StatusPrepared, branches, -1)
	})
	if err == storage.ErrNotFound {
		msg := fmt.Sprintf("no trans with gid: %s status: %s found", branch.Gid, dtmcli.StatusPrepared)
		logger.Errorf(msg)
		return dtmcli.ErrorMessage2Error(msg, dtmcli.ErrFailure)
	}
	logger.Infof("LockGlobalSaveBranches result: %v: gid: %s old status: %s branches: %s",
		err, branch.Gid, dtmcli.StatusPrepared, dtmimp.MustMarshalString(branches))
	return err
}
