package http

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"finanzas-api/internal/users/domain"
	"finanzas-api/internal/users/port/in"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type stubUserService struct {
	createFunc   func(ctx context.Context, cmd in.CreateUserCommand) (*domain.User, error)
	getByIDFunc  func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	updateFunc   func(ctx context.Context, cmd in.UpdateUserCommand) (*domain.User, error)
	deleteFunc   func(ctx context.Context, id uuid.UUID) error
	listFunc     func(ctx context.Context, q in.ListUsersQuery) ([]*domain.User, error)
	getByEmailFn func(ctx context.Context, email string) (*domain.User, error)
}

func (s *stubUserService) Create(ctx context.Context, cmd in.CreateUserCommand) (*domain.User, error) {
	return s.createFunc(ctx, cmd)
}
func (s *stubUserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.getByIDFunc(ctx, id)
}
func (s *stubUserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.getByEmailFn(ctx, email)
}
func (s *stubUserService) Update(ctx context.Context, cmd in.UpdateUserCommand) (*domain.User, error) {
	return s.updateFunc(ctx, cmd)
}
func (s *stubUserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.deleteFunc(ctx, id)
}
func (s *stubUserService) List(ctx context.Context, q in.ListUsersQuery) ([]*domain.User, error) {
	return s.listFunc(ctx, q)
}

func newTestContext(method, target, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var reader *strings.Reader
	if body != "" {
		reader = strings.NewReader(body)
	} else {
		reader = strings.NewReader("")
	}
	c.Request = httptest.NewRequest(method, target, reader)
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

var fixedUser = &domain.User{
	ID:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
	Email:     "a@b.com",
	FirstName: "Ana",
	LastName:  "Pérez",
	Role:      "user",
	IsActive:  true,
	CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	UpdatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
}

func TestCreateUser_Success_NoPasswordInResponse(t *testing.T) {
	svc := &stubUserService{
		createFunc: func(ctx context.Context, cmd in.CreateUserCommand) (*domain.User, error) {
			return fixedUser, nil
		},
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("POST", "/users", `{"email":"a@b.com","first_name":"Ana","last_name":"Pérez","password":"secret1"}`)

	h.CreateUser(c)

	if w.Code != 201 {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	if strings.Contains(w.Body.String(), "secret1") || strings.Contains(w.Body.String(), "password") {
		t.Fatalf("response leaked the password: %s", w.Body.String())
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	svc := &stubUserService{}
	h := NewUserHandler(svc)
	c, w := newTestContext("POST", "/users", `not-json`)

	h.CreateUser(c)

	if w.Code != 400 {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateUser_ValidationError_PopulatesFields(t *testing.T) {
	svc := &stubUserService{}
	h := NewUserHandler(svc)
	c, w := newTestContext("POST", "/users", `{"email":"not-an-email","first_name":"","last_name":"","password":"secret1"}`)

	h.CreateUser(c)

	if w.Code != 422 {
		t.Fatalf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"fields"`) {
		t.Fatalf("expected fields in body, got %s", w.Body.String())
	}
}

func TestCreateUser_EmailTaken_Returns409(t *testing.T) {
	svc := &stubUserService{
		createFunc: func(ctx context.Context, cmd in.CreateUserCommand) (*domain.User, error) {
			return nil, domain.ErrEmailTaken
		},
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("POST", "/users", `{"email":"a@b.com","first_name":"Ana","last_name":"Pérez","password":"secret1"}`)

	h.CreateUser(c)

	if w.Code != 409 {
		t.Fatalf("expected 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateUser_GenericError_Returns500WithoutLeakingMessage(t *testing.T) {
	svc := &stubUserService{
		createFunc: func(ctx context.Context, cmd in.CreateUserCommand) (*domain.User, error) {
			return nil, errors.New("conexión perdida con postgres en 10.0.0.9")
		},
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("POST", "/users", `{"email":"a@b.com","first_name":"Ana","last_name":"Pérez","password":"secret1"}`)

	h.CreateUser(c)

	if w.Code != 500 {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
	if strings.Contains(w.Body.String(), "10.0.0.9") {
		t.Fatalf("internal error message leaked: %s", w.Body.String())
	}
}

func TestGetUser_NotFound_Returns404(t *testing.T) {
	svc := &stubUserService{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		},
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("GET", "/users/11111111-1111-1111-1111-111111111111", "")
	c.Params = gin.Params{{Key: "id", Value: "11111111-1111-1111-1111-111111111111"}}

	h.GetUser(c)

	if w.Code != 404 {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetUser_MalformedID_Returns400(t *testing.T) {
	svc := &stubUserService{}
	h := NewUserHandler(svc)
	c, w := newTestContext("GET", "/users/not-a-uuid", "")
	c.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}

	h.GetUser(c)

	if w.Code != 400 {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListUsers_DefaultPagination(t *testing.T) {
	var gotQuery in.ListUsersQuery
	svc := &stubUserService{
		listFunc: func(ctx context.Context, q in.ListUsersQuery) ([]*domain.User, error) {
			gotQuery = q
			return []*domain.User{fixedUser}, nil
		},
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("GET", "/users", "")

	h.ListUsers(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if gotQuery.Limit != 10 || gotQuery.Offset != 0 {
		t.Fatalf("expected default limit=10 offset=0, got %+v", gotQuery)
	}
}

func TestListUsers_NonNumericParams_DefaultToZero(t *testing.T) {
	var gotQuery in.ListUsersQuery
	svc := &stubUserService{
		listFunc: func(ctx context.Context, q in.ListUsersQuery) ([]*domain.User, error) {
			gotQuery = q
			return nil, nil
		},
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("GET", "/users?limit=abc&offset=xyz", "")

	h.ListUsers(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if gotQuery.Limit != 0 || gotQuery.Offset != 0 {
		t.Fatalf("expected limit/offset 0 on non-numeric input, got %+v", gotQuery)
	}
}
