package pipeline

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type SortOrder int

const (
	Asc  SortOrder = 1
	Desc           = -1
)

func Paginate(from time.Time, to *time.Time, lim, os int) mongo.Pipeline {
	return mongo.Pipeline{
		match(from, to),
		sortBy(Asc),
		offset(os),
		limit(lim),
	}
}

func match(from time.Time, to *time.Time) bson.D {
	span := bson.D{{"$gte", from}}
	if to != nil {
		span = append(span, bson.E{Key: "$lt", Value: *to})
	}
	return bson.D{
		{"$match",
			bson.D{
				{"created_at",
					span,
				},
			},
		},
	}
}

func sortBy(sortBy SortOrder) bson.D {
	return bson.D{{"$sort", bson.D{{"created_at", sortBy}}}}
}

func limit(l int) bson.D {
	return bson.D{{"$limit", l}}

}

func offset(o int) bson.D {
	return bson.D{{"$skip", o}}
}
