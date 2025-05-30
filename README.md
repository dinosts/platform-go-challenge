Thought Process:

29/05/25

- Which file structure to follow?
  I've decided to follow the golang-standards/project-layout convention, as I find it clean, widely recognized, and scalable.
  It provides a good foundation for organizing code in a way that's maintainable and familiar to other Go developers.

- Which Framework to use?
  I'm using Chi, which is a router that's close to Go's standard net/http package.
  I like its idiomatic and minimalistic design, which keeps things simple.

- Architecture and Dependency injection?
  I chose a layered architecture for this project, separating the code into a presentation layer, service layer, and repository layer.
  I’m trying to approach the implementation with a procedural mindset rather than following object-oriented design patterns.

- Stateless or stateful authentication?
  I chose JWT-based stateless authentication to keep the service simple and decoupled from any central session store.
  I'm skipping refresh tokens altogether and instead issuing short-lived JWTs.
  This enables fast request handling, as the server can validate tokens locally without needing to query an external store for each authorization check.
