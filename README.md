# A Golang Project for GWI's Challenge

This project is a solution to a challenge where I needed to build an API that lets users manage their favourite assets.
Users can add, update, view, and delete favourites like charts, insights, or audience profiles.

## How to run tests

I’ve included a Makefile in the project to simplify common development tasks. To run the test suite, simply execute:

```
make test
```

This command will run all unit and e2e tests for the project.
Make sure you have any necessary dependencies installed before running it.

## How to Test the Application

### 1. Build and Run the Docker Container

Build the Docker image with:

```bash
docker build -t gwi_api .
```

Run the container and expose port 3008:

```bash
docker run -p 3008:3008 gwi_api
```

### 2. Seeded Development Environment

By default, the container runs with environment variables set for a development environment that automatically seeds the database with test data. This allows you to start testing immediately without additional setup.

### 3. Seeded Data

Here is the data that is pre-populated in the database on startup:

```json
{
  "users": [
    {
      "id": "a3973a1c-a77b-4a04-a296-ddec19034419",
      "email": "test@test.com",
      "password": "pass" // hashed in the actual database
    }
  ],
  "charts": [
    {
      "id": "11111111-1111-1111-1111-111111111111",
      "title": "test chart",
      "x_axis_title": "commit number",
      "y_axis_title": "lines of code",
      "data": [
        { "x": 1, "y": 100 },
        { "x": 2, "y": 300 },
        { "x": 3, "y": 500 }
      ]
    }
  ],
  "insights": [
    {
      "id": "22222222-2222-2222-2222-222222222222",
      "text": "40% of millennials spend more than 3 hours on social media daily"
    },
    {
      "id": "22222222-2222-2222-2222-222222222223",
      "text": "100% of zoomers spend more than 8 hours on watching memes"
    }
  ],
  "audiences": [
    {
      "id": "33333333-3333-3333-3333-333333333333",
      "gender": "Male",
      "birth_country": "United Kingdom",
      "age_group": "25-34",
      "social_media_hours": 3.5,
      "purchases_last_month": 7
    }
  ],
  "favourites": [
    {
      "id": "44444444-4444-4444-4444-444444444444",
      "user_id": "a3973a1c-a77b-4a04-a296-ddec19034419",
      "asset_id": "11111111-1111-1111-1111-111111111111",
      "asset_type": "chart",
      "description": "Main performance chart"
    },
    {
      "id": "55555555-5555-5555-5555-555555555555",
      "user_id": "a3973a1c-a77b-4a04-a296-ddec19034419",
      "asset_id": "22222222-2222-2222-2222-222222222222",
      "asset_type": "insight",
      "description": "Great for Q2 presentation"
    },
    {
      "id": "66666666-6666-6666-6666-666666666666",
      "user_id": "a3973a1c-a77b-4a04-a296-ddec19034419",
      "asset_id": "33333333-3333-3333-3333-333333333333",
      "asset_type": "audience",
      "description": "Target audience for campaign"
    }
  ]
}
```

### 4. Authenticate

Use the login endpoint with the pre-seeded user's email (`test@test.com`) and password (`pass`) to obtain a JWT token.

### 5. Test Favourite Endpoints

With the obtained token, you can now call the protected favourite endpoints. Each endpoint includes detailed documentation on usage.

## Some of my thoughts while implementing this

29/05/25

- Which file structure to follow?<br/>
  I've decided to follow the golang-standards/project-layout convention, as I find it clean, widely recognized, and scalable.
  It provides a good foundation for organizing code in a way that's maintainable and familiar to other Go developers.
  For my internals folder I will follow a Domain Driven Design having each file under its underlying domain.

- Which Framework to use?<br/>
  I'm using Chi, which is a router that's close to Go's standard net/http package.
  I like its idiomatic and minimalistic design, which keeps things simple.

- Architecture and Dependency injection?<br/>
  I chose a layered architecture for this project, separating the code into a presentation layer, service layer, and repository layer.
  I’m trying to approach the implementation with a procedural mindset rather than following object-oriented design patterns.
  Separating code into services and repositories feels like object-oriented design, I think it’s needed to keep things organized and clear.
  It’s hard to avoid OOP patterns completely in this case.

- Stateless or stateful authentication?<br/>
  I chose JWT-based stateless authentication to keep the service simple and decoupled from any central session store.
  I'm skipping refresh tokens altogether and instead issuing short-lived JWTs.
  This enables fast request handling, as the server can validate tokens locally without needing to query an external store for each authorization check.

31/05/25

- Which Database to use?<br/>
  Since the challenge does not require using a real database, I will implement an in-memory store for simplicity.
  Thanks to the current layered architecture, the code remains loosely coupled, making it straightforward to switch to a real database in the future without significant refactoring.
  In a real life scenario I would probably select a nosql database since assets have varying structures and it does not require complex relational queries.

- Use a separate domain for each asset type instead of a single asset domain? <br/>
  Even though this adds a bit of repetition, it keeps things cleaner and easier to work with.
  Each asset has its own space, so if we ever want to add something specific, like chart versioning or special filters for audiences,
  we can do it without messing with other parts of the code.

03/06/25

- The decision of building everything from scratch.<br/>
  I decided to build most stuff from scratch, relying mostly on Go’s built-in packages.
  This was a choice to get more hands-on work with Golang.
  Had a lot of fun :)
