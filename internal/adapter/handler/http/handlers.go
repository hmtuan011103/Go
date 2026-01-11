package http

import "github.com/gostructure/app/internal/adapter/handler/http/user"

type Handlers struct {
	User *user.UserHandler
	Auth *user.AuthHandler
}
