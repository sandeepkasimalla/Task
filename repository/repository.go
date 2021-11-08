package repository

import (
	"context"
	"errors"
	"time"

	"github.com/go-chassis/openlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersRepo struct {
	DbClient     *mongo.Client
	DatabaseName string
}

//Insert function will insert the data into database
func (dr *UsersRepo) Insert(meta map[string]interface{}) (map[string]interface{}, int, error) {
	collection := dr.DbClient.Database(dr.DatabaseName).Collection("users")

	meta["createdon"] = time.Now().Unix()
	res, err := collection.InsertOne(context.Background(), meta)
	if err != nil {
		openlog.Error(err.Error())
		return make(map[string]interface{}), 500, errors.New("Internal Server Error")
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	result, _, err := dr.Find(id)
	if err != nil {
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}

//Find function will find the document by id returns that document
func (dr *UsersRepo) Find(id string) (map[string]interface{}, int, error) {
	collection := dr.DbClient.Database(dr.DatabaseName).Collection("users")
	result := make(map[string]interface{})
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		openlog.Error(err.Error())
		return result, 400, errors.New("Invalid Id")
	}
	err = collection.FindOne(context.TODO(), bson.M{"_id": docID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, 404, errors.New("No Documents Found")
		}
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}

//FindAndUpdate function will find the documentation by id and update the data for that document
func (dr *UsersRepo) FindAndUpdate(id string, document map[string]interface{}) (map[string]interface{}, int, error) {
	collection := dr.DbClient.Database(dr.DatabaseName).Collection("users")
	update := make(map[string]interface{})
	update["$set"] = document
	result := make(map[string]interface{})
	docID, convErr := primitive.ObjectIDFromHex(id)
	if convErr != nil {
		openlog.Error(convErr.Error())
		return result, 400, errors.New("Invalid ID")
	}
	document["updatedon"] = time.Now().Unix()
	err := collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": docID}, update).Decode(&result)
	if err != nil {
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	result, _, err = dr.Find(id)
	if err != nil {
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}

//FindByFiltersAndPagenation will find the document by filters and pagenation and returns the data
func (dr *UsersRepo) FindByFiltersAndPagenation(page int64, size int64, filters map[string]interface{}, sort map[string]interface{}) ([]map[string]interface{}, int64, error) {

	options := *options.Find()

	collection := dr.DbClient.Database(dr.DatabaseName).Collection("users")
	var result []map[string]interface{}

	total, err := collection.CountDocuments(context.TODO(), filters)
	if err != nil {
		openlog.Error(err.Error())
		return result, total, errors.New("Internal Server Error")
	}
	cursor, err := collection.Find(context.TODO(), filters, options.SetLimit(size), options.SetSkip((page-1)*size), options.SetSort(sort))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, total, nil
		}
		openlog.Error(err.Error())
		return result, total, errors.New("Internal Server Error")
	}
	for cursor.Next(context.TODO()) {
		var doc map[string]interface{}
		err := cursor.Decode(&doc)
		if err != nil {
			openlog.Error(err.Error())
			return result, total, errors.New("Internal Server Error")
		}
		result = append(result, doc)
	}
	return result, total, nil
}

//FindByFilters will find the document by filters and pagenation and returns the data
func (dr *UsersRepo) FindByFilters(filters map[string]interface{}, sort map[string]interface{}) ([]map[string]interface{}, int, error) {

	options := *options.Find()

	collection := dr.DbClient.Database(dr.DatabaseName).Collection("users")
	var result []map[string]interface{}

	cursor, err := collection.Find(context.TODO(), filters, options.SetSort(sort))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, 200, nil
		}
		return result, 500, errors.New("Internal server error")
	}
	for cursor.Next(context.TODO()) {
		var doc map[string]interface{}
		err := cursor.Decode(&doc)
		if err != nil {
			openlog.Error(err.Error())
			return result, 500, errors.New("Internal Server Error")
		}
		result = append(result, doc)
	}
	return result, 200, nil
}

//IsBoilerPlateNameExists checks for any Boilerplate conflict
func (dr *UsersRepo) IsEmailExists(email string) (int, error) {
	collection := dr.DbClient.Database(dr.DatabaseName).Collection("users")
	result := make(map[string]interface{})
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 404, nil
		}
		openlog.Error(err.Error())
		return 500, errors.New("Internal Server Error")
	}
	return 409, errors.New("Email already exists")
}

//Delete will allows to delete the  document
func (dr *UsersRepo) Delete(id string) (map[string]interface{}, int, error) {
	collection := dr.DbClient.Database(dr.DatabaseName).Collection("users")
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		openlog.Error(err.Error())
		return nil, 400, errors.New("Invalid id")
	}
	result := make(map[string]interface{})
	err = collection.FindOneAndDelete(context.TODO(), bson.M{"_id": docID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, 404, errors.New("No Documents Found for given id")
		}
		openlog.Error(err.Error())
		return nil, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}
