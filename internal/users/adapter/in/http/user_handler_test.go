package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"finanzas-api/internal/httpx"
	"finanzas-api/internal/users/application"
	"finanzas-api/internal/users/domain"
	"finanzas-api/internal/users/port/in"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type stubUserService struct {
	createFunc   func(ctx context.Context, cmd in.CreateUserCommand) error
	getByIDFunc  func(ctx context.Context, actor in.Actor, id uuid.UUID) (*domain.User, error)
	updateFunc   func(ctx context.Context, actor in.Actor, cmd in.UpdateUserCommand) (*domain.User, error)
	deleteFunc   func(ctx context.Context, id uuid.UUID) error
	listFunc     func(ctx context.Context, q in.ListUsersQuery) ([]*domain.User, error)
	getByEmailFn func(ctx context.Context, email string) (*domain.User, error)
}

func (s *stubUserService) Create(ctx context.Context, cmd in.CreateUserCommand) error {
	return s.createFunc(ctx, cmd)
}
func (s *stubUserService) GetByID(ctx context.Context, actor in.Actor, id uuid.UUID) (*domain.User, error) {
	return s.getByIDFunc(ctx, actor, id)
}
func (s *stubUserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.getByEmailFn(ctx, email)
}
func (s *stubUserService) Update(ctx context.Context, actor in.Actor, cmd in.UpdateUserCommand) (*domain.User, error) {
	return s.updateFunc(ctx, actor, cmd)
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

func TestCreateUser_Success_Returns202WithoutUserData(t *testing.T) {
	svc := &stubUserService{
		createFunc: func(ctx context.Context, cmd in.CreateUserCommand) error { return nil },
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("POST", "/users", `{"email":"a@b.com","first_name":"Ana","last_name":"Pérez","password":"rio-cometa-azul"}`)

	h.CreateUser(c)

	if w.Code != 202 {
		t.Fatalf("expected 202, got %d: %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, leak := range []string{"rio-cometa-azul", "password", "hash", `"id"`, `"user"`, "a@b.com"} {
		if strings.Contains(body, leak) {
			t.Fatalf("response leaked %q: %s", leak, body)
		}
	}
	if !strings.Contains(body, registrationAcceptedMessage) {
		t.Fatalf("expected the generic message, got %s", body)
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

// TestCreateUser_ExistingEmail_SameResponseAsNewEmail usa el servicio real
// de application con un repositorio en memoria: la única diferencia entre
// los dos registros es si el correo ya existe.
func TestCreateUser_ExistingEmail_SameResponseAsNewEmail(t *testing.T) {
	register := func(existing bool) (int, string, *memUserRepository) {
		repo := &memUserRepository{exists: existing}
		h := NewUserHandler(application.NewUserService(repo, fakeHasher{}))
		c, w := newTestContext("POST", "/users", `{"email":"a@b.com","first_name":"Ana","last_name":"Pérez","password":"rio-cometa-azul"}`)
		h.CreateUser(c)
		return w.Code, w.Body.String(), repo
	}

	newCode, newBody, newRepo := register(false)
	takenCode, takenBody, takenRepo := register(true)

	if newCode != 202 || takenCode != 202 {
		t.Fatalf("expected 202 for both, got new=%d taken=%d", newCode, takenCode)
	}
	if newBody != takenBody {
		t.Fatalf("bodies differ:\nnew:   %s\ntaken: %s", newBody, takenBody)
	}
	if newRepo.saved != 1 || takenRepo.saved != 0 {
		t.Fatalf("expected one save for the new email and none for the taken one, got new=%d taken=%d", newRepo.saved, takenRepo.saved)
	}
}

func TestCreateUser_GenericError_Returns500WithoutLeakingMessage(t *testing.T) {
	svc := &stubUserService{
		createFunc: func(ctx context.Context, cmd in.CreateUserCommand) error {
			return errors.New("conexión perdida con postgres en 10.0.0.9")
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
		getByIDFunc: func(ctx context.Context, actor in.Actor, id uuid.UUID) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		},
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("GET", "/users/11111111-1111-1111-1111-111111111111", "")
	c.Params = gin.Params{{Key: "id", Value: "11111111-1111-1111-1111-111111111111"}}
	setActor(c, fixedUser.ID.String(), "user")

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

// setActor simula lo que deja el AuthGuard en el contexto.
func setActor(c *gin.Context, id, role string) {
	c.Set(httpx.CtxUserID, id)
	c.Set(httpx.CtxUserRole, role)
}

// --- Autorización (GHSA-j36w-x5rp-5f3w) ---

func TestGetUser_WithoutActor_Returns401(t *testing.T) {
	svc := &stubUserService{
		getByIDFunc: func(ctx context.Context, actor in.Actor, id uuid.UUID) (*domain.User, error) {
			t.Fatal("service must not be called without an authenticated actor")
			return nil, nil
		},
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("GET", "/users/"+fixedUser.ID.String(), "")
	c.Params = gin.Params{{Key: "id", Value: fixedUser.ID.String()}}

	h.GetUser(c)

	if w.Code != 401 {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetUser_PassesActorFromContext(t *testing.T) {
	var got in.Actor
	svc := &stubUserService{
		getByIDFunc: func(ctx context.Context, actor in.Actor, id uuid.UUID) (*domain.User, error) {
			got = actor
			return fixedUser, nil
		},
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("GET", "/users/"+fixedUser.ID.String(), "")
	c.Params = gin.Params{{Key: "id", Value: fixedUser.ID.String()}}
	setActor(c, fixedUser.ID.String(), "user")

	h.GetUser(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if got.ID != fixedUser.ID || got.Role != "user" {
		t.Fatalf("actor not propagated: %+v", got)
	}
}

func TestUpdateUser_Forbidden_Returns403(t *testing.T) {
	svc := &stubUserService{
		updateFunc: func(ctx context.Context, actor in.Actor, cmd in.UpdateUserCommand) (*domain.User, error) {
			return nil, domain.ErrForbidden
		},
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("PUT", "/users/"+fixedUser.ID.String(), `{"role":"admin"}`)
	c.Params = gin.Params{{Key: "id", Value: fixedUser.ID.String()}}
	setActor(c, fixedUser.ID.String(), "user")

	h.UpdateUser(c)

	if w.Code != 403 {
		t.Fatalf("expected 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateUser_WithoutActor_Returns401(t *testing.T) {
	svc := &stubUserService{
		updateFunc: func(ctx context.Context, actor in.Actor, cmd in.UpdateUserCommand) (*domain.User, error) {
			t.Fatal("service must not be called without an authenticated actor")
			return nil, nil
		},
	}
	h := NewUserHandler(svc)
	c, w := newTestContext("PUT", "/users/"+fixedUser.ID.String(), `{"first_name":"X"}`)
	c.Params = gin.Params{{Key: "id", Value: fixedUser.ID.String()}}
	setActor(c, "not-a-uuid", "user")

	h.UpdateUser(c)

	if w.Code != 401 {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateUser_RoleInBodyIsIgnored(t *testing.T) {
	var got in.CreateUserCommand
	svc := &stubUserService{
		createFunc: func(ctx context.Context, cmd in.CreateUserCommand) error {
			got = cmd
			return nil
		},
	}
	h := NewUserHandler(svc)
	body := `{"email":"x@y.com","first_name":"X","last_name":"Y","password":"secret-123","role":"admin"}`
	c, w := newTestContext("POST", "/users", body)

	h.CreateUser(c)

	if w.Code != 202 {
		t.Fatalf("expected 202, got %d: %s", w.Code, w.Body.String())
	}
	if got.Email != "x@y.com" {
		t.Fatalf("command not built from body: %+v", got)
	}
}

// --- SPEC 01: el dominio es la única fuente de las reglas ---

func TestCreateUser_DomainRulesReachTheResponse(t *testing.T) {
	cases := []struct {
		name, field, body, want string
	}{
		{
			name:  "contraseña de 6 caracteres",
			field: "password",
			body:  `{"email":"a@b.com","first_name":"Ana","last_name":"Pérez","password":"secret"}`,
			want:  domain.ValidatePassword("secret", "a@b.com", "Ana", "Pérez"),
		},
		{
			name:  "contraseña que contiene el nombre",
			field: "password",
			body:  `{"email":"a@b.com","first_name":"Ana","last_name":"Pérez","password":"cometa-ana-azul"}`,
			want:  domain.ValidatePassword("cometa-ana-azul", "a@b.com", "Ana", "Pérez"),
		},
		{
			name:  "correo con nombre visible",
			field: "email",
			body:  `{"email":"Ana <a@b.co>","first_name":"Ana","last_name":"Pérez","password":"rio-cometa-azul"}`,
			want:  "formato de email inválido",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &memUserRepository{}
			h := NewUserHandler(application.NewUserService(repo, fakeHasher{}))
			c, w := newTestContext("POST", "/users", tc.body)

			h.CreateUser(c)

			if w.Code != 422 {
				t.Fatalf("expected 422, got %d: %s", w.Code, w.Body.String())
			}
			var resp struct {
				Error struct {
					Fields map[string]string `json:"fields"`
				} `json:"error"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("invalid JSON body: %v", err)
			}
			if tc.want == "" || resp.Error.Fields[tc.field] != tc.want {
				t.Fatalf("expected fields.%s = %q, got %v", tc.field, tc.want, resp.Error.Fields)
			}
			if repo.saved != 0 {
				t.Fatal("expected nothing saved on invalid data")
			}
		})
	}
}

// --- Dobles para probar el handler con el servicio real (SPEC 01) ---

// memUserRepository implementa out.UserRepository en memoria: solo lo que
// usa el registro. exists simula un correo ya registrado.
type memUserRepository struct {
	exists bool
	saved  int
}

func (r *memUserRepository) Save(ctx context.Context, u *domain.User) error {
	r.saved++
	return nil
}
func (r *memUserRepository) Update(ctx context.Context, u *domain.User) error { return nil }
func (r *memUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, domain.ErrUserNotFound
}
func (r *memUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, domain.ErrUserNotFound
}
func (r *memUserRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (r *memUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (r *memUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.exists, nil
}

type fakeHasher struct{}

func (fakeHasher) Hash(plain string) (string, error) { return "hashed", nil }
