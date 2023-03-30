package controllers

import (
	"context"
	"gin-mongodb/configs"
	"gin-mongodb/models"
	"gin-mongodb/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	// validation request body
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data: map[string]interface{}{
				"error": err.Error()},
		})
		return
	}

	// validate required fields
	if validationError := validate.Struct(&user); validationError != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data: map[string]interface{}{
				"error": validationError.Error()},
		})
		return
	}

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	_, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"error": err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, response.UserResponse{
		Status:  http.StatusCreated,
		Message: "success"})
}

func GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.Param("userId")
	var responseUser models.User

	objId, _ := primitive.ObjectIDFromHex(userId)
	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&responseUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: map[string]interface{}{
				"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, response.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &response.UserResponseData{
			User: responseUser,
		},
	})
}

func EditUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.Param("userId")
	var requestBodyUserData models.User
	objId, _ := primitive.ObjectIDFromHex(userId)

	// validate request body
	if err := c.BindJSON(&requestBodyUserData); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data: map[string]interface{}{
				"error": err.Error()},
		})
		return
	}

	// validate required fields
	if validationError := validate.Struct(&requestBodyUserData); validationError != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data: map[string]interface{}{
				"error": validationError.Error()},
		})
		return
	}

	updateData := bson.M{
		"name":     requestBodyUserData.Name,
		"location": requestBodyUserData.Location,
		"title":    requestBodyUserData.Title,
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": updateData})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data: map[string]interface{}{
				"error": err.Error()},
		})
		return
	}

	// get updated user detail
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error",
				Data: map[string]interface{}{
					"error": err.Error()},
			})
			return
		}
	}

	c.JSON(http.StatusOK, response.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &response.UserResponseData{
			User: updatedUser,
		},
	})
}

func DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.Param("userId")
	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data: map[string]interface{}{
				"error": err.Error()},
		})
		return
	}

	if result.DeletedCount < 1 {
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "Error",
			Data: map[string]interface{}{
				"error": "user with specified ID not found"},
		})
		return
	}
	c.JSON(http.StatusCreated, response.UserResponse{
		Status:  http.StatusCreated,
		Message: "success"})
}

func GetAllUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User
	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data: map[string]interface{}{
				"error": err.Error()},
		})
		return
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error",
				Data: map[string]interface{}{
					"error": err.Error()},
			})
			return
		}

		users = append(users, singleUser)
	}

	c.JSON(http.StatusOK,
		response.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: &response.UserResponseData{
				User: users,
			}},
	)
}
