package application

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"testing"

	"finanzas-api/internal/users/domain"
	"finanzas-api/internal/users/port/in"

	"github.com/google/uuid"
)

type stubUserRepository struct {
	saveFunc          func(ctx context.Context, u *domain.User) error
	updateFunc        func(ctx context.Context, u *domain.User) error
	findByIDFunc      func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	findByEmailFunc   func(ctx context.Context, email string) (*domain.User, error)
	deleteFunc        func(ctx context.Context, id uuid.UUID) error
	listFunc          func(ctx context.Context, limit, offset int) ([]*domain.User, error)
	existsByEmailFunc func(ctx context.Context, email string) (bool, error)

	saveCalled          bool
	existsByEmailCalled bool
}

func (s *stubUserRepository) Save(ctx context.Context, u *domain.User) error {
	s.saveCalled = true
	if s.saveFunc != nil {
		return s.saveFunc(ctx, u)
	}
	return nil
}

func (s *stubUserRepository) Update(ctx context.Context, u *domain.User) error {
	if s.updateFunc != nil {
		return s.updateFunc(ctx, u)
	}
	return nil
}

func (s *stubUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if s.findByIDFunc != nil {
		return s.findByIDFunc(ctx, id)
	}
	return nil, domain.ErrUserNotFound
}

func (s *stubUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if s.findByEmailFunc != nil {
		return s.findByEmailFunc(ctx, email)
	}
	return nil, domain.ErrUserNotFound
}

func (s *stubUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if s.deleteFunc != nil {
		return s.deleteFunc(ctx, id)
	}
	return nil
}

