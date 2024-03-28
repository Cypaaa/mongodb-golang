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

type ChoiceHandler struct {
	mc *mongodb.MongoClient
	c  *mongodb.MongoClientCollection
}

func NewChoiceHandler(mc *mongodb.MongoClient, c *mongodb.MongoClientCollection) *ChoiceHandler {
	return &ChoiceHandler{mc, c}
}

func (h *ChoiceHandler) Get(ctx *fiber.Ctx) error {
	choice := &models.Choice{}
	// get id
	id := ctx.Params("id")
	// check if id is not null
	if len(id) < 1 {
		return fiber.NewError(http.StatusBadRequest, "null id provided")
	}

	// get the choice from its id
	err := h.c.FindID(id, choice)
	if err != nil {
		return fiber.ErrNotFound
	}

	// return it
	return ctx.JSON(choice)
}

func (h *ChoiceHandler) Create(ctx *fiber.Ctx) error {
	var req models.Choice

	// get request to req struct
	err := ctx.BodyParser(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}

	// check if name is not empty
	if len(req.Name) < 1 {
		return fiber.NewError(http.StatusBadRequest, "The name of the choice cannot be empty")
	} else if req.QuestionID == primitive.NilObjectID {
		return fiber.NewError(http.StatusBadRequest, "null question id provided")
	}

	// check if question exists AND has_choice == true
	question := &models.Question{}
	err = h.c.Database().Collection("question").FindObjectID(req.QuestionID, question)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid question id provided")
	} else if !question.HasChoice {
		return fiber.NewError(http.StatusBadRequest, "question has_choice is false")
	}

	// insert request
	req.ID = primitive.NewObjectID()
	result, err := h.c.InsertOne(&req)

	// if there is an error return code 500
	if err != nil {
		log.Println("Choice -> InsertOne: ", err)
		return fiber.ErrInternalServerError
	}

	// else return the inserted ID
	return ctx.JSON(result.InsertedID)
}

func (h *ChoiceHandler) Update(ctx *fiber.Ctx) error {
	var req struct {
		ID       primitive.ObjectID `json:"id"`
		Name     string             `json:"name"`
		Password string             `json:"password"`
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

	// get the choice from its id
	choice := models.Choice{}
	err = h.c.FindObjectID(req.ID, &choice)
	if err != nil {
		return fiber.ErrNotFound
	}

	question := models.Question{}
	err = h.c.FindObjectID(choice.QuestionID, &question)
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

	choice.Name = req.Name
	return h.c.UpdateOne(bson.M{"_id": choice.ID}, choice)
}

func (h *ChoiceHandler) Delete(ctx *fiber.Ctx) error {
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
		return fiber.NewError(http.StatusBadRequest, "null choice id provided")
	}

	choice := models.Choice{}
	err = h.c.FindObjectID(req.ID, &choice)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid choice id provided")
	}

	question := models.Question{}
	err = h.c.FindObjectID(choice.QuestionID, &question)
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

	h.c.Database().Collection("answer").DeleteMany(bson.M{"choice_id": req.ID})

	return h.c.DeleteOne(bson.M{"_id": req.ID})
}
