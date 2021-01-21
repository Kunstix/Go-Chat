package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/kunstix/gochat/models"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (user *User) GetId() string {
	return user.Id
}

func (user *User) GetName() string {
	return user.Name
}

type UserRepository struct {
	Db *sql.DB
}

func (repo *UserRepository) AddUser(user models.User) error {
	existingUser := repo.FindUserByUsername(user.GetName())

	if existingUser != nil {
		stmt, err := repo.Db.Prepare("INSERT INTO user(id, name) values(?,?)")
		checkErr(err)

		_, err = stmt.Exec(user.GetId(), user.GetName())
		checkErr(err)
		return nil
	} else {
		return errors.New("User already exists")
	}
}

func (repo *UserRepository) AddRegisteredUser(name string, password string) {
	stmt, err := repo.Db.Prepare("INSERT into user (id, name, password) VALUES(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(uuid.New().String(), name, password)
	checkErr(err)
}

func (repo *UserRepository) RemoveUser(user models.User) {
	stmt, err := repo.Db.Prepare("DELETE FROM user WHERE id = ? AND password IS NOT NULL")
	checkErr(err)

	_, err = stmt.Exec(user.GetId())
	checkErr(err)
}

func (repo *UserRepository) FindUserById(ID string) models.User {

	row := repo.Db.QueryRow("SELECT id, name FROM user where id = ? LIMIT 1", ID)

	var user User
	if err := row.Scan(&user.Id, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return &user
}

func (repo *UserRepository) FindUserByUsername(username string) *User {

	row := repo.Db.QueryRow("SELECT id, name, password FROM user where name = ? LIMIT 1", username)

	var user User
	if err := row.Scan(&user.Id, &user.Name, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}
	return &user
}

func (repo *UserRepository) GetAllRegisteredUsers() []models.User {

	rows, err := repo.Db.Query("SELECT id, name FROM user WHERE password IS NULL")

	if err != nil {
		log.Fatal(err)
	}
	var users []models.User
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name)
		users = append(users, &user)
	}

	return users
}

func (repo *UserRepository) GetAllUsers() []models.User {

	rows, err := repo.Db.Query("SELECT id, name FROM user")

	if err != nil {
		log.Fatal(err)
	}
	var users []models.User
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name)
		users = append(users, &user)
	}

	return users
}
