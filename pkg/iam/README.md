# IAM - Identity and Access Management package

## Key features

- [x] Download and cached JWKs
- [x] Parse JWT from string to [jwt.Token](https://pkg.go.dev/github.com/form3tech-oss/jwt-go@v3.2.2+incompatible#Token)
- [x] Get [UserProfile](types.go) from token
- [ ] RBAC implementation

## Folder Structure

```
.
├── README.md
├── context.go
├── errors.go
├── jwks_client.go
├── jwt.go
├── types.go
├── user_profile.go
└── validator.go
```

* `types.go` contains definitions of the `Validator` interface, `UserProfile`, `jwks`, and `jwk` struct.
* `validator.go` contains `New` function, `authZeroValidator` - that polling download JWKs and cached public key for validated JWT, and its methods.
* `jwks_client.go` contains definitions of the `jwksClient` - that have methods for get JWKs and parse to public key.
* `jwt.go` contains method for parse JWT from string to `jwt.Token`.
* `user_profile.go` contains method for get user profile from `jwt.Token`.
* `context.go` contains method for get/set user profile from context.
