package handler

import (
	"atlas/pkg/models"
	"atlas/pkg/mongodb"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionHandler struct {
	mc *mongodb.MongoClient
	c  *mongodb.MongoClientCollection
}

func NewQuestionHandler(mc *mongodb.MongoClient, c *mongodb.MongoClientCollection) *QuestionHandler {
	return &QuestionHandler{mc, c}
}

func (h *QuestionHandler) Get(ctx *fiber.Ctx) error {
	question := &models.Question{}
	id := ctx.Params("id")
	if len(id) < 1 {
		return fiber.NewError(http.StatusBadRequest, "null id provided")
	}

	err := h.c.FindID(id, question)
	if err != nil {
		return fiber.ErrNotFound
	}

	return ctx.JSON(question)
}

func (h *QuestionHandler) Create(ctx *fiber.Ctx) error {
	var req models.Question

	// get request to req struct
	err := ctx.BodyParser(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}

	// verify our properties are valid
	if len(req.Name) < 1 {
		return fiber.NewError(http.StatusBadRequest, "empty question name provided")
	} else if req.PollID == primitive.NilObjectID {
		return fiber.NewError(http.StatusBadRequest, "null question id provided")
	}

	// check if poll exists
	err = h.c.Database().Collection("poll").FindObjectID(req.PollID, &models.Poll{})
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid poll id provided")
	}

	if !req.HasChoice && req.IsMultiple {
		req.IsMultiple = false
	}

	// insert request
	req.ID = primitive.NewObjectID()
	result, err := h.c.InsertOne(req)

	// if there is an error return code 500
	if err != nil {
		log.Println("Question -> InsertOne: ", err)
		return fiber.ErrInternalServerError
	}

	// else return the inserted ID
	return ctx.JSON(result.InsertedID)
}

func (h *QuestionHandler) Update(ctx *fiber.Ctx) error {
	var req struct {
		ID         primitive.ObjectID `json:"id"`
		Name       string             `json:"name"`
		IsMultiple bool               `json:"is_multiple"`
		Password   string             `json:"password"`
	}

	// get request to req struct
	err := ctx.BodyParser(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}

	// check if name is not empty
	if len(req.Name) < 1 {
		return fiber.NewError(http.StatusBadRequest, "The name of the choice cannot be empty")
	}

	question := models.Question{}
	err = h.c.FindObjectID(req.ID, &question)
	if err == nil {
		poll := models.Poll{}
		err = h.c.Database().Collection("poll").FindObjectID(question.PollID, &poll)
		// if found and password invalid
		if err == nil && req.Password != poll.Password {
			return fiber.ErrUnauthorized
		}
		// else if found and password match then authorized
		// or not found then authorized by default
	}
	// if question is null it means there is no poll -> authorized

	question.Name = req.Name
	question.IsMultiple = req.IsMultiple
	return h.c.UpdateOne(bson.M{"_id": question.ID}, question)
}

func (h *QuestionHandler) Delete(ctx *fiber.Ctx) error {
	var req struct {
		ID       primitive.ObjectID
		Password string
	}

	// get request to req struct
	err := ctx.BodyParser(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}

	// check if the question id is not null
	if req.ID == primitive.NilObjectID {
		return fiber.NewError(http.StatusBadRequest, "null question id provided")
	}

	question := models.Question{}
	err = h.c.FindObjectID(req.ID, &question)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid question id provided")
	}

	poll := models.Poll{}
	err = h.c.Database().Collection("poll").FindObjectID(question.PollID, &poll)
	// if found and password invalid
	if err == nil && req.Password != poll.Password {
		return fiber.ErrUnauthorized
	}
	// else if found and password match then authorized
	// or not found then authorized by default
	h.c.Database().Collection("answer").DeleteMany(bson.M{"question_id": req.ID})
	h.c.Database().Collection("choice").DeleteMany(bson.M{"question_id": req.ID})

	return h.c.DeleteOne(bson.M{"_id": req.ID})
}
