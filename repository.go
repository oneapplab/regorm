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
	First(model *T, conds ...interface{}) error        // Select query with limit 1
	FirstOrFail(model *T, conds ...interface{}) error  // Select query with limit 1 and return error if finds nothing
	Find(model *[]T, conds ...interface{}) error       // Select query
	FindOrFail(model *[]T, conds ...interface{}) error // Select query and return error if finds nothing
	Create(model *T) (*T, error)                       // Insert model
	BatchCreate(models []*T) (int64, error)            // Batch Insert based on slice of model
	Update(model *T) error                             // Update a model
	Delete(model *T) (int64, error)                    // Delete a record
	GetDB() *gorm.DB                                   // Get Database Instance
}

// Repository a generic struct which should be embed by other repositories
// to access repository methods
//
// sample usage as embed repository:
//
//	type SampleRepository struct {
//		Repository[SampleModel]
//	}
type Repository[T IBaseModel] struct {
	IRepository[T]

	Database *gorm.DB
}

// InitRepository use this in cases you don't want to embed Repository in your Repository structs
// sample usage as declare as repository:
//
// sampleRepository := InitRepository[SampleModel](db)
func InitRepository[T IBaseModel](database *gorm.DB) IRepository[T] {
	return &Repository[T]{
		Database: database,
	}
}

// IBaseModel is interface which models should implement
type IBaseModel interface {
	TableName() string
}

// First finds the first record ordered by primary key, matching given conditions
func (r *Repository[T]) First(model *T, conds ...interface{}) error {
	res := r.Database.First(&model, conds...)

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return res.Error
	}

	return nil
}

// FirstOrFail finds the first record ordered by primary key, matching given conditions
func (r *Repository[T]) FirstOrFail(model *T, conds ...interface{}) error {
	res := r.Database.First(&model, conds...)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

// Find finds the all the records ordered by primary key, matching given conditions
func (r *Repository[T]) Find(models *[]T, conds ...interface{}) error {
	res := r.Database.Find(&models, conds...)

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return res.Error
	}

	return nil
}

// FindOrFail finds the all the records ordered by primary key, matching given conditions
func (r *Repository[T]) FindOrFail(models *[]T, conds ...interface{}) error {
	res := r.Database.Find(&models, conds...)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

// Create inserts value, returning the inserted data's primary key in value's id
func (r *Repository[T]) Create(model *T) (*T, error) {
	res := r.Database.Create(model)

	if res.Error != nil {
		return nil, res.Error
	}

	return model, nil
}

// BulkCreate Create inserts value, returning the inserted data's primary key in value's id
func (r *Repository[T]) BulkCreate(models []*T) (int64, error) {
	res := r.Database.Create(models)

	if res.Error != nil {
		return res.RowsAffected, res.Error
	}

	return res.RowsAffected, nil
}

// Update Save updates value in database. If value doesn't contain a matching primary key, value is inserted.
func (r *Repository[T]) Update(model *T) error {
	res := r.Database.Save(model)

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
	res := r.Database.Delete(model)

	if res.Error != nil {
		return res.RowsAffected, res.Error
	}

	return res.RowsAffected, nil
}

// GetDB return *gorm.DB for other methods which this repository doesn't support it
func (r *Repository[T]) GetDB() *gorm.DB {
	return r.Database
}
