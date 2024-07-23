# GoQuiz

This is a simple Quiz game which allows the user to create a new quiz game with N questions, answer each question and get to know their results. Users can also see how well they were compared to each other in form of a leaderboard

The project is composed by a **REST API** and a **CLI** that interact with the API and add interactivity to the game.

It was done using Go v1.22.4. `echo` was choosen as the for the web framework for the API. While `cobra` was used for the CLI

There is a postman collection for convenience, with the endpoints already configured. It can be found in the `GoQuiz.postman_collection.json` file

**There are some assumptions to the game**

- There can be only one ongoing quiz per user.
- Each user needs a username and password before they can start a quiz
- Leaderboards are divided by the number of questions. Meaning that if two different users decided to play a game with 4 questions and 8 questions respectively, they will be ranked in different leaderboards. This was made for fairness

**Implementation details**

- There is no database, all data is stored in memory. There is a mock in-memory database under `server/db/mockDatabase.go`. Therefore, every time you restart the application the data is reloaded. There is some mock data that is loaded on intialization to help with testing the application.
- The repositories (which interact with the database) implement interfaces instead of being directly used, and are injected. This allows for more easily changing the implementation of the app repositories to use a real database in the future.
- The business logic is handled inside the handlers itself for simplicity instead of having a new service or commands package/layer, depending on the architecture choosen.
- User passwords are stored as plain text for simplicity

## REST API

You can play the quiz interacting with a simple REST API. A brief explanation of each endpoint is given below.

Authorization is enforced on some endpoints using a Basic Authorization middleware (username and password in the Authorization HEADER)

### Endpoitns

#### **GET** /user/:id

Endpoint to get an user, it get all the user details (password included) and does not use a DTO to facilitate debugging

```json
Reponse example:
{
    "Id": 1,
    "Username": "user1",
    "Password": "pass1"
}
```

#### **POST** /user

Create a new user in the application. It expects a json payload with a username and password. It returns if it succeeded or there is already a user using that username

```json
Payload example
{
    "username": "user1",
    "password": "pass1"
}
```

#### **GET** /game/leaderboard/:numQuestions

Gets the current quizzes leaderboards. Leaderboards are divided by the number of questions a quiz had.

```json
Response example
[
    {
        "username": "user0",
        "highestScore": 75
    },
    {
        "username": "user1",
        "highestScore": 50
    },
    {
        "username": "user2",
        "highestScore": 25
    }
]
```

#### **GET** /game

**Required Basic Auth**
Get the current ongoing game for the authenticated user. If no game is found it returns a 404 response

```json
Response example
{
    "id": "11721689360",
    "username": "user1",
    "questions": [
        {
            "id": 3,
            "description": "What is a common term for free spins or bonus rounds in slots?",
            "options": [
                "Betting Circles",
                "Win Loops",
                "Free Rolls",
                "Scatter Bonuses"
            ]
        },
        ...
    ],
    "createdDate": "2024-07-22T21:02:40.7517045+02:00"
}
```

#### **POST** /game/:numQuestions

**Required Basic Auth**
Creates a new quiz with _numQuestions_ questions. If there is an ongoing game, a new game will not be created and the endpoint will return an unsuccessful status.

```json
Response example
{
  "id": "",
  "username": "user1",
  "questions": [
    {
      "id": 3,
      "description": "What is a common term for free spins or bonus rounds in slots?",
      "options": ["Betting Circles", "Win Loops", "Free Rolls", "Scatter Bonuses"]
    },
    ...
  ],
  "createdDate": "0001-01-01T00:00:00Z"
}
```

#### **POST** /game/answers

**Required Basic Auth**
Receives the answers for the current ongoing game of the user. If there is no ongoing game the endpoint will return an unsuccessful status.

```json
Payload example
[
    {
        "questionId":  13,
        "selectedOption": 3
    },
    ...
]

Response example
{
    "id": "31721681532",
    "username": "user1",
    "questionsResults": [
        {
            "question": {
                "id": 13,
                "description": "What does the term 'all-in' mean in poker?",
                "options": [
                    "Folding all cards",
                    "Raising the minimum bet",
                    "Betting all of one's chips",
                    "Splitting the pot"
                ]
            },
            "correctionOption": 2,
            "selectedOption": 3,
            "isCorrect": false
        },
        ...
    ],
    "scorePercentage": 25,
    "createdDate": "2024-07-22T22:52:12.1890397+02:00",
    "completedDate": "2024-07-22T23:19:39.1147689+02:00",
    "percentileScore": 0,
    "rankingPosition": 3
}
```

### CLI

The CLI is build using cobra. The commands below can be used to interact with the application

#### user new [username]

This can be used to create a new use in the game. It will ask for a password to be set.

#### quiz [username]

This can be used to start a new game, or continue one that you had started before.

#### leaderboard [numQuestions]

This endpoint will show you the leaderboard for quizzes with the size of the numQuestion questions given.

### Future improvementes

Since this is the first iteraction of the GoQuiz, there is a huge margin for improvement. Some of the items below could be taken into consideration on future versions:

- Use a real database instead of storing data in memory
- Allowing users or admins to create new questions, since currently there is a prefixed number of question that can be used on the games
- Shuffle questions options, currently the options are in a preset order
- Implement more unit testing to check all the handlers
- Implement caching, and ideally a cached repository to improve performance on similar request.
- Depending on if this quiz will be hosted on a distributed manner with multiple instances running, it could be beneficial to implement a distributed cache using `Redis` for example, which would also allow for distributed locking to prevent concurrency issues.
- Depending on the increase complexity of future iterations, could be beneficial to layer more that application, either using a pattern like CQRS or even layered application, removing the business logic from the handlers.

There is much more that can be improved. As I get more experienced with the Go programming language itself I will be able to spot a great deal of things that I might have overlooked in this project.
