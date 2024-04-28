package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// (POST /user) Register user endpoint
func (s *Server) Register(ctx echo.Context) error {
	var request generated.RegisterJSONRequestBody
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	errors := s.ValidateCreateUser(request)
	if len(errors.Messages) != 0 {
		return ctx.JSON(http.StatusBadRequest, errors)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 8)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
	}

	user, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), request.PhoneNumber)
	if err != nil && err != sql.ErrNoRows {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
	} else if user.Id != 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "phone number is already registered"})
	}

	result, err := s.Repository.CreateUser(ctx.Request().Context(), repository.CreateUserInput{
		Name:     request.FullName,
		Phone:    request.PhoneNumber,
		Password: string(password),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
	}

	return ctx.JSON(http.StatusOK, generated.RegisterResponse{
		Id: result.Id,
	})
}


// (POST /login) User authentication endpoint
func (s *Server) Login(ctx echo.Context) error {
	var request generated.LoginJSONRequestBody
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	errors := s.ValidateLoginUser(request)
	if len(errors.Messages) != 0 {
		return ctx.JSON(http.StatusBadRequest, errors)
	}

	user, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), request.PhoneNumber)
	if err != nil && err != sql.ErrNoRows {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "incorrect password or phone number"})
	}

	token, err := s.GenerateJWT(user.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
	}

	return ctx.JSON(http.StatusOK, generated.LoginResponse{
		Token: token,
	})
}


// (GET /user) Get user endpoint
func (s *Server) GetUser(ctx echo.Context, params generated.GetUserParams) error {
	if params.Authorization == "" {
		return ctx.String(http.StatusForbidden, "unauthorized")
	}

	token, err := s.ValidateJWT(params.Authorization)
	if !token.Valid || err != nil {
		return ctx.String(http.StatusForbidden, "unauthorized")
	}

	idClaims, err := s.GetJWTClaims(token, "id")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
	}

	id, err := strconv.Atoi(idClaims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
	}

	user, err := s.Repository.GetUserById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
	}

	return ctx.JSON(http.StatusOK, generated.User{
		FullName:    user.Name,
		PhoneNumber: user.Phone,
	})
}

// (PUT /user) Update user endpoint
func (s *Server) UpdateUser(ctx echo.Context, params generated.UpdateUserParams) error {
	if params.Authorization == "" {
		return ctx.String(http.StatusForbidden, "unauthorized")
	}

	token, err := s.ValidateJWT(params.Authorization)
	if !token.Valid || err != nil {
		return ctx.String(http.StatusForbidden, "unauthorized")
	}

	idClaims, err := s.GetJWTClaims(token, "id")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
	}

	id, err := strconv.Atoi(idClaims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
	}

	user, err := s.Repository.GetUserById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
	}

	var request generated.UpdateUserJSONRequestBody
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	errors := s.ValidateUpdateUser(request)
	if len(errors.Messages) != 0 {
		return ctx.JSON(http.StatusBadRequest, errors)
	}
	
	if request.PhoneNumber != "" {
		existingUser, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), request.PhoneNumber)
		if err != nil && err != sql.ErrNoRows {
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "internal error"})
		} else if existingUser.Id != 0 && request.PhoneNumber != user.Phone {
			return ctx.JSON(http.StatusConflict, generated.ErrorResponse{Message: "phone number is already registered"})
		}
		user.Phone = request.PhoneNumber
	}

	if request.FullName != "" {
		user.Name = request.FullName
	}

	result, err := s.Repository.UpdateUserById(ctx.Request().Context(), repository.UpdateUserInput{
		Id:    user.Id,
		Name:  user.Name,
		Phone: user.Phone,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, generated.User{
		FullName:    result.Name,
		PhoneNumber: result.Phone,
	})
}
