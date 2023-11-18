package pipeline

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestPaginate(t *testing.T) {
	from := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)
	limit := 10
	offset := 5

	t.Run("WithToNull", func(t *testing.T) {
		expectedPipeline := mongo.Pipeline{
			{
				{"$match", bson.D{
					{"created_at", bson.D{
						{"$gte", from},
					}},
				}},
			},
			{
				{"$sort", bson.D{
					{"created_at", 1},
				}},
			},
			{
				{"$limit", 10},
			},
			{
				{"$skip", 5},
			},
		}

		pipeline := Paginate(from, nil, limit, offset)

		assert.Equal(t, expectedPipeline, pipeline)
	})

	t.Run("WithToValue", func(t *testing.T) {
		expectedPipeline := mongo.Pipeline{
			{
				{"$match", bson.D{
					{"created_at", bson.D{
						{"$gte", from},
						{"$lt", to},
					}},
				}},
			},
			{
				{"$sort", bson.D{
					{"created_at", 1},
				}},
			},
			{
				{"$limit", 10},
			},
			{
				{"$skip", 5},
			},
		}

		pipeline := Paginate(from, &to, limit, offset)

		assert.Equal(t, expectedPipeline, pipeline)
	})

	t.Run("ZeroLimitAndOffset", func(t *testing.T) {
		expectedPipeline := mongo.Pipeline{
			{
				{"$match", bson.D{
					{"created_at", bson.D{
						{"$gte", from},
						{"$lt", to},
					}},
				}},
			},
			{
				{"$sort", bson.D{
					{"created_at", 1},
				}},
			},
			{
				{"$skip", 0},
			},
		}

		pipeline := Paginate(from, &to, 0, 0)

		assert.EqualValues(t, expectedPipeline, pipeline)
	})
}
