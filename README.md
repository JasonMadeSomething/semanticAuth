# Semantic Authenticator

Semantic Authenticator is an experiment in using OpenAI’s embedding models to power a cosine-similarity-based login system. Instead of matching passwords exactly, it matches them semantically.

---

## What It Does

* Users register with any string — a phrase, a memory, a feeling.
* That string is embedded via OpenAI and stored as a vector in MongoDB.
* During login, user input is embedded again and compared via cosine similarity.
* If the guess is *close enough* to the original, you get in.
* All login attempts are logged with similarity scores

---

## Why It Exists

This started as a joke about \*\*semantic security\*\*, but it turned into an actual working system with:

* Vector caching for efficiency (via MongoDB)
* Configurable similarity threshold per login attempt
* Full audit logging of inputs and scores
* A clean, minimal Go backend ready for deployment

---

## Tech Stack

* **Go**
* **MongoDB** (via Docker)
* **OpenAI API** for `text-embedding-3-small`
* **Cosine Similarity** for the actual login math
* **Chi** router with CORS for frontend integration
* **Postman-tested REST API** (Frontend WIP)

---

## Endpoints

### `POST /register`

Registers a new user.

#### Request Body

```json
{
  "username": "steve",
  "password": "my grandma’s lasagna recipe"
}
```

---

### `POST /login`

Logs in by comparing the semantic similarity to the stored vector.

#### Request Body

```json
{
  "username": "steve",
  "password": "lasagna recipe",
  "threshold": 0.88
}
```

---

### `POST /report`

Fetches recent login attempts for a user.

#### Request Body

```json
{
  "username": "steve",
  "threshold": 0.88
}
```

#### Example Response

```json
[
  {
    "input": "lasagna recipe",
    "similarity": 0.645,
    "timestamp": "2025-07-26T01:42:36.137Z",
    "passed": false
  },
  {
    "input": "grandma’s lasagna",
    "similarity": 0.899,
    "timestamp": "2025-07-26T01:42:05.547Z",
    "passed": true
  }
]
```

---

## Setup (Dev)

```bash
git clone https://github.com/JasonMadeSomething/semanticAuth.git
cd semantic-auth

go mod tidy
cp .env.example .env  # Add your OpenAI key
docker-compose up -d

go run main.go
```

---

## Sample Playground Inputs

* "My first pet's favorite jazz song"
* "The thing she said under the stars"
* "Fear. But the good kind."
* "Tuesdays at grandma’s"

---

---

## Disclaimers

This is not production-grade security. It’s a tech demo / philosophical experiment/art piece in login form.

---

## Contact

Built by @jasonmadesomething&#x20;
