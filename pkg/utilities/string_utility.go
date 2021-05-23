package utilities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

func ToObjectID(id string) (primitive.ObjectID, error){
	if objectId, err := primitive.ObjectIDFromHex(id); err != nil{
		return primitive.ObjectID{}, err
	}else{
		return objectId, nil
	}
}

func IsStringEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
