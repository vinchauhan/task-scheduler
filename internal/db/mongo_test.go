package db_test

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vinchauhan/task-scheduler/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)
var mongoClient *mongo.Client
func TestPlug(t *testing.T) {
	Convey("Setup", t, func() {
		mongoDb := db.New("mongodb://localhost:9999")
		badMongoDb := db.New("mongodb://")
		Convey("On a valid url , Plug method will return a mongo connection handle", func() {
			con, _ := mongoDb.Plug()
			So(con, ShouldHaveSameTypeAs, mongoClient)
		})

		Convey("On a invalid url provided , Plug method will return an error", func() {
			_, err := badMongoDb.Plug()
			So(err.Error(), ShouldResemble, errors.New("error parsing uri: must have at least 1 host").Error())
		})

	})
}
