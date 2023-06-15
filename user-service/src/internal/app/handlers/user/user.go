package user

import (
	"context"
	"errors"

	"github.com/Skaifai/gophers-microservice/user-service/internal/app/models/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/helpers"
	"github.com/Skaifai/gophers-microservice/user-service/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrServerError = errors.New("internal server error")
)

type service interface {
	GetAll(ctx context.Context, offset int64, limit int64) (_ []user.User, err error)
	GetByUserID(ctx context.Context, id string) (_ *user.User, err error)
	DeleteUserByID(ctx context.Context, id string) (err error)
	UpdateUser(ctx context.Context, u *user.User) (_ *user.User, err error)
	Registrate(ctx context.Context, u *user.User) (*user.User, error)
	Activate(ctx context.Context, activation_link string) (activated bool, err error)
	Login(ctx context.Context, key, userAgent, password string) (accessToken string, refreshToken string, err error)
	IsLogged(ctx context.Context, tokenString string) (bool, error)
	GetByToken(ctx context.Context, accessToken string) (_ *user.User, err error)
	Refresh(ctx context.Context, refreshToken string, userAgent string) (accessToken string, err error)
	Logout(ctx context.Context, refreshToken string) error
}

type handler struct {
	proto.UnimplementedUserServiceServer
	UserService service
}

func New(usvc service) *handler {
	return &handler{
		UserService: usvc,
	}
}

func (h *handler) GetAllUsers(ctx context.Context, req *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	us, err := h.UserService.GetAll(ctx, req.GetOffset(), req.GetLimit())
	if err != nil {
		return &proto.GetAllUsersResponse{Status: 404}, errors.New("There are some error occured")
	}

	users := []*proto.User{}

	for _, u := range us {
		user := protoFromModel(&u)
		users = append(users, user)
	}

	return &proto.GetAllUsersResponse{Users: users, Status: 200}, nil
}

func (h *handler) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	u, err := h.UserService.GetByUserID(ctx, helpers.Itoa64(req.GetId()))
	if err != nil {
		return &proto.GetUserResponse{Status: 404}, err
	}
	return &proto.GetUserResponse{User: protoFromModel(u), Status: 200}, nil
}

func (h *handler) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	u, err := h.UserService.UpdateUser(ctx, protoToModel(req.GetUser()))

	if err != nil {
		return &proto.UpdateUserResponse{Status: 500}, err
	}

	return &proto.UpdateUserResponse{User: protoFromModel(u), Status: 200}, nil
}

func (h *handler) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	err := h.UserService.DeleteUserByID(ctx, helpers.Itoa64(req.GetId()))
	if err != nil {
		return &proto.DeleteUserResponse{Status: 500}, nil
	}

	return &proto.DeleteUserResponse{Status: 200}, nil
}

func (h *handler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	accessToken, refreshToken, err := h.UserService.Login(ctx, req.GetKey(), req.GetUserAgent(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &proto.LoginResponse{Tokens: &proto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}}, nil
}

func (h *handler) Registration(ctx context.Context, req *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
	user := &user.User{
		Username:    req.GetUsername(),
		Email:       req.GetEmail(),
		Password:    req.GetPassword(),
		FirstName:   req.GetFirstName(),
		LastName:    req.GetLastName(),
		PhoneNumber: req.GetPhoneNumber(),
		DOB:         req.GetDOB().AsTime(),
		Address:     req.GetAddress(),
		AboutMe:     req.GetAboutMe(),
		ProfPicUrl:  req.GetProfPicURL(),
	}

	u, err := h.UserService.Registrate(ctx, user)
	if err != nil {
		return &proto.RegistrationResponse{Status: 400}, err
	}

	res := protoFromModel(u)
	return &proto.RegistrationResponse{Userdata: res, Status: 200}, nil
}

func (h *handler) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	err := h.UserService.Logout(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &proto.LogoutResponse{Status: 200}, nil
}

func (h *handler) Refresh(ctx context.Context, req *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	accessToken, err := h.UserService.Refresh(ctx, req.GetRefreshToken(), req.GetUserAgent())
	if err != nil {
		return nil, err
	}

	return &proto.RefreshResponse{AccessToken: &proto.Tokens{AccessToken: accessToken, RefreshToken: req.GetRefreshToken()}}, nil
}

func (h *handler) Activate(ctx context.Context, req *proto.ActivateRequest) (*proto.ActivateResponse, error) {
	activation_string := req.GetActivationString()

	activated, err := h.UserService.Activate(ctx, activation_string)
	if err != nil {
		return &proto.ActivateResponse{Activated: false}, err
	}

	return &proto.ActivateResponse{Activated: activated}, nil
}

func (h *handler) IsLogged(ctx context.Context, req *proto.IsLoggedRequest) (*proto.IsLoggedResponse, error) {
	is_logged, err := h.UserService.IsLogged(ctx, req.GetRefreshToken())
	if err != nil {
		return &proto.IsLoggedResponse{IsLogged: false}, err
	}

	return &proto.IsLoggedResponse{IsLogged: is_logged}, nil
}

func (h *handler) GetUserByToken(ctx context.Context, req *proto.GetUserByTokenRequest) (*proto.GetUserByTokenResponse, error) {
	user, err := h.UserService.GetByToken(ctx, req.AccessToken)
	if err != nil {
		return &proto.GetUserByTokenResponse{Status: 404}, err
	}

	return &proto.GetUserByTokenResponse{
		User: protoFromModel(user),
	}, nil
}

func (h *handler) mustEmbedUnimplementedUserServiceServer() {}

func protoToModel(u *proto.User) *user.User {
	return &user.User{
		ID:               u.Id,
		Role:             u.UserRole,
		Username:         u.Username,
		Email:            u.Email,
		Password:         u.Password,
		RegistrationDate: u.RegistrationDate.AsTime(),
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		PhoneNumber:      u.PhoneNumber,
		DOB:              u.DOB.AsTime(),
		Address:          u.Address,
		AboutMe:          u.AboutMe,
		ProfPicUrl:       u.ProfPicURL,
		Activated:        u.Activated,
		Version:          u.Version,
	}
}

func protoFromModel(u *user.User) *proto.User {
	registrationDate := timestamppb.New(u.RegistrationDate)
	DOB := timestamppb.New(u.DOB)

	return &proto.User{
		Id:               u.ID,
		UserRole:         u.Role,
		Username:         u.Username,
		Email:            u.Email,
		Password:         u.Password,
		RegistrationDate: registrationDate,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		PhoneNumber:      u.PhoneNumber,
		DOB:              DOB,
		Address:          u.Address,
		AboutMe:          u.AboutMe,
		ProfPicURL:       u.ProfPicUrl,
		Activated:        u.Activated,
		Version:          u.Version,
	}
}
