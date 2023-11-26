# AlchemistTemplate

## How to run

1. Run `cat .env-example > .env.dev`
2. Update your config at `.env.dev`
3. Run `make setup`
4. Run `make dev`

## Features

- [x] Read application configurations from .env file
- [x] Dependencies injection container
- [x] Logger implementation
- [x] Integrate with Open telemetry
- [x] Integrate with Auth0
- [ ] Users management
- [ ] Unit testing
- [ ] Integrate CI/CD

## Key technologies

- Go programing language
- PostgreSQL database
- Docker
- Auth0
- Hexagonal architecture

## References

- [Learn how to start using RS256 for signing and verifying your JWTs.](https://auth0.com/blog/navigating-rs256-and-jwks/)
- [Hexagonal Architecture with Go - A Thorough Exploration of Backend Engineering and Distributed System](https://github.com/LordMoMA/Hexagonal-Architecture)
- [Building RESTful API with Hexagonal Architecture in Go](https://dev.to/bagashiz/building-restful-api-with-hexagonal-architecture-in-go-1mij)
- [REST guidelines suggest using a specific HTTP method on a particular type of call made to the server i.e. GET, POST, PUT or DELETE.](https://restfulapi.net/http-methods/)
- [go-jwt-middleware](https://github.com/auth0/go-jwt-middleware)
- [Add authorization to a Go application.](https://auth0.com/docs/quickstart/backend/golang/interactive)