func (s *stubUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if s.listFunc != nil {
		return s.listFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (s *stubUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	s.existsByEmailCalled = true
	if s.existsByEmailFunc != nil {
		return s.existsByEmailFunc(ctx, email)
	}
	return false, nil
}

type stubPasswordHasher struct {
	hashFunc func(plain string) (string, error)
}

func (h *stubPasswordHasher) Hash(plain string) (string, error) {
	if h.hashFunc != nil {
		return h.hashFunc(plain)
	}
	return "hashed:" + plain, nil
}

func TestCreate_HashesPassword_NeverStoresPlainText(t *testing.T) {
	var user *domain.User
	repo := &stubUserRepository{saveFunc: func(ctx context.Context, u *domain.User) error {
		user = u
		return nil
	}}
	hasher := &stubPasswordHasher{}
	svc := NewUserService(repo, hasher)

	err := svc.Create(context.Background(), in.CreateUserCommand{
		Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Password: "plain-text-password",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if user.PasswordHash == "plain-text-password" {
		t.Fatal("Create stored the plain text password instead of the hash")
	}
	if user.PasswordHash != "hashed:plain-text-password" {
		t.Fatalf("expected hashed password to reach the repository, got %q", user.PasswordHash)
	}
	if !repo.saveCalled {
		t.Fatal("expected repo.Save to be called")
	}
}

func TestCreate_EmailAlreadyTaken_ReturnsNilWithoutSaving(t *testing.T) {
	repo := &stubUserRepository{
		existsByEmailFunc: func(ctx context.Context, email string) (bool, error) { return true, nil },
	}
	svc := NewUserService(repo, &stubPasswordHasher{})

	err := svc.Create(context.Background(), in.CreateUserCommand{
		Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Password: validPassword,
	})

	if err != nil {
		t.Fatalf("expected nil for an existing email (A5), got %v", err)
	}
	if repo.saveCalled {
		t.Fatal("expected repo.Save NOT to be called when email is taken")
	}
}

func TestCreate_InvalidData_DoesNotCallRepository(t *testing.T) {
	repo := &stubUserRepository{}
	svc := NewUserService(repo, &stubPasswordHasher{})

	err := svc.Create(context.Background(), in.CreateUserCommand{Email: "", FirstName: "", LastName: ""})

	var ve *domain.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %v", err)
	}
	if repo.saveCalled || repo.existsByEmailCalled {
		t.Fatal("expected repository not to be touched on invalid data")
	}
}

func TestCreate_HasherError_Propagates(t *testing.T) {
	boom := errors.New("hasher exploded")
	repo := &stubUserRepository{}
	hasher := &stubPasswordHasher{hashFunc: func(string) (string, error) { return "", boom }}
	svc := NewUserService(repo, hasher)

	err := svc.Create(context.Background(), in.CreateUserCommand{
		Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Password: validPassword,
	})

	if !errors.Is(err, boom) {
		t.Fatalf("expected hasher error to propagate, got %v", err)
	}
	if repo.saveCalled {
		t.Fatal("expected repo.Save NOT to be called when hashing fails")
	}
}

func TestGetByID_ZeroID_DoesNotTouchRepository(t *testing.T) {
	repo := &stubUserRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
			t.Fatal("repository should not be called with a zero ID")
			return nil, nil
		},
	}
	svc := NewUserService(repo, &stubPasswordHasher{})

	_, err := svc.GetByID(context.Background(), adminActor, uuid.Nil)
	if !errors.Is(err, domain.ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
}

func TestGetByID_NotFound_Propagates(t *testing.T) {
	repo := &stubUserRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		},
	}
	svc := NewUserService(repo, &stubPasswordHasher{})

	_, err := svc.GetByID(context.Background(), adminActor, uuid.New())
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdate_EmailChangedToTakenEmail(t *testing.T) {
	existing := &domain.User{ID: uuid.New(), Email: "old@b.com", FirstName: "Ana", LastName: "Pérez", IsActive: true}
	repo := &stubUserRepository{
		findByIDFunc:      func(ctx context.Context, id uuid.UUID) (*domain.User, error) { return existing, nil },
		existsByEmailFunc: func(ctx context.Context, email string) (bool, error) { return true, nil },
	}
	svc := NewUserService(repo, &stubPasswordHasher{})

	newEmail := "taken@b.com"
	_, err := svc.Update(context.Background(), in.Actor{ID: existing.ID, Role: domain.RoleUser}, in.UpdateUserCommand{ID: existing.ID, Email: &newEmail})

	if !errors.Is(err, domain.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUpdate_EmailUnchanged_DoesNotCheckExistence(t *testing.T) {
	existing := &domain.User{ID: uuid.New(), Email: "same@b.com", FirstName: "Ana", LastName: "Pérez", IsActive: true}
	repo := &stubUserRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) { return existing, nil },
	}
	svc := NewUserService(repo, &stubPasswordHasher{})

	newFirstName := "Ana María"
	_, err := svc.Update(context.Background(), in.Actor{ID: existing.ID, Role: domain.RoleUser}, in.UpdateUserCommand{ID: existing.ID, FirstName: &newFirstName})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if repo.existsByEmailCalled {
		t.Fatal("expected ExistsByEmail NOT to be called when email did not change")
	}
}

func TestDelete_ZeroID(t *testing.T) {
	svc := NewUserService(&stubUserRepository{}, &stubPasswordHasher{})
	if err := svc.Delete(context.Background(), uuid.Nil); !errors.Is(err, domain.ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
}

func TestList_ClampsLimit(t *testing.T) {
	cases := []struct {
		name          string
		limit, offset int
		expectLimit   int
		expectErr     bool
	}{
		{"zero limit defaults to 10", 0, 0, 10, false},
		{"limit above 100 clamps to 100", 500, 0, 100, false},
		{"negative limit errors", -1, 0, 0, true},
		{"negative offset errors", 10, -1, 0, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotLimit int
			repo := &stubUserRepository{
				listFunc: func(ctx context.Context, limit, offset int) ([]*domain.User, error) {
					gotLimit = limit
					return nil, nil
				},
			}
			svc := NewUserService(repo, &stubPasswordHasher{})

			_, err := svc.List(context.Background(), in.ListUsersQuery{Limit: tc.limit, Offset: tc.offset})

			if tc.expectErr {
				var ve *domain.ValidationError
				if !errors.As(err, &ve) {
					t.Fatalf("expected *ValidationError, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("List returned error: %v", err)
			}
			if gotLimit != tc.expectLimit {
				t.Fatalf("expected limit %d, got %d", tc.expectLimit, gotLimit)
			}
		})
	}
}

var adminActor = in.Actor{ID: uuid.New(), Role: domain.RoleAdmin}

// validPassword cumple la política de contraseñas de SPEC 01 y no contiene
// el correo ni el nombre de los usuarios de prueba.
const validPassword = "rio-cometa-azul"

func ptr[T any](v T) *T { return &v }

// --- Autorización (GHSA-j36w-x5rp-5f3w) ---

func TestCreate_AlwaysAssignsUserRole(t *testing.T) {
	var user *domain.User
	repo := &stubUserRepository{saveFunc: func(ctx context.Context, u *domain.User) error {
		user = u
		return nil
	}}
	svc := NewUserService(repo, &stubPasswordHasher{})

	err := svc.Create(context.Background(), in.CreateUserCommand{
		Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Password: validPassword,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if user.Role != domain.RoleUser {
		t.Fatalf("expected role %q, got %q", domain.RoleUser, user.Role)
	}
}

func TestGetByID_OtherUser_ReturnsNotFoundWithoutTouchingRepository(t *testing.T) {
	repo := &stubUserRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
			t.Fatal("repository should not be called for another user's ID")
			return nil, nil
		},
	}
	svc := NewUserService(repo, &stubPasswordHasher{})
	actor := in.Actor{ID: uuid.New(), Role: domain.RoleUser}

	_, err := svc.GetByID(context.Background(), actor, uuid.New())
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetByID_Self_Allowed(t *testing.T) {
	self := &domain.User{ID: uuid.New(), Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Role: domain.RoleUser}
	repo := &stubUserRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) { return self, nil },
	}
	svc := NewUserService(repo, &stubPasswordHasher{})

	got, err := svc.GetByID(context.Background(), in.Actor{ID: self.ID, Role: domain.RoleUser}, self.ID)
	if err != nil || got != self {
		t.Fatalf("expected own user, got %v, %v", got, err)
	}
}

func TestGetByID_Admin_CanReadAnyUser(t *testing.T) {
	other := &domain.User{ID: uuid.New(), Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Role: domain.RoleUser}
	repo := &stubUserRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) { return other, nil },
	}
	svc := NewUserService(repo, &stubPasswordHasher{})

	if _, err := svc.GetByID(context.Background(), adminActor, other.ID); err != nil {
		t.Fatalf("admin should read any user, got %v", err)
	}
}

