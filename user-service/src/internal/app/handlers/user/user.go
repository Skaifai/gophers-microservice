package user

import (
	"context"
	"errors"

	"github.com/Skaifai/gophers-microservice/user-service/internal/app/models/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/helpers"
	"github.com/Skaifai/gophers-microservice/user-service/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service interface {
	GetAll(ctx context.Context, offset int64, limit int64) (_ []user.User, err error)
	GetByUserID(ctx context.Context, id string) (_ *user.User, err error)
	DeleteUserByID(ctx context.Context, id string) (err error)
	UpdateUser(ctx context.Context, u *user.User) (_ *user.User, err error)
}

type handler struct {
	proto.UnimplementedUserServiceServer
	usvc service
}

func New(usvc service) *handler {
	return &handler{
		usvc: usvc,
	}
}

func (h *handler) GetAllUsers(ctx context.Context, req *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	us, err := h.usvc.GetAll(ctx, req.GetOffset(), req.GetLimit())
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
	u, err := h.usvc.GetByUserID(ctx, helpers.Itoa64(req.GetId()))
	if err != nil {
		return &proto.GetUserResponse{Status: 404}, err
	}
	return &proto.GetUserResponse{User: protoFromModel(u), Status: 200}, nil
}

func (h *handler) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	u, err := h.usvc.UpdateUser(ctx, protoToModel(req.GetUser()))

	if err != nil {
		return &proto.UpdateUserResponse{Status: 500}, err
	}

	return &proto.UpdateUserResponse{User: protoFromModel(u), Status: 200}, nil
}

func (h *handler) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	err := h.usvc.DeleteUserByID(ctx, helpers.Itoa64(req.GetId()))
	if err != nil {
		return &proto.DeleteUserResponse{Status: 500}, nil
	}

	return &proto.DeleteUserResponse{Status: 200}, nil
}

func (h *handler) Login(ctx context.Context, _ *proto.LoginRequest) (*proto.LoginResponse, error) {
	return nil, nil
}

func (h *handler) Registration(ctx context.Context, _ *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
	return nil, nil
}

func (h *handler) Logout(ctx context.Context, _ *proto.LogoutRequest) (*proto.LogoutResponse, error) {

	return nil, nil
}

func (h *handler) Refresh(ctx context.Context, _ *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	return nil, nil
}

func (h *handler) Activate(ctx context.Context, req *proto.ActivateRequest) (*proto.ActivateResponse, error) {
	return nil, nil
}

func (h *handler) IsLogged(ctx context.Context, _ *proto.IsLoggedRequest) (*proto.IsLoggedResponse, error) {
	return nil, nil
}

func (h *handler) GetUserByToken(ctx context.Context, _ *proto.GetUserByTokenRequest) (*proto.GetUserByTokenResponse, error) {
	return nil, nil
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
