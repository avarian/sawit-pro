package handler

import "github.com/SawitProRecruitment/UserService/repository"

type Server struct {
	Repository repository.RepositoryInterface
	Secret string
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Secret string
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		Secret: opts.Secret,
	}
}
