package application

import (
	"context"
	"errors"
	"testing"

	"finanzas-api/internal/auth/domain"
	"finanzas-api/internal/auth/port/in"
)

type stubUserFinder struct {
	findByEmailFunc func(ctx context.Context, email string) (*domain.Credentials, error)
}

func (s *stubUserFinder) FindByEmail(ctx context.Context, email string) (*domain.Credentials, error) {
	return s.findByEmailFunc(ctx, email)
}

type stubPasswordVerifier struct {
	matchesFunc func(plain, hash string) bool
}

func (s *stubPasswordVerifier) Matches(plain, hash string) bool {
	if s.matchesFunc != nil {
		return s.matchesFunc(plain, hash)
	}
	return true
}

type stubTokenProviderApp struct {
	issueFunc func(claims domain.Claims) (domain.Token, error)
}

func (s *stubTokenProviderApp) Issue(claims domain.Claims) (domain.Token, error) {
	if s.issueFunc != nil {
		return s.issueFunc(claims)
	}
	return domain.Token{Value: "token-for-" + claims.Subject}, nil
}

func (s *stubTokenProviderApp) Verify(raw string) (domain.Claims, error) {
	return domain.Claims{}, nil
}

var activeCredentials = &domain.Credentials{
	UserID: "user-1", Email: "a@b.com", PasswordHash: "hashed", Role: "admin", Active: true,
}

func TestLogin_Success_IssuesTokenWithUserAndRole(t *testing.T) {
	var issuedClaims domain.Claims
	finder := &stubUserFinder{findByEmailFunc: func(ctx context.Context, email string) (*domain.Credentials, error) {
		return activeCredentials, nil
	}}
	tokens := &stubTokenProviderApp{issueFunc: func(claims domain.Claims) (domain.Token, error) {
		issuedClaims = claims
		return domain.Token{Value: "ok"}, nil
	}}
	svc := NewAuthService(finder, &stubPasswordVerifier{}, tokens)

	token, err := svc.Login(context.Background(), in.LoginCommand{Email: "a@b.com", Password: "secret1"})
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if token.Value != "ok" {
		t.Fatalf("expected token value 'ok', got %q", token.Value)
	}
	if issuedClaims.Subject != "user-1" || issuedClaims.Role != "admin" {
		t.Fatalf("expected claims with the credentials' user and role, got %+v", issuedClaims)
	}
}

func TestLogin_UnknownEmail_And_WrongPassword_ReturnSameError(t *testing.T) {
	unknownEmail := &stubUserFinder{findByEmailFunc: func(ctx context.Context, email string) (*domain.Credentials, error) {
		return nil, domain.ErrInvalidCredentials
	}}
	svc1 := NewAuthService(unknownEmail, &stubPasswordVerifier{}, &stubTokenProviderApp{})
	_, err1 := svc1.Login(context.Background(), in.LoginCommand{Email: "ghost@b.com", Password: "whatever"})

	wrongPassword := &stubUserFinder{findByEmailFunc: func(ctx context.Context, email string) (*domain.Credentials, error) {
		return activeCredentials, nil
	}}
	svc2 := NewAuthService(wrongPassword, &stubPasswordVerifier{matchesFunc: func(string, string) bool { return false }}, &stubTokenProviderApp{})
	_, err2 := svc2.Login(context.Background(), in.LoginCommand{Email: "a@b.com", Password: "wrong"})

	if !errors.Is(err1, domain.ErrInvalidCredentials) || !errors.Is(err2, domain.ErrInvalidCredentials) {
		t.Fatalf("expected both to be ErrInvalidCredentials, got %v and %v", err1, err2)
	}
}

func TestLogin_InactiveUser(t *testing.T) {
	inactive := &domain.Credentials{UserID: "user-1", Email: "a@b.com", PasswordHash: "hashed", Role: "user", Active: false}
	finder := &stubUserFinder{findByEmailFunc: func(ctx context.Context, email string) (*domain.Credentials, error) {
		return inactive, nil
	}}
	svc := NewAuthService(finder, &stubPasswordVerifier{}, &stubTokenProviderApp{})

	_, err := svc.Login(context.Background(), in.LoginCommand{Email: "a@b.com", Password: "secret1"})
	if !errors.Is(err, domain.ErrUserInactive) {
		t.Fatalf("expected ErrUserInactive, got %v", err)
	}
}

func TestLogin_TokenProviderError_Propagates(t *testing.T) {
	boom := errors.New("token provider exploded")
	finder := &stubUserFinder{findByEmailFunc: func(ctx context.Context, email string) (*domain.Credentials, error) {
		return activeCredentials, nil
	}}
	tokens := &stubTokenProviderApp{issueFunc: func(claims domain.Claims) (domain.Token, error) {
		return domain.Token{}, boom
	}}
	svc := NewAuthService(finder, &stubPasswordVerifier{}, tokens)

	_, err := svc.Login(context.Background(), in.LoginCommand{Email: "a@b.com", Password: "secret1"})
	if !errors.Is(err, boom) {
		t.Fatalf("expected token provider error to propagate, got %v", err)
	}
}
