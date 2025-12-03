package endpoint

//go:generate go run github.com/youbuwei/doeot-go/cmd/bizgen -module user

import (
	"errors"

	"github.com/youbuwei/doeot-go/internal/user/app"
	"github.com/youbuwei/doeot-go/internal/user/domain"
	"github.com/youbuwei/doeot-go/pkg/biz"
	"github.com/youbuwei/doeot-go/pkg/errs"
	"github.com/youbuwei/doeot-go/pkg/validate"
)

// GetUserReq describes the input of GetUser endpoint.
// `validate:"gt=0"` is used by the validator package for pre-flight checks.
type GetUserReq struct {
	ID int64 `path:"id" json:"id" validate:"gt=0"`
}

// GetUserResp is the output DTO of GetUser.
type GetUserResp struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Role  string `json:"role"`
	Phone string `json:"phone"`
}

type GetUserListReq struct{}

type GetUserListResp struct{}

// CreateUserReq describes the input of CreateUser endpoint.
type CreateUserReq struct {
	Name  string `json:"name" validate:"required,min=3"`
	Age   int    `json:"age" validate:"gte=0,lte=120"`
	Role  string `json:"role" validate:"omitempty,oneof=normal admin"`
	Phone string `json:"phone" validate:"required,mobile"`
}

// Validate implements validate.Custom for CreateUserReq.
// Here we enforce some cross-field business rules.
func (r *CreateUserReq) Validate() error {
	if r.Role == "admin" && r.Age < 18 {
		return errs.BadRequest("admin 用户必须年满 18 岁")
	}
	return nil
}

// CreateUserResp is the output DTO of CreateUser.
type CreateUserResp struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Role  string `json:"role"`
	Phone string `json:"phone"`
}

// UserEndpoint groups all user-related endpoints.
// Its dependencies are injected by wiring code in the user module.
type UserEndpoint struct {
	Svc *app.UserService
}

// GetUser is a demo handler which will be wired to both HTTP and RPC
// via annotations + code generation.
//
// @Route  GET /user/:id
// @RPC    User.Get
// @Auth   login
// @Desc   获取用户详情
// @Tags   user
func (e *UserEndpoint) GetUser(ctx biz.Context, req *GetUserReq) (*GetUserResp, error) {
	// By the time we arrive here, basic request validation (gt=0) is already
	// performed by the generated wrapper using validate.Struct.
	u, err := e.Svc.GetUser(ctx.RequestContext(), req.ID)
	if errors.Is(err, domain.ErrUserNotFound) {
		return nil, errs.NotFound("user not found")
	}
	if err != nil {
		return nil, errs.Internal("failed to get user").WithCause(err)
	}

	return &GetUserResp{
		ID:    u.ID,
		Name:  u.Name,
		Age:   u.Age,
		Role:  u.Role,
		Phone: u.Phone,
	}, nil
}

// GetUserList get user list
// @Route  GET /users
// @RPC    Users.Get
// @Auth   login
// @Desc   获取用户列表
// @Tags   users
func (e *UserEndpoint) GetUserList(ctx biz.Context, req *GetUserListReq) ([]*GetUserResp, error) {
	l, err := e.Svc.GetUserList(ctx.RequestContext())
	if errors.Is(err, domain.ErrUserNotFound) {
		return nil, errs.NotFound("user not found")
	}
	if err != nil {
		return nil, errs.Internal("failed to get user list").WithCause(err)
	}

	res := make([]*GetUserResp, 0, len(l))
	for i := range l {
		res = append(res, &GetUserResp{
			ID:    l[i].ID,
			Name:  l[i].Name,
			Age:   l[i].Age,
			Role:  l[i].Role,
			Phone: l[i].Phone,
		})
	}
	return res, nil
}

// CreateUser creates a new user based on validated request.
//
// @Route  POST /user
// @RPC    User.Create
// @Auth   login
// @Desc   创建用户
// @Tags   user
func (e *UserEndpoint) CreateUser(ctx biz.Context, req *CreateUserReq) (*CreateUserResp, error) {
	// At this point, tag-based and custom Validate() have already run.
	u := &domain.User{
		Name:  req.Name,
		Age:   req.Age,
		Role:  req.Role,
		Phone: req.Phone,
	}

	created, err := e.Svc.CreateUser(ctx.RequestContext(), u)
	if err != nil {
		return nil, errs.Internal("failed to create user").WithCause(err)
	}

	return &CreateUserResp{
		ID:    created.ID,
		Name:  created.Name,
		Age:   created.Age,
		Role:  created.Role,
		Phone: created.Phone,
	}, nil
}

// Ensure CreateUserReq satisfies validate.Custom at compile time.
var _ validate.Custom = (*CreateUserReq)(nil)