func TestGetByID_ActorWithNilID_CannotReadAnyone(t *testing.T) {
	svc := NewUserService(&stubUserRepository{}, &stubPasswordHasher{})

	_, err := svc.GetByID(context.Background(), in.Actor{Role: domain.RoleUser}, uuid.New())
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdate_OtherUser_ReturnsNotFound(t *testing.T) {
	repo := &stubUserRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
			t.Fatal("repository should not be called for another user's ID")
			return nil, nil
		},
	}
	svc := NewUserService(repo, &stubPasswordHasher{})
	actor := in.Actor{ID: uuid.New(), Role: domain.RoleUser}

	_, err := svc.Update(context.Background(), actor, in.UpdateUserCommand{ID: uuid.New(), FirstName: ptr("X")})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdate_UserCannotChangeOwnRoleOrStatus(t *testing.T) {
	cases := map[string]in.UpdateUserCommand{
		"role":      {Role: ptr(domain.RoleAdmin)},
		"is_active": {IsActive: ptr(false)},
	}
	for name, cmd := range cases {
		t.Run(name, func(t *testing.T) {
			self := uuid.New()
			repo := &stubUserRepository{
				findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
					t.Fatal("repository should not be called when the change is forbidden")
					return nil, nil
				},
			}
			svc := NewUserService(repo, &stubPasswordHasher{})
			cmd.ID = self

			_, err := svc.Update(context.Background(), in.Actor{ID: self, Role: domain.RoleUser}, cmd)
			if !errors.Is(err, domain.ErrForbidden) {
				t.Fatalf("expected ErrForbidden, got %v", err)
			}
		})
	}
}

func TestUpdate_AdminCanPromoteUser(t *testing.T) {
	target := &domain.User{ID: uuid.New(), Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Role: domain.RoleUser, IsActive: true}
	repo := &stubUserRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) { return target, nil },
	}
	svc := NewUserService(repo, &stubPasswordHasher{})

	got, err := svc.Update(context.Background(), adminActor, in.UpdateUserCommand{ID: target.ID, Role: ptr(domain.RoleAdmin)})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if got.Role != domain.RoleAdmin {
		t.Fatalf("expected role admin, got %q", got.Role)
	}
}

func TestUpdate_InvalidRole_IsValidationError(t *testing.T) {
	target := &domain.User{ID: uuid.New(), Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Role: domain.RoleUser}
	repo := &stubUserRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) { return target, nil },
	}
	svc := NewUserService(repo, &stubPasswordHasher{})

	_, err := svc.Update(context.Background(), adminActor, in.UpdateUserCommand{ID: target.ID, Role: ptr("superadmin")})
	var ve *domain.ValidationError
	if !errors.As(err, &ve) || ve.Fields["role"] == "" {
		t.Fatalf("expected validation error on role, got %v", err)
	}
	if target.Role != domain.RoleUser {
		t.Fatalf("role must not change on validation error, got %q", target.Role)
	}
}

// --- SPEC 01: orden del registro ---

func TestCreate_InvalidPassword_IsValidationErrorWithoutHashing(t *testing.T) {
	repo := &stubUserRepository{}
	hasherCalled := false
	hasher := &stubPasswordHasher{hashFunc: func(string) (string, error) {
		hasherCalled = true
		return "hashed", nil
	}}
	svc := NewUserService(repo, hasher)

	err := svc.Create(context.Background(), in.CreateUserCommand{
		Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Password: "corta",
	})

	var ve *domain.ValidationError
	if !errors.As(err, &ve) || ve.Fields["password"] == "" {
		t.Fatalf("expected validation error on password, got %v", err)
	}
	if hasherCalled || repo.existsByEmailCalled || repo.saveCalled {
		t.Fatal("expected neither hasher nor repository to be called on invalid data")
	}
}

