/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/dtm-labs/dtm/test/busi"
	"github.com/dtm-labs/logger"
	"github.com/stretchr/testify/assert"
)

// BarrierModel barrier model for gorm
type BarrierModel struct {
	dtmutil.ModelBase
	dtmcli.BranchBarrier
}

// TableName gorm table name
func (BarrierModel) TableName() string { return dtmimp.BarrierTableName }

func TestBaseSqlDB(t *testing.T) {
	asserts := assert.New(t)
	db := dtmutil.DbGet(busi.BusiConf)
	barrier := &dtmcli.BranchBarrier{
		TransType: "saga",
		Gid:       "gid2",
		BranchID:  "branch_id2",
		Op:        dtmimp.OpAction,
		BarrierID: 1,
	}
	db.Must().Exec(fmt.Sprintf("insert into %s(trans_type, gid, branch_id, op, barrier_id, reason) values('saga', 'gid1', 'branch_id1', 'action', '01', 'saga')", dtmimp.BarrierTableName))
	tx, err := db.ToSQLDB().Begin()
	asserts.Nil(err)
	err = barrier.Call(tx, func(tx *sql.Tx) error {
		logger.Debugf("rollback gid2")
		return fmt.Errorf("gid2 error")
	})
	asserts.Error(err, fmt.Errorf("gid2 error"))
	dbr := db.Model(&BarrierModel{}).Where("gid=?", "gid1").Find(&[]BarrierModel{})
	asserts.Equal(dbr.RowsAffected, int64(1))
	dbr = db.Model(&BarrierModel{}).Where("gid=?", "gid2").Find(&[]BarrierModel{})
	asserts.Equal(dbr.RowsAffected, int64(0))
	barrier.BarrierID = 0
	err = barrier.CallWithDB(db.ToSQLDB(), func(tx *sql.Tx) error {
		logger.Debugf("submit gid2")
		return nil
	})
	asserts.Nil(err)
	dbr = db.Model(&BarrierModel{}).Where("gid=?", "gid2").Find(&[]BarrierModel{})
	asserts.Equal(dbr.RowsAffected, int64(1))
}

func TestBaseHttp(t *testing.T) {
	resp, err := dtmcli.GetRestyClient().R().SetQueryParam("panic_string", "1").Post(busi.Busi + "/TestPanic")
	assert.Nil(t, err)
	assert.Contains(t, resp.String(), "panic_string")
	resp, err = dtmcli.GetRestyClient().R().SetQueryParam("panic_error", "1").Post(busi.Busi + "/TestPanic")
	assert.Nil(t, err)
	assert.Contains(t, resp.String(), "panic_error")
}
