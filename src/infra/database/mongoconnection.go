package database

import (
	"context"
	"entity/src/apperrors"
	"entity/src/infra/logger"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
 * Type Implementation
 */
type mongoConn struct {
	once   sync.Once
	db     *mongo.Database
	params struct {
		dbName           string
		conStr           string
		conTimeOut       time.Duration
		upsertGetTimeOut time.Duration
	}
}

func (mc *mongoConn) Insert(data interface{}, table string) (id interface{}, err *apperrors.Error) {
	if cerr := initConn(); cerr == nil {
		if res, dberr := mc.db.Collection(table).InsertOne(getContext(mc.params.upsertGetTimeOut), data); dberr == nil {
			id = res.InsertedID
		} else {
			err = getDataError([2]string{"mongodb insert one error", dberr.Error()}, &dberr, data)
		}
	} else {
		err = getConnError(cerr)
	}

	return id, err
}

// Get a document from database and inject its content into result param, that SHOULD be an address reference (&result)
func (mc *mongoConn) Get(id interface{}, table string, result interface{}) (err *apperrors.Error) {
	if cerr := initConn(); cerr == nil {
		ctx := getContext(mc.params.upsertGetTimeOut)
		bid, _ := primitive.ObjectIDFromHex(id.(string))
		errmsg := [2]string{"mongodb findOne error", "error on get data"}

		if raw, ferr := mc.db.Collection(table).FindOne(ctx, bson.M{"_id": bid}).DecodeBytes(); ferr == nil {
			if rerr := setResult(raw, result); rerr != nil {
				err = getDataError(errmsg, &rerr, nil)
			}
		} else {
			if isNotFoundError(ferr) {
				err = getNotFoundError(id)
			} else {
				err = getDataError(errmsg, &ferr, nil)
			}
		}
	} else {
		err = getConnError(cerr)
	}

	return err
}

func NewDatabaseConnection() DataBaseHandler {
	return &singCon
}

/*
 * Type Auxiliar Funcs
 */
func getConnError(err error) *apperrors.Error {
	logger.Error("mongodb connection error", err)
	res := apperrors.NewInfraError("error trying to connect on database", nil)

	return &res
}

func getDataError(msgs [2]string, err *error, data interface{}) *apperrors.Error {
	logger.Error(msgs[0], *err)
	arr := apperrors.NewDataError(msgs[1], data)

	return &arr
}

func getNotFoundError(id interface{}) *apperrors.Error {
	arr := apperrors.NewNotFoundError("data with informed id not found", id)

	return &arr
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "no documents in result")
}

func setResult(raw bson.Raw, result interface{}) (err error) {
	var bys []byte
	if bys, err = bson.Marshal(raw); err == nil {
		err = bson.Unmarshal(bys, result)
	}

	return err
}

/*
 * Singleton Implementation
 */
var singCon mongoConn

func initConn() error {
	singCon.once.Do(func() {
		db := os.Getenv("MONGO_DB")
		if db == "" {
			log.Fatal("env MONGO_DB not informed")
		}

		conSt := os.Getenv("MONGO_CONN_STR")
		if conSt == "" {
			log.Fatal("env MONGO_CONN_STR not informed")
		}

		tmOut := func() int64 {
			var def int64
			def = 5

			if t := os.Getenv("MONGO_TIMEOUT"); t != "" {
				def, _ = strconv.ParseInt(t, 0, 0)
			}

			return def
		}()

		singCon.params = struct {
			dbName           string
			conStr           string
			conTimeOut       time.Duration
			upsertGetTimeOut time.Duration
		}{db, conSt, time.Duration(tmOut), 2}

	})

	return connect()
}

func connect() (err error) {
	if !isConnected() {
		opts := options.Client().ApplyURI(singCon.params.conStr)

		var cli *mongo.Client

		if cli, err = mongo.Connect(getContext(singCon.params.conTimeOut), opts); err == nil {
			if err = cli.Ping(getContext(2), nil); err == nil {
				logger.Info("database successfully connected")
				singCon.db = cli.Database(singCon.params.dbName)
			}
		}
	}

	return err
}

func isConnected() bool {
	return singCon.db != nil && func() bool {
		err := singCon.db.Client().Ping(getContext(2), nil)
		if err != nil {
			logger.Error("database check connection err", err)
		}

		return err == nil
	}()
}

func getContext(tmout time.Duration) (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), tmout*time.Second)
	return ctx
}