func TestCreate_MergesUserAndPasswordFieldErrors(t *testing.T) {
	svc := NewUserService(&stubUserRepository{}, &stubPasswordHasher{})

	err := svc.Create(context.Background(), in.CreateUserCommand{
		Email: "no-es-correo", FirstName: "", LastName: "Pérez", Password: "corta",
	})

	var ve *domain.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %v", err)
	}
	for _, field := range []string{"email", "first_name", "password"} {
		if ve.Fields[field] == "" {
			t.Errorf("expected field %q in %v", field, ve.Fields)
		}
	}
}

func TestCreate_ExistingEmail_HashesBeforeCheckingExistence(t *testing.T) {
	var calls []string
	repo := &stubUserRepository{
		existsByEmailFunc: func(ctx context.Context, email string) (bool, error) {
			calls = append(calls, "exists")
			return true, nil
		},
	}
	hasher := &stubPasswordHasher{hashFunc: func(plain string) (string, error) {
		calls = append(calls, "hash")
		return "hashed", nil
	}}
	svc := NewUserService(repo, hasher)

	_ = svc.Create(context.Background(), in.CreateUserCommand{
		Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Password: validPassword,
	})

	if len(calls) != 2 || calls[0] != "hash" || calls[1] != "exists" {
		t.Fatalf("expected exactly [hash exists], got %v", calls)
	}
	if repo.saveCalled {
		t.Fatal("expected repo.Save NOT to be called when email exists")
	}
}

func TestCreate_RoleOfSavedUserIsUser(t *testing.T) {
	var saved *domain.User
	repo := &stubUserRepository{saveFunc: func(ctx context.Context, u *domain.User) error {
		saved = u
		return nil
	}}
	svc := NewUserService(repo, &stubPasswordHasher{})

	if err := svc.Create(context.Background(), in.CreateUserCommand{
		Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Password: validPassword,
	}); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if saved == nil || saved.Role != domain.RoleUser {
		t.Fatalf("expected saved user with role %q, got %+v", domain.RoleUser, saved)
	}
	if saved.EmailVerifiedAt != nil {
		t.Fatalf("expected EmailVerifiedAt nil on registration, got %v", saved.EmailVerifiedAt)
	}
}

// --- SPEC 01: registro indistinguible ---

func TestCreate_SaveEmailTaken_ReturnsNil(t *testing.T) {
	repo := &stubUserRepository{saveFunc: func(ctx context.Context, u *domain.User) error {
		return domain.ErrEmailTaken
	}}
	svc := NewUserService(repo, &stubPasswordHasher{})

	err := svc.Create(context.Background(), in.CreateUserCommand{
		Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Password: validPassword,
	})
	if err != nil {
		t.Fatalf("expected nil when Save reports a concurrent duplicate, got %v", err)
	}
	if !repo.saveCalled {
		t.Fatal("expected repo.Save to be called")
	}
}

func TestCreate_SaveOtherError_Propagates(t *testing.T) {
	boom := errors.New("db down")
	repo := &stubUserRepository{saveFunc: func(ctx context.Context, u *domain.User) error { return boom }}
	svc := NewUserService(repo, &stubPasswordHasher{})

	err := svc.Create(context.Background(), in.CreateUserCommand{
		Email: "a@b.com", FirstName: "Ana", LastName: "Pérez", Password: validPassword,
	})
	if !errors.Is(err, boom) {
		t.Fatalf("expected Save error to propagate, got %v", err)
	}
}

func TestCreate_NeverLogsTheEmail(t *testing.T) {
	const email = "registro-privado@ejemplo.co"
	repos := map[string]*stubUserRepository{
		"correo existente": {existsByEmailFunc: func(ctx context.Context, e string) (bool, error) { return true, nil }},
		"carrera en Save":  {saveFunc: func(ctx context.Context, u *domain.User) error { return domain.ErrEmailTaken }},
		"registro nuevo":   {},
	}
	for name, repo := range repos {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			prev, prevFlags := log.Writer(), log.Flags()
			log.SetOutput(&buf)
			t.Cleanup(func() { log.SetOutput(prev); log.SetFlags(prevFlags) })

			svc := NewUserService(repo, &stubPasswordHasher{})
			if err := svc.Create(context.Background(), in.CreateUserCommand{
				Email: email, FirstName: "Ana", LastName: "Pérez", Password: validPassword,
			}); err != nil {
				t.Fatalf("Create returned error: %v", err)
			}
			if strings.Contains(buf.String(), email) {
				t.Fatalf("email leaked to logs: %q", buf.String())
			}
		})
	}
}
