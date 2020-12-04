package db

import (
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4"
	"github.com/maxshend/tiny_goauth/models"
)

// CreateUser creates a new record in users database table
func (s *datastore) CreateUser(user *models.User) error {
	tr, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tr.Rollback(ctx)

	err = tr.QueryRow(
		ctx,
		"INSERT INTO users(email, password) VALUES($1, $2) RETURNING id, created_at",
		user.Email, user.Password,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}

	for _, role := range user.Roles {
		batch.Queue(
			"INSERT INTO user_roles(user_id, role_id) VALUES($1, (SELECT id FROM roles WHERE name = $2))",
			user.ID, role,
		)
	}

	br := tr.SendBatch(ctx, batch)

	for i := 0; i < len(user.Roles); i++ {
		ct, err := br.Exec()
		if err != nil {
			return err
		}
		if ct.RowsAffected() != 1 {
			return zeroInsertedRows
		}
	}

	if err = br.Close(); err != nil {
		return err
	}
	if err = tr.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *datastore) UserExistsWithField(fl validator.FieldLevel) (bool, error) {
	result := false
	err := s.db.QueryRow(
		ctx,
		"SELECT EXISTS (SELECT 1 FROM users WHERE "+fl.FieldName()+" = $1)", fl.Field().String(),
	).Scan(&result)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (s *datastore) UserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		ctx,
		"SELECT users.id AS id, email, password, created_at, ARRAY_REMOVE(ARRAY_AGG(roles.name), NULL) AS roles FROM users "+
			"LEFT JOIN user_roles ON users.id = user_roles.user_id "+
			"LEFT JOIN roles ON user_roles.role_id = roles.id WHERE email = $1 GROUP BY users.id "+
			"LIMIT 1",
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.Roles)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *datastore) DeleteUser(id int64) error {
	commandTag, err := s.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return zeroDeleteRows
	}

	return nil
}

func (s *datastore) GetRoles() (roles []string, err error) {
	rows, err := s.db.Query(ctx, "SELECT name FROM roles")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var role string
		if err = rows.Scan(&role); err != nil {
			return
		}

		roles = append(roles, role)
	}

	return
}

func (s *datastore) CreateRole(name string) error {
	ct, err := s.db.Exec(ctx, "INSERT INTO roles(name) VALUES($1)", name)
	if err != nil {
		return err
	}
	if ct.RowsAffected() != 1 {
		return zeroInsertedRows
	}

	return nil
}
