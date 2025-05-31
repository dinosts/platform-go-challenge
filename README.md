Thought Process:

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
