package repository

import (
	"context"
)

func (r *Repository) CreateUser(ctx context.Context, input CreateUserInput) (output User, err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return
	}
	defer tx.Commit()

	query := `INSERT INTO users(full_name, phone_number, password) values($1, $2, $3) RETURNING id`
	err = tx.QueryRowContext(ctx, query, input.FullName, input.PhoneNumber, input.Password).Scan(&output.Id)
	if err != nil {
		tx.Rollback()
		return
	}
	return
}

func (r *Repository) UpdateUserById(ctx context.Context, input UpdateUserInput) (output User, err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return
	}
	defer tx.Commit()

	query := `UPDATE users SET full_name=$1, phone_number=$2, updated_at=NOW() WHERE id = $3 RETURNING id,full_name,phone_number`
	err = r.Db.QueryRowContext(ctx, query, input.FullName, input.PhoneNumber, input.Id).Scan(
		&output.Id,
		&output.FullName,
		&output.PhoneNumber,
	)
	if err != nil {
		tx.Rollback()
		return
	}
	return
}

func (r *Repository) GetUserById(ctx context.Context, id int) (output User, err error) {
	query := `SELECT id, full_name, phone_number, password, updated_at, created_at FROM users WHERE id = $1`
	err = r.Db.QueryRowContext(ctx, query, id).Scan(
		&output.Id,
		&output.FullName,
		&output.PhoneNumber,
		&output.Password,
		&output.UpdatedAt,
		&output.CreatedAt,
	)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (output User, err error) {
	query := `SELECT id, full_name, phone_number, password, updated_at, created_at FROM users u WHERE phone_number = $1`
	err = r.Db.QueryRowContext(ctx, query, phoneNumber).Scan(
		&output.Id,
		&output.FullName,
		&output.PhoneNumber,
		&output.Password,
		&output.UpdatedAt,
		&output.CreatedAt,
	)
	if err != nil {
		return
	}
	return
}
