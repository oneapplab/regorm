package regorm

import (
	"gorm.io/gorm"
)

// IRepository a generic interface for repositories
// [T IBaseModel] is generic type which T is type based on IBaseModel interface
// sample:
//
//	type ISampleRepository interface {
//		data.IRepository[SampleModel]
//	}
type IRepository[T IBaseModel] interface {
	First(model *T, conds interface{}) error // Select query based on id
	Create(model *T) (*T, error)             // Insert model
	BatchCreate(models []*T) (int64, error)  // Batch Insert based on slice of model
	Update(model *T) error                   // Update a model
	Delete(model *T) (int64, error)          // Delete a record
	GetDB() *gorm.DB                         // Get Database Instance
}

// Repository a generic struct which should be embed by other repositories
// to access repository methods
//
// sample usage as embed repository:
//
//	type SampleRepository struct {
//		Repository[SampleModel]
//	}
//
// sample usage as declare as repository:
//
// sampleRepository := InitRepository[SampleModel](db)
//

type Repository[T IBaseModel] struct {
	IRepository[T]

	database *gorm.DB
}

func InitRepository[T IBaseModel](database *gorm.DB) IRepository[T] {
	return &Repository[T]{
		database: database,
	}
}

// IBaseModel is interface which models should implement
type IBaseModel interface {
	TableName() string
}

// First finds the first record ordered by primary key, matching given conditions conds
func (r *Repository[T]) First(model *T, conds interface{}) error {
	res := r.database.First(&model, conds)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *Repository[T]) Find(models *[]T, conds interface{}) error {
	res := r.database.Find(&models, conds)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

// Create inserts value, returning the inserted data's primary key in value's id
func (r *Repository[T]) Create(model *T) (*T, error) {
	res := r.database.Create(model)

	if res.Error != nil {
		return nil, res.Error
	}

	return model, nil
}

// Create inserts value, returning the inserted data's primary key in value's id
func (r *Repository[T]) BulkCreate(models []*T) (int64, error) {
	res := r.database.Create(models)

	if res.Error != nil {
		return res.RowsAffected, res.Error
	}

	return res.RowsAffected, nil
}

// Save updates value in database. If value doesn't contain a matching primary key, value is inserted.
func (r *Repository[T]) Update(model *T) error {
	res := r.database.Save(model)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

// Delete deletes value matching given conditions.
// If value contains primary key it is included in the conditions.
// If value includes a deleted_at field, then Delete performs a soft delete
// instead by setting deleted_at with the current time if null.
func (r *Repository[T]) Delete(model *T) (int64, error) {
	res := r.database.Delete(model)

	if res.Error != nil {
		return res.RowsAffected, res.Error
	}

	return res.RowsAffected, nil
}

// GetDB return *gorm.DB for other methods which this repository doesn't support it
func (r *Repository[T]) GetDB() *gorm.DB {
	return r.database
}
