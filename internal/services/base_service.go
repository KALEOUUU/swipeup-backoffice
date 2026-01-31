package services

import (
	"gorm.io/gorm"
)

// BaseService provides common CRUD operations to avoid code duplication
type BaseService[T any] struct {
	db *gorm.DB
}

func NewBaseService[T any](db *gorm.DB) *BaseService[T] {
	return &BaseService[T]{db: db}
}

func (s *BaseService[T]) Create(entity *T) error {
	return s.db.Create(entity).Error
}

func (s *BaseService[T]) Update(entity *T) error {
	return s.db.Session(&gorm.Session{FullSaveAssociations: false}).Save(entity).Error
}

func (s *BaseService[T]) Delete(id uint) error {
	var entity T
	return s.db.Delete(&entity, id).Error
}

func (s *BaseService[T]) FindByID(id uint, preloads ...string) (*T, error) {
	var entity T
	query := s.db
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *BaseService[T]) FindAll(preloads ...string) ([]T, error) {
	var entities []T
	query := s.db
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.Find(&entities).Error
	return entities, err
}

func (s *BaseService[T]) FindWithCondition(condition interface{}, preloads ...string) ([]T, error) {
	var entities []T
	query := s.db
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.Where(condition).Find(&entities).Error
	return entities, err
}

func (s *BaseService[T]) Count(condition interface{}) (int64, error) {
	var count int64
	var entity T
	err := s.db.Model(&entity).Where(condition).Count(&count).Error
	return count, err
}

func (s *BaseService[T]) GetDB() *gorm.DB {
	return s.db
}
