# Une partie de ce README.md a été faire avec Copilot Chat<br />Donc si quelque chose ne fonctionne pas c'est la faute de Copilot<br />Toutes les routes ont été testé et fonctionne<br />Pour le front j'en avait marre donc il y a que la gestion de `poll`

Le dns de mongo db est: `mongodb://localhost:27017/`
Il faut juste faire une db nommée `poll` et le reste est automatique (normalement?)


# API Routes

## POST /poll

Creates a new poll.

**Request Body:**

```json
{
  "name": "<String>",
  "password": "<String>"
}
```

**Returns:**
- 200: object ID
- 400: name missing or empty
- 500: internal error while contacting back-end


## GET /poll/:id

Get a poll.

**Returns:**
- 200: object ID
- 400: invalid id
- 404: not found
- 500: internal error while contacting back-end


## PUT /poll

Update a poll name (only).

**Request Body:**

```json
{
    "id": "<String>",
    "name": "<String>",
    "password": "<String>"
}
```

**Returns:**
- 200: object ID
- 400: name missing or empty
- 401: wrong password
- 500: internal error while contacting back-end


## DELETE /poll

Delete a poll.

**Request Body:**

```json
{
    "id": "<String>",
    "password": "<String>"
}
```

**Returns:**
- 200: success
- 400: invalid id
- 401: wrong password
- 404: not found
- 500: internal error while contacting back-end




## POST /question

Creates a new question.

**Request Body:**

```json
{
  "poll_id": "<String>",
  "name": "<String>",
  "has_choice": "<Boolean>",
  "is_multiple": "<Boolean>"
}
```

**Returns:**
- 200: object ID
- 400: property missing or empty
- 500: internal error while contacting back-end


## GET /question/:id

Get a question.

**Returns:**
- 200: 
```json
{
  "id": "<String>",
  "poll_id": "<String>",
  "name": "<String>",
  "has_choice": "<Boolean>",
  "is_multiple": "<Boolean>"
}
```
- 400: invalid id
- 404: not found
- 500: internal error while contacting back-end


## PUT /question

Update a question name (only).
password is the poll password

**Request Body:**

```json
{
  "id": "<String>",
  "name": "<String>",
  "password": "<String>"
}
```

**Returns:**
- 200: object ID
- 400: name missing or empty
- 401: wrong password
- 500: internal error while contacting back-end


## DELETE /question

Delete a question.
password is the poll password

**Request Body:**

```json
{
    "id": "<String>",
    "password": "<String>"
}
```

**Returns:**
- 200: success
- 400: invalid id
- 401: wrong password
- 404: not found
- 500: internal error while contacting back-end

### POST /choice
Creates a new choice.

**Request Body:**
```json
{
  "question_id": "<ObjectId>",
  "name": "<String>",
  "password": "<String>"
}
```

**Returns:**
- 200: object ID
- 400: property missing or empty
- 500: internal error while contacting back-end

### GET /choice/:id
Get a choice.

**Returns:**
- 200: 
```json
{
  "id": "<ObjectId>",
  "question_id": "<ObjectId>",
  "name": "<String>"
}
```
- 400: invalid id
- 404: not found
- 500: internal error while contacting back-end

### PUT /choice
Update a choice name (only).
password is the poll password

**Request Body:**
```json
{
  "id": "<ObjectId>",
  "name": "<String>",
  "password": "<String>"
}
```

**Returns:**
- 200: object ID
- 400: name missing or empty
- 401: wrong password
- 500: internal error while contacting back-end

### DELETE /choice
Delete a choice.
password is the poll password

**Request Body:**
```json
{
    "id": "<ObjectId>",
    "password": "<String>"
}
```

**Returns:**
- 200: success
- 400: invalid id
- 401: wrong password
- 404: not found
- 500: internal error while contacting back-end



## POST /answer

Creates a new answer.

**Request Body:**

```json
{
  "question_id": "<ObjectId>",
  "choice_id": "<ObjectId>",
  "value": "<String>"
}
```

This is handled by the `Create` function in the [`AnswerHandler`](command:_github.copilot.openSymbolInFile?%5B%22pkg%2Fapi%2Fhandler%2Fanswer_handler.go%22%2C%22AnswerHandler%22%5D "pkg/api/handler/answer_handler.go") class.

## DELETE /answer

Deletes an existing answer.

**Request Body:**

```json
{
  "id": "<ObjectId>",
  "password": "<String>"
}
```

This is handled by the `Delete` function in the [`AnswerHandler`](command:_github.copilot.openSymbolInFile?%5B%22pkg%2Fapi%2Fhandler%2Fanswer_handler.go%22%2C%22AnswerHandler%22%5D "pkg/api/handler/answer_handler.go") class.

## GET /answer/:id

Gets an answer by id.

**Request Body:**

None

This is handled by the `Get` function in the [`AnswerHandler`](command:_github.copilot.openSymbolInFile?%5B%22pkg%2Fapi%2Fhandler%2Fanswer_handler.go%22%2C%22AnswerHandler%22%5D "pkg/api/handler/answer_handler.go") class.
