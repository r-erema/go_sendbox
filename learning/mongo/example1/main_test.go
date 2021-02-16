package example1

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
)

var db *mongo.Database

func TestMain(m *testing.M) {
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://erema:LUNJTG30LcpedGXz@cluster0.unmvy.mongodb.net/test?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database("test")
	if err = db.Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	code := m.Run()
	os.Exit(code)
}

func TestInsertResult(t *testing.T) {
	users := db.Collection("users")
	fixtures := fixtures()
	expectedDocsCount := len(fixtures)
	insertResult, err := users.InsertMany(context.TODO(), fixtures)
	require.NoError(t, err)
	assert.Equal(t, expectedDocsCount, len(insertResult.InsertedIDs))

	selectResult, err := users.Find(context.Background(), bson.D{})
	require.NoError(t, err)
	assert.Equal(t, expectedDocsCount, selectResult.RemainingBatchLength())
}

func fixtures() (fixtures []interface{}) {
	for i := 1; i <= 5; i++ {
		fixtures = append(fixtures, bson.D{
			{Key: "id", Value: 100},
			{Key: "balance", Value: i},
		})
	}
	for i := 70; i <= 195; i += 5 {
		fixtures = append(fixtures, bson.D{
			{Key: "id", Value: 200},
			{Key: "balance", Value: i},
		})
	}
	for i := 20; i <= 22; i++ {
		fixtures = append(fixtures, bson.D{
			{Key: "id", Value: 300},
			{Key: "balance", Value: i},
		})
	}
	for i := 210; i <= 270; i += 10 {
		fixtures = append(fixtures, bson.D{
			{Key: "id", Value: 400},
			{Key: "balance", Value: i},
		})
	}
	return
}
