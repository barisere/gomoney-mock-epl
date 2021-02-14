package database

import "go.mongodb.org/mongo-driver/mongo"

func IsDuplicateKeyError(err error) bool {
	writeException, ok := err.(mongo.WriteException)
	if !ok {
		return false
	}
	return anyWriteError(writeException.WriteErrors)
}

func anyWriteError(errs mongo.WriteErrors) bool {
	for _, err := range errs {
		if err.Code == 11000 {
			return true
		}
	}
	return false
}
