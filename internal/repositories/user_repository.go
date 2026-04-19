package repositories

import (
	"database/sql"
	"time"

	"iiceekiingfx.com/internal/models"
)

type UserRepository struct {
	db *Database
}

func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password, first_name, last_name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.DB.Exec(query,
		user.ID, user.Email, user.Password,
		user.FirstName, user.LastName, user.Role,
		time.Now(), time.Now(),
	)
	
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at
		FROM users WHERE email = $1
	`
	
	user := &models.User{}
	err := r.db.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password,
		&user.FirstName, &user.LastName, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return user, err
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at
		FROM users WHERE id = $1
	`
	
	user := &models.User{}
	err := r.db.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Password,
		&user.FirstName, &user.LastName, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return user, err
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET first_name = $2, last_name = $3, role = $4, updated_at = $5
		WHERE id = $1
	`
	
	_, err := r.db.DB.Exec(query,
		user.ID, user.FirstName, user.LastName,
		user.Role, time.Now(),
	)
	
	return err
}
