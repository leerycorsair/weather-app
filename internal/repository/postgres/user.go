package postgres

import (
	"errors"
	"fmt"
	"weather-app/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("insert into %s (login, password, email) values ($1, $2, $3) returning id", UsersTable)
	row := r.db.QueryRow(query, user.Login, user.Password, user.Email)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) GetUser(login, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("select id, login, password, email from %s where login=$1 and password=$2", UsersTable)
	err := r.db.Get(&user, query, login, password)
	return user, err
}

func (r *UserRepository) GetFavorites(userId int) ([]int, error) {
	var favorites []int
	query := fmt.Sprintf("select city_id from %s where user_id=$1", UsersCitiesTable)
	err := r.db.Select(&favorites, query, userId)
	return favorites, err
}

func (r *UserRepository) AddFavorite(userId int, cityId int) (int, error) {
	var id int
	query := fmt.Sprintf("insert into %s (user_id, city_id) values ($1, $2) returning id", UsersCitiesTable)
	row := r.db.QueryRow(query, userId, cityId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) DeleteFavorite(userId int, cityId int) error {
	query := fmt.Sprintf("delete from %s where user_id=$1 and city_id=$2", UsersCitiesTable)
	result, err := r.db.Exec(query, userId, cityId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows deleted")
	}

	return nil
}
