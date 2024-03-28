package handler

import (
	"atlas/pkg/models"
	"atlas/pkg/mongodb"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PollHandler struct {
	mc *mongodb.MongoClient
	c  *mongodb.MongoClientCollection
}

func NewPollHandler(mc *mongodb.MongoClient, c *mongodb.MongoClientCollection) *PollHandler {
	return &PollHandler{mc: mc, c: c}
}

// Get by id
func (h *PollHandler) Get(ctx *fiber.Ctx) error {
	poll := &models.Poll{}
	id := ctx.Params("id")
	if len(id) < 1 {
		return fiber.NewError(http.StatusBadRequest, "id cannot be empty")
	}

	err := h.c.FindID(id, poll)
	if err != nil {
		return fiber.ErrNotFound
	}

	poll.Password = ""
	b, err := json.Marshal(poll)
	if err != nil {
		return fiber.ErrBadRequest
	}

	_, err = ctx.Write(b)
	return err
}

func (h *PollHandler) Create(ctx *fiber.Ctx) error {
	var req models.Poll

	// if no error on parse
	// AND name is not nil
	err := ctx.BodyParser(&req)
	if err != nil || len(req.Name) < 1 {
		return fiber.ErrBadRequest
	}

	// password if not required so the poll can be edited by everyone
	req.ID = primitive.NewObjectID()
	result, err := h.c.InsertOne(req)
	if err != nil {
		log.Println("Poll -> InsertOne: ", err)
		if mongo.IsDuplicateKeyError(err) {
			return fiber.ErrConflict
		}
		log.Println("[ERROR](Create Poll)", err)
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(result.InsertedID)
}

func (h *PollHandler) Update(ctx *fiber.Ctx) error {
	var req struct {
		ID       primitive.ObjectID `json:"id"`
		Name     string             `json:"name"`
		Password string             `json:"password"`
	}

	// get request to req struct
	err := json.Unmarshal(ctx.Body(), &req)
	if err != nil {
		return fiber.ErrBadRequest
	}

	// check if name is not empty
	if len(req.Name) < 1 {
		return fiber.NewError(http.StatusBadRequest, "The name of the choice cannot be empty")
	}

	poll := models.Poll{}
	err = h.c.Database().Collection("poll").FindObjectID(req.ID, &poll)
	// if found and password invalid
	if err == nil && req.Password != poll.Password {
		return fiber.ErrUnauthorized
	} else if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id provided")
	}
	// else if found and password match then authorized
	// or not found then authorized by default

	err = h.c.UpdateOne(bson.M{"_id": poll.ID}, bson.D{{Key: "name", Value: req.Name}})
	return err
}

func (h *PollHandler) Delete(ctx *fiber.Ctx) error {
	var req struct {
		ID       primitive.ObjectID
		Password string
	}

	b, _ := io.ReadAll(bytes.NewBuffer(ctx.Body()))

	// get request to req struct
	err := json.Unmarshal(b, &req)
	if err != nil {
		return fiber.ErrBadRequest
	}

	// check if the question id is not null
	if req.ID == primitive.NilObjectID {
		return fiber.NewError(http.StatusBadRequest, "null question id provided")
	}

	poll := models.Poll{}
	err = h.c.Database().Collection("poll").FindObjectID(req.ID, &poll)
	// if not found then doesn't exist
	if err != nil {
		return nil
		// if found and password valid
	} else if req.Password == poll.Password {
		questions := []models.Question{}
		h.c.Find(bson.M{"poll_id": req.ID}, &questions)
		for _, q := range questions {
			h.c.Database().Collection("answer").DeleteMany(bson.M{"question_id": q.ID})
			h.c.Database().Collection("choice").DeleteMany(bson.M{"question_id": q.ID})
			h.c.Database().Collection("question").DeleteOne(bson.M{"_id": q.ID})
		}
		return h.c.DeleteOne(bson.M{"_id": req.ID})
	}
	// else unauthorized

	return fiber.ErrUnauthorized
}
