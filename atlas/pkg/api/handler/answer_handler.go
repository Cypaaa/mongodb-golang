package handler

import (
	"atlas/pkg/models"
	"atlas/pkg/mongodb"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnswerHandler struct {
	mc *mongodb.MongoClient
	c  *mongodb.MongoClientCollection
}

func NewAnswerHandler(mc *mongodb.MongoClient, c *mongodb.MongoClientCollection) *AnswerHandler {
	return &AnswerHandler{mc, c}
}

func (h *AnswerHandler) Get(ctx *fiber.Ctx) error {
	answer := &models.Answer{}
	id := ctx.Params("id")
	if len(id) < 1 {
		return fiber.NewError(http.StatusBadRequest, "null id provided")
	}

	err := h.c.FindID(id, answer)
	if err != nil {
		return fiber.ErrNotFound
	}

	return ctx.JSON(answer)
}

func (h *AnswerHandler) Create(ctx *fiber.Ctx) error {
	var req models.Answer

	// get request to req struct
	err := ctx.BodyParser(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}

	// check if the question id is not null
	if req.QuestionID == primitive.NilObjectID {
		return fiber.NewError(http.StatusBadRequest, "null question id provided")
	}

	// check if question exists
	question := models.Question{}

	err = h.c.Database().Collection("question").FindObjectID(req.QuestionID, &question)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid question id provided")
	}

	// if question has choice but no choice is selected
	if question.HasChoice && req.ChoiceID == primitive.NilObjectID {
		return fiber.NewError(http.StatusBadRequest, "null choice id provided")
	}

	// if question has no choice but a choice is selected
	if !question.HasChoice && req.ChoiceID != primitive.NilObjectID {
		return fiber.NewError(http.StatusBadRequest, "question has_choice is false")
	}

	// if question is not multiple choice and a choice is already in mongodb
	if !question.IsMultiple {
		answer := models.Answer{}
		err = h.c.Database().Collection("answer").FindOne(bson.M{"question_id": req.QuestionID}, &answer)
		if err == nil {
			return fiber.NewError(http.StatusBadRequest, "question is not multiple choice")
		}
	}

	// check if there is a choice id
	if req.ChoiceID != primitive.NilObjectID {
		// if so check if the choice exists
		choice := models.Choice{}
		err = h.c.Database().Collection("choice").FindObjectID(req.ChoiceID, &choice)
		if err != nil {
			return fiber.NewError(http.StatusBadRequest, "invalid choice id provided")
		}

		// check if the choice belongs to the question
		if choice.QuestionID != req.QuestionID {
			return fiber.NewError(http.StatusBadRequest, "choice does not belong to the question")
		}

		req.Value = choice.Name
	} else if len(req.Value) < 1 {
		return fiber.NewError(http.StatusBadRequest, "null value provided")
	}

	// insert request
	req.ID = primitive.NewObjectID()
	result, err := h.c.InsertOne(req)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(result.InsertedID)
}

func (h *AnswerHandler) Delete(ctx *fiber.Ctx) error {
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

	// get answer
	answer := models.Answer{}
	err = h.c.FindObjectID(req.ID, &answer)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid answer id provided")
	}

	question := models.Question{}
	err = h.c.FindObjectID(answer.QuestionID, &question)
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

	return h.c.DeleteOne(bson.M{"_id": req.ID})
}
