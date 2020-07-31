package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

func BatchBulkCreate(BulkCreate func(*[]bson.M) error, documents *[]bson.M, limit int) error {
	counter := 0
	totalDocuments := len(*documents)

	for counter != totalDocuments {
		var nElements int

		if counter+limit > totalDocuments {
			nElements = totalDocuments - counter
		} else {
			nElements = limit
		}

		docsToCreate := make([]bson.M, nElements)

		copy(docsToCreate, (*documents)[counter:counter+nElements])

		err := BulkCreate(&docsToCreate)

		if err != nil {
			return err
		}

		counter += len(docsToCreate)
	}

	return nil
}
