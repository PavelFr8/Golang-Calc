package orchestrator

import (
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