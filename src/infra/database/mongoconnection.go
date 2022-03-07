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
	if err = initConn(); err == nil {
		if res, dberr := mc.db.Collection(table).InsertOne(getContext(mc.params.upsertGetTimeOut), data); dberr == nil {
			id = res.InsertedID
		} else {
			err = getDataError([2]string{"mongodb insert one error", dberr.Error()}, dberr, data)
		}
	}

	return
}

// Get a document from database and inject its content into result param, that SHOULD be an address reference (&result)
func (mc *mongoConn) Get(id interface{}, table string, result interface{}, filters ...map[string]interface{}) (err *apperrors.Error) {
	if err = initConn(); err == nil {
		fils := getFiltersWithId(id, filters...)
		ctx := getContext(mc.params.upsertGetTimeOut)
		errmsg := [2]string{"mongodb findOne error", "error on get data"}

		if raw, ferr := mc.db.Collection(table).FindOne(ctx, fils).DecodeBytes(); ferr == nil {
			err = setResult(raw, result, errmsg)
		} else {
			err = getError(ferr, id.(string), errmsg)
		}
	}

	return
}

// Updates a document into database and inject its new content into result param, that SHOULD be an address reference (&result)
func (mc *mongoConn) Update(id interface{}, data interface{}, table string, result interface{}, filters ...map[string]interface{}) (err *apperrors.Error) {
	if err = initConn(); err == nil {
		fils := getFiltersWithId(id, filters...)
		ctx := getContext(mc.params.upsertGetTimeOut)
		errmsg := [2]string{"mongodb findOneAndUpdate error", "error on update data"}
		updoc := bson.M{"$set": data}
		afopt := options.After
		upsopt := false
		opt := options.FindOneAndUpdateOptions{ReturnDocument: &afopt, Upsert: &upsopt}

		if raw, uerr := mc.db.Collection(table).FindOneAndUpdate(ctx, fils, updoc, &opt).DecodeBytes(); uerr == nil {
			err = setResult(raw, result, errmsg)
		} else {
			err = getError(uerr, id.(string), errmsg)
		}
	}

	return
}

func NewDatabaseConnection() DataBaseHandler {
	return &singCon
}

/*
 * Type Auxiliar Funcs
 */
func getFilters(fils ...map[string]interface{}) (res primitive.D) {
	for _, fil := range fils {
		for k, v := range fil {
			res = append(res, bson.E{Key: k, Value: v})
		}
	}

	return
}

func getFiltersWithId(id interface{}, fils ...map[string]interface{}) primitive.D {
	bid, _ := primitive.ObjectIDFromHex(id.(string))
	fils = append(fils, map[string]interface{}{"_id": bid})

	return getFilters(fils...)
}

func getError(err error, id string, msgs [2]string) *apperrors.Error {
	if isNotFoundError(err) {
		return getNotFoundError(id)
	}

	return getDataError(msgs, err, nil)
}

func getDataError(msgs [2]string, err error, data interface{}) *apperrors.Error {
	logger.Error(msgs[0], err)
	res := apperrors.NewDataError(msgs[1], data)

	return &res
}

func getNotFoundError(id interface{}) *apperrors.Error {
	res := apperrors.NewNotFoundError("data with informed id not found", id)

	return &res

}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "no documents in result")
}

func setResult(raw bson.Raw, result interface{}, msgs [2]string) (err *apperrors.Error) {
	var bys []byte
	var perr error

	if bys, perr = bson.Marshal(raw); perr == nil {
		perr = bson.Unmarshal(bys, result)
	}

	if perr != nil {
		err = getDataError(msgs, perr, nil)
	}

	return
}

/*
 * Singleton Implementation
 */
var singCon mongoConn

func initConn() *apperrors.Error {
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

func connect() (aperr *apperrors.Error) {
	if !isConnected() {
		opts := options.Client().ApplyURI(singCon.params.conStr)
		var cli *mongo.Client
		var err error

		if cli, err = mongo.Connect(getContext(singCon.params.conTimeOut), opts); err == nil {
			if err = cli.Ping(getContext(2), nil); err == nil {
				logger.Info("database successfully connected")
				singCon.db = cli.Database(singCon.params.dbName)
			}
		}

		if err != nil {
			logger.Error("mongodb connection error", err)
			res := apperrors.NewInfraError("error trying to connect on database", nil)
			aperr = &res
		}
	}

	return
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
