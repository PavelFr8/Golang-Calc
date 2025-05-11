package orchestrator

import (
	"github.com/PavelFr8/Golang-Calc/pkg/hash"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateExpression(expr *Expression) error {
	return r.db.Create(expr).Error
}

func (r *Repository) CreateTask(task *Task) error {
	return r.db.Create(task).Error
}

func (r *Repository) CreateUser(login string, password string) error {
    hash_password, err := hash.Generate(password)
    if err != nil {
        return err
    }
    user := &User{
        Login: login,
        Password: hash_password,
    }
	return r.db.Create(user).Error
}

func (r *Repository) GetMaxTaskID() uint {
    var maxIDObject Task
    err := r.db.Order("id desc").First(&maxIDObject).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return 0 
        }
        panic("FIND MAX ID ERROR: " + err.Error())
    }
    return maxIDObject.ID
}

func (r *Repository) GetMaxExpressionID() uint {
    var maxIDObject Expression
    err := r.db.Order("id desc").First(&maxIDObject).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return 0 
        }
        panic("FIND MAX ID ERROR: " + err.Error())
    }
    return maxIDObject.ID
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&Expression{}, &Task{}, &User{}); err != nil {
		panic("failed to migrate database")
	}
	return db
}
