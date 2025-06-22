package repository

import (
	"errors"
	"finanzas-api/internal/users/domain"
	"sync"
	"time"
)

type userRepositoryMemory struct {
	users  map[uint]*domain.User
	emails map[string]uint
	nextID uint
	mutex  sync.RWMutex
}

func NewUserMemoryRepository() domain.UserRepository {
	return &userRepositoryMemory{
		users:  make(map[uint]*domain.User),
		emails: make(map[string]uint),
		nextID: 1,
	}
}

func (r *userRepositoryMemory) Create(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Verificar si el email ya existe
	if _, exists := r.emails[user.Email]; exists {
		return errors.New("email already exists")
	}

	// Asignar ID y timestamps
	user.ID = r.nextID
	r.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Guardar usuario
	r.users[user.ID] = user
	r.emails[user.Email] = user.ID

	return nil
}

func (r *userRepositoryMemory) GetByID(id uint) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Verificar soft delete
	if !user.DeletedAt.Time.IsZero() {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *userRepositoryMemory) GetByEmail(email string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	userID, exists := r.emails[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	user := r.users[userID]
	if !user.DeletedAt.Time.IsZero() {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *userRepositoryMemory) Update(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	existingUser, exists := r.users[user.ID]
	if !exists || !existingUser.DeletedAt.Time.IsZero() {
		return errors.New("user not found")
	}

	// Si cambió el email, actualizar índice
	if existingUser.Email != user.Email {
		// Verificar que el nuevo email no exista
		if _, emailExists := r.emails[user.Email]; emailExists {
			return errors.New("email already exists")
		}

		// Actualizar índice de emails
		delete(r.emails, existingUser.Email)
		r.emails[user.Email] = user.ID
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user

	return nil
}

func (r *userRepositoryMemory) Delete(id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user, exists := r.users[id]
	if !exists || !user.DeletedAt.Time.IsZero() {
		return errors.New("user not found")
	}

	// Soft delete
	user.DeletedAt.Time = time.Now()
	user.DeletedAt.Valid = true
	user.UpdatedAt = time.Now()

	return nil
}

func (r *userRepositoryMemory) List(limit, offset int) ([]*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var users []*domain.User
	count := 0
	skipped := 0

	for _, user := range r.users {
		// Saltar usuarios eliminados
		if !user.DeletedAt.Time.IsZero() {
			continue
		}

		// Aplicar offset
		if skipped < offset {
			skipped++
			continue
		}

		// Aplicar limit
		if limit > 0 && count >= limit {
			break
		}

		users = append(users, user)
		count++
	}

	return users, nil
}

func (r *userRepositoryMemory) EmailExists(email string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	userID, exists := r.emails[email]
	if !exists {
		return false, nil
	}

	user := r.users[userID]
	// Solo existe si no está eliminado
	return user.DeletedAt.Time.IsZero(), nil
}
