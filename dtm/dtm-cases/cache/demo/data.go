package demo

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/dtm-labs/dtmcli/logger"
	"github.com/dtm-labs/rockscache"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
)

const dbKey = "key1"

var rdbKey = "key1"

var rdb = redis.NewClient(&redis.Options{
	Addr:     "en.dtm.pub:6379",
	Username: "root",
	Password: "",
})

var dc = rockscache.NewClient(rdb, rockscache.NewDefaultOptions())

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "dtm:passwd123dtm@tcp(en.dtm.pub:3306)/cache1?charset=utf8")
	logger.FatalIfError(err)
}

func clearData() {
	err := rdb.Del(rdb.Context(), DataKey).Err()
	logger.FatalIfError(err)
	_, err = db.Exec("delete from cache1.ver")
	logger.FatalIfError(err)
}

func initData(key string, value string, mode string) {
	clearData()
	SetCacheValue(key, value, mode)
	SetDBValue(&DBRow{
		K:        key,
		V:        value,
		TimeCost: "",
	})
}

func ensure(condition bool, format string, v ...interface{}) {
	hint := "ok"
	if !condition {
		hint = "failed"
	}
	logger.Infof("ensure: %s for %s", hint, fmt.Sprintf(format, v...))
	if !condition {
		os.Exit(1)
	}
}

const DataKey = "key1"

type DBRow struct {
	K        string
	V        string
	TimeCost string
}

func Post(url string, body map[string]interface{}) *resty.Response {
	logger.Infof("posting: %s, %s", url, body)
	r, err := resty.New().R().SetBody(body).Post(url)
	logger.FatalIfError(err)
	logger.FatalfIf(r.StatusCode() != 200, "post failed: %s", r.String())
	return r
}

func Get(url string) *resty.Response {
	logger.Infof("getting: %s, %s", url)
	r, err := resty.New().R().Get(url)
	logger.FatalIfError(err)
	logger.FatalfIf(r.StatusCode() != 200, "post failed: %s", r.String())
	return r
}

func QueryRow(db *sql.DB, sql string, args ...interface{}) *sql.Row {
	logger.Infof("query: %s %v", sql, args)
	return db.QueryRow(sql, args...)
}

func GetDBValue(key string) DBRow {
	var row DBRow
	r := QueryRow(db, "select k, v, time_cost from cache1.ver where k=?", key)
	err := r.Scan(&row.K, &row.V, &row.TimeCost)
	logger.FatalIfError(err)
	return row
}

func SetDBValue(row *DBRow) {
	tx, err := db.Begin()
	logger.FatalIfError(err)
	_, err = Exec(tx, "insert into cache1.ver(k, v, time_cost) values(?, ?, ?) on duplicate key update v=values(v), time_cost=values(time_cost)", row.K, row.V, row.TimeCost)
	logger.FatalIfError(err)
	err = tx.Commit()
	logger.FatalIfError(err)
}

func GetCacheValue(key string, mode string) string {
	var res string
	var err error
	if mode == "rockscache" {
		res, err = dc.RawGet(key)
	} else {
		res, err = rdb.Get(rdb.Context(), key).Result()
	}
	logger.Debugf("get cache: %s, %s, %v %s", key, res, err, mode)
	logger.FatalIfError(err)
	return res
}

func SetCacheValue(key string, value string, mode string) {
	var err error
	if mode == "rockscache" {
		err = dc.RawSet(key, value, 300*time.Second)
	} else {
		err = rdb.Set(rdb.Context(), key, value, 300*time.Second).Err()
	}
	logger.FatalIfError(err)
	logger.Debugf("set cache: %s, %s %s", key, value, mode)
}

func MustMapBodyFrom(c *gin.Context) map[string]interface{} {
	var body map[string]interface{}
	err := c.BindJSON(&body)
	logger.FatalIfError(err)
	return body
}

func Exec(tx *sql.Tx, sql string, args ...interface{}) (sql.Result, error) {
	logger.Infof("exec: %s %v", sql, args)
	return tx.Exec(sql, args...)
}

func UpdateInTx(tx *sql.Tx, row *DBRow) error {
	_, err := Exec(tx, "insert into cache1.ver(k, v, time_cost) values(?, ?, ?) on duplicate key update v=values(v), time_cost=values(time_cost)", row.K, row.V, row.TimeCost)
	return err
}

func Fetch(mode string, key string, expire time.Duration, fn func() (string, error)) (string, error) {
	if mode == "rockscache" {
		return dc.Fetch(key, expire, fn)
	} else {
		return NormalFetch(key, expire, fn)
	}
}

func NormalFetch(key string, expire time.Duration, fn func() (string, error)) (string, error) {
	res, err := rdb.Get(rdb.Context(), key).Result()
	if err == redis.Nil {
		res, err = fn()
		if err != nil {
			return "", err
		}
		logger.Debugf("set %s to %s %v", key, res, expire)
		_, err = rdb.Set(rdb.Context(), key, res, expire).Result()
	}
	return res, err
}

func DeleteCacheValue(key string) {
	err := rdb.Del(rdb.Context(), key).Err()
	logger.Debugf("delete cache: %s, %v", key, err)
	logger.FatalIfError(err)
}
