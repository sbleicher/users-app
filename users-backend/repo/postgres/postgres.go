package postgres

import (
	"fmt"
	"os"
	"users-backend/model"
	"users-backend/repo"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var (
	_ repo.UserRepo = new(PostgresRepo)
)

type PostgresRepo struct {
	db *pg.DB
}

func NewPostgresRepo() (*PostgresRepo, func()) {
	opt, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db := pg.Connect(opt)

	err = createSchema(db)
	if err != nil {
		panic(err)
	}

	// Return the repo and a cleanup function to close the connection
	return &PostgresRepo{db: db}, func() {
		if err := db.Close(); err != nil {
			// log.Fatalf("Failed to close DB connection: %v", err)
			fmt.Printf("Failed to close DB connection: %v", err)
		}
		// log.Println("Database connection closed.")
		fmt.Println("Database connection closed.")
	}
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*model.User)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        true,
			IfNotExists: true,
			Varchar:     255,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresRepo) GetById(user_id int) (*model.User, error) {
	var user model.User
	err := r.db.Model(&user).
		Where("user_id = ?", user_id).
		Select()
	if err != nil {
		print(err)
		return nil, err
		// return nil, ErrUserNotFound{
		// 	Message: fmt.Sprintf("User with username %s is not found", username),
		// }
	}
	return &user, nil
}

func (r *PostgresRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Model(&user).
		Where("user_name = ?", username).
		Select()
	if err != nil {
		return nil, err
		// return nil, ErrUserNotFound{
		// 	Message: fmt.Sprintf("User with username %s is not found", username),
		// }
	}
	return &user, nil
}

func (r *PostgresRepo) GetAll() (*[]model.User, error) {
	var users []model.User
	err := r.db.Model(&users).Select()
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *PostgresRepo) Create(user *model.User) (int, error) {
	_, err := r.db.Model(user).Insert()
	if err != nil {
		return -1, err
	}
	return user.UserID, nil
}

func (r *PostgresRepo) Update(user *model.User) (int, error) {
	u := &model.User{UserID: user.UserID}
	err := r.db.Model(u).WherePK().Select()
	if err != nil {
		return -1, err
	}

	u.UserName = user.UserName
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Email = user.Email
	u.UserStatus = user.UserStatus
	u.Department = user.Department

	_, err = r.db.Model(u).WherePK().Update()
	if err != nil {
		return -1, err
	}

	return u.UserID, nil
}

func (r *PostgresRepo) Delete(user_id int) error {
	user := &model.User{UserID: user_id}
	_, err := r.db.Model(user).WherePK().Delete()
	return err
}
