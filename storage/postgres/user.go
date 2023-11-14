package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/husanmusa/auth_pro_service/genproto/auth_service"
	"github.com/husanmusa/auth_pro_service/pkg/helper"
	"github.com/husanmusa/auth_pro_service/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) storage.UserI {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUser(ctx context.Context, req *pb.User) error {
	req.Id = uuid.NewString()

	_, err := r.db.Exec(ctx, `
	insert into users (
		id,
	    name,
		username,
	    password,
		role
	) values ($1,$2,$3,$4,$5)`,
		req.Id,
		req.Name,
		req.Username,
		req.Password,
		req.Role,
	)
	if err != nil {
		return fmt.Errorf("create backdoor %w", err)
	}
	return nil
}
func (r *userRepo) GetUser(ctx context.Context, req *pb.ById) (*pb.User, error) {
	var (
		user pb.User
	)
	err := r.db.QueryRow(ctx, `
	select
		id,
		name,
		username,
		role
	from users
	where id = $1`, req.Id).Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("get backdoor %w", err)
	}

	return &user, nil
}
func (r *userRepo) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	var (
		query = `
		select
			id,
			name,
			username,
			role
		from users`
		filter = " WHERE true "
		page   = " OFFSET 0"
		limit  = " LIMIT 10"
		order  = " ORDER BY created_at DESC"
		params = make(map[string]interface{})
		resp   pb.GetUserListResponse
	)

	if req.Limit > 0 {
		params["limit"] = req.Limit
		limit = ` LIMIT :limit`
	}

	if req.Offset > 0 {
		params["offset"] = req.Offset
		page = ` OFFSET :offset`
	}

	query += filter + order + limit + page
	query, arr := helper.ReplaceQueryParams(query, params)

	rows, err := r.db.Query(ctx, query, arr...)
	if err != nil {
		return nil, fmt.Errorf("get all user query %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			user pb.User
		)
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Username,
			&user.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("get backdoor %w", err)
		}

		resp.Users = append(resp.Users, &user)
	}
	query = `select count(1) from users` + filter
	query, arr = helper.ReplaceQueryParams(query, params)
	err = r.db.QueryRow(ctx, query, arr...).Scan(
		&resp.Count,
	)
	if err != nil {
		return nil, fmt.Errorf("get all user count %w", err)
	}
	return &resp, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, in *pb.User) error {
	query := `update users set
			name= :name,
			username= :username,
			role= :role
	   where
		id = :id`

	params := map[string]interface{}{
		"name":     in.Name,
		"username": in.Username,
		"role":     in.Role,
	}

	q, arr := helper.ReplaceQueryParams(query, params)
	result, err := r.db.Exec(ctx, q, arr...)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, in *pb.ById) error {
	query := `delete from users where id = $1`

	result, err := r.db.Exec(ctx, query, in.Id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*pb.User, error) {
	var (
		user pb.User
	)
	err := r.db.QueryRow(ctx, `
	select
		id,
		name,
		username,
		password,
		role
	from users
	where username = $1`, username).Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("get by username %w", err)
	}

	return &user, nil
}
