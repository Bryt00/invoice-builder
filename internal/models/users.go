package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type User struct {
	Name              string     `gorm:"type:varchar(255); not null" json:"name"`
	Email             string     `gorm:"unique;type:citext;not null" json:"email"`
	PasswordHash      string     `gorm:"type:varchar(255); not null" json:"-"`
	IsProfileComplete bool       `gorm:"default:false" json:"is_profile_complete"`
	IsActivated       bool       `gorm:"default:false" json:"is_activated"`
	ActivationToken   string     `gorm:"type:varchar(255)" json:"-"`
	ActivationExpiry  *time.Time `json:"-"`
	ResetToken        string     `gorm:"type:varchar(255)" json:"-"`
	ResetExpiry       *time.Time `json:"-"`
	TokenVersion      int        `gorm:"default:0;not null" json:"-"`
	RoleID            uuid.UUID  `gorm:"type:uuid;not null" json:"role_id"`
	Role              Role       `json:"role"`
	BaseModel
}

func (u *User) IsAdmin() bool {
	if u == nil {
		return false
	}
	return u.Role.Name == "admin" || u.Role.Name == "superadmin" || u.Role.Name == "Admin"
}
type UserInterface interface {
	Insert(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error

	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByActivationToken(ctx context.Context, token string) (*User, error)
	GetByResetToken(ctx context.Context, token string) (*User, error)
	ActivateUser(ctx context.Context, id uuid.UUID) error
	Authenticate(ctx context.Context, email, plaintxtPassword string) (uuid.UUID, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	GetAllUsers(ctx context.Context, search string, page, limit int) ([]*User, int64, error)
	UpdateRole(ctx context.Context, userID, roleID uuid.UUID) error
	ToggleActivation(ctx context.Context, userID uuid.UUID, isActivated bool) error
	GetSystemUserStats(ctx context.Context) (total, active, unverified int64, err error)
	IncrementTokenVersion(ctx context.Context, id uuid.UUID) error
}

type UserModel struct {
	DB *gorm.DB
}

func (m *UserModel) Insert(ctx context.Context, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 12)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	user.ID = uuid.New()

	// Assign default "User" role if none is set
	if user.RoleID == uuid.Nil {
		var defaultRole Role
		if err := m.DB.WithContext(ctx).Where("name = ?", "User").First(&defaultRole).Error; err == nil {
			user.RoleID = defaultRole.ID
		}
	}

	createdUser := m.DB.WithContext(ctx).Create(user).Error
	if createdUser != nil {
		var pgErr *pgconn.PgError
		if errors.As(createdUser, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateEmail
		}
		return createdUser
	}
	return nil
}
func (m *UserModel) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	result := m.DB.WithContext(ctx).Preload("Role").First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, fmt.Errorf("db err: %w", result.Error)
	}
	return &user, nil
}
func (m *UserModel) Update(ctx context.Context, user *User) error {
	return m.DB.WithContext(ctx).Save(user).Error
}
func (m *UserModel) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// GetByEmail Retrieves the user email
func (m *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	res := m.DB.First(&user, "email = ?", email)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, fmt.Errorf("db err: %w", res.Error)
	}

	return &user, nil
}

// GetByActivationToken retrieves the user by active activation token
func (m *UserModel) GetByActivationToken(ctx context.Context, token string) (*User, error) {
	var user User
	res := m.DB.WithContext(ctx).First(&user, "activation_token = ? AND activation_expiry > ?", token, time.Now())
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, fmt.Errorf("db err: %w", res.Error)
	}
	return &user, nil
}

// ActivateUser marks the user as activated and clears activation token fields
func (m *UserModel) ActivateUser(ctx context.Context, id uuid.UUID) error {
	res := m.DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(map[string]any{
		"is_activated":      true,
		"activation_token":  nil,
		"activation_expiry": nil,
	})
	return res.Error
}

// GetByResetToken retrieves the user by active reset token
func (m *UserModel) GetByResetToken(ctx context.Context, token string) (*User, error) {
	var user User
	res := m.DB.WithContext(ctx).First(&user, "reset_token = ? AND reset_expiry > ?", token, time.Now())
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, fmt.Errorf("db err: %w", res.Error)
	}
	return &user, nil
}

func (m *UserModel) Authenticate(ctx context.Context, email, plaintxtPassword string) (uuid.UUID, error) {
	user, err := m.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return uuid.Nil, ErrInvalidCredentials
		}
		return uuid.Nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plaintxtPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return uuid.Nil, ErrInvalidCredentials
		}
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (m *UserModel) UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	res := m.DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(map[string]any{
		"password_hash": string(hashedPassword),
		"reset_token":   nil,
		"reset_expiry":  nil,
	})
	if res.Error != nil {
		return res.Error
	}

	// Invalidate all existing tokens by bumping the token version.
	return m.IncrementTokenVersion(ctx, id)
}

// ExistsByID Check if user already exists
func (m *UserModel) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	user, err := m.GetByID(ctx, id)
	if err != nil {
		return false, nil
	}
	return user != nil, nil
}

func (m *UserModel) GetAllUsers(ctx context.Context, search string, page, limit int) ([]*User, int64, error) {
	var users []*User
	var total int64

	query := m.DB.WithContext(ctx).Model(&User{}).Preload("Role")
	if search != "" {
		s := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", s, s)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (m *UserModel) UpdateRole(ctx context.Context, userID, roleID uuid.UUID) error {
	err := m.DB.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("role_id", roleID).Error
	if err != nil {
		return err
	}

	// Invalidate all existing tokens since the role changed.
	return m.IncrementTokenVersion(ctx, userID)
}

// IncrementTokenVersion bumps the user's token_version to invalidate all
// previously issued JWTs.
func (m *UserModel) IncrementTokenVersion(ctx context.Context, id uuid.UUID) error {
	return m.DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		UpdateColumn("token_version", gorm.Expr("token_version + 1")).Error
}

func (m *UserModel) ToggleActivation(ctx context.Context, userID uuid.UUID, isActivated bool) error {
	return m.DB.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("is_activated", isActivated).Error
}

func (m *UserModel) GetSystemUserStats(ctx context.Context) (total, active, unverified int64, err error) {
	if err = m.DB.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
		return 0, 0, 0, err
	}
	if err = m.DB.WithContext(ctx).Model(&User{}).Where("is_activated = ?", true).Count(&active).Error; err != nil {
		return 0, 0, 0, err
	}
	unverified = total - active
	return total, active, unverified, nil
}
