package repositories

import (
	"menu-tree-system/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MenuRepository defines the interface for menu data access
type MenuRepository interface {
	FindAll() ([]models.Menu, error)
	FindByID(id uuid.UUID) (*models.Menu, error)
	FindChildren(parentID *uuid.UUID) ([]models.Menu, error)
	Create(menu *models.Menu) error
	Update(menu *models.Menu) error
	Delete(id uuid.UUID) error
	GetMaxSortOrder(parentID *uuid.UUID) (int, error)
	UpdateSortOrders(menus []models.Menu) error
}

type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository creates a new MenuRepository instance
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

// FindAll returns all menus ordered by parent and sort order
func (r *menuRepository) FindAll() ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Preload("Parent").Order("sort_order ASC").Find(&menus).Error
	return menus, err
}

// FindByID returns a single menu by ID with its parent loaded
func (r *menuRepository) FindByID(id uuid.UUID) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.Preload("Parent").First(&menu, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// FindChildren returns direct children of a parent
func (r *menuRepository) FindChildren(parentID *uuid.UUID) ([]models.Menu, error) {
	var menus []models.Menu
	query := r.db.Order("sort_order ASC")

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Find(&menus).Error
	return menus, err
}

// Create inserts a new menu item
func (r *menuRepository) Create(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

// Update saves changes to an existing menu item
func (r *menuRepository) Update(menu *models.Menu) error {
	return r.db.Save(menu).Error
}

// Delete removes a menu item by ID (children cascade via DB constraint)
func (r *menuRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Menu{}, "id = ?", id).Error
}

// GetMaxSortOrder returns the highest sort_order among siblings
func (r *menuRepository) GetMaxSortOrder(parentID *uuid.UUID) (int, error) {
	var maxOrder int
	query := r.db.Model(&models.Menu{}).Select("COALESCE(MAX(sort_order), -1)")

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Scan(&maxOrder).Error
	return maxOrder, err
}

// UpdateSortOrders batch updates sort_order for multiple menus
func (r *menuRepository) UpdateSortOrders(menus []models.Menu) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, menu := range menus {
			if err := tx.Model(&models.Menu{}).
				Where("id = ?", menu.ID).
				Update("sort_order", menu.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}