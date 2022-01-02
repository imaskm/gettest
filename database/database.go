package database

import (
	"context"
	"log"
	"os"

	"github.com/imaskm/getir/types"
	"github.com/imaskm/getir/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	conn *mongo.Client
}

const (
	DatabaseName         = "getir-case-study"
	RecordCollectionName = "records"
)

var mongoURI string

func init() {

	mongoURI = os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("mongo db uri is required, export MONGO_URI")
	}
}

func NewMongoSession() (*MongoDB, error) {
	client, err := getMongodbclient()
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		// Can't connect to Mongo server
		log.Fatal(err)
	}
	m := new(MongoDB)
	m.conn = client

	return m, nil

}

func getMongodbclient() (*mongo.Client, error) {
	mgoOptns := options.Client().ApplyURI(mongoURI)

	client, err := mongo.NewClient(mgoOptns)
	if err != nil {
		log.Fatal(err)
	}
	client.Connect(context.Background())

	return client, nil
}

// GetRecords returns all the record filtered on the request
func (m *MongoDB) GetRecords(rr *types.RecordRequest) ([]*types.Record, error) {
	startDate, err := utility.ConvertStringDateToTime(rr.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := utility.ConvertStringDateToTime(rr.EndDate)
	if err != nil {
		return nil, err
	}

	f := bson.M{
		"createdAt": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	curr, err := m.conn.Database(DatabaseName).Collection(RecordCollectionName).Find(context.Background(), f)
	if err != nil {
		return nil, err
	}

	var records []*types.Record
	var record = new(types.RecordDb)

	defer curr.Close(context.Background())
	for curr.Next(context.Background()) {
		err = curr.Decode(record)
		if err != nil {
			log.Fatal(err)
		}

		s := utility.GetSum(record.Counts)
		if s >= rr.MinCount && s <= rr.MaxCount {
			r := new(types.Record)
			r.CreatedAt = record.CreatedAt
			r.Key = record.Key
			r.TotalCount = s
			records = append(records, r)
		}

	}

	if err := curr.Err(); err != nil {
		log.Fatal(err)
	}

	return records, nil
}
