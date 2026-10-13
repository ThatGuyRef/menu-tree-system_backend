package services

import (
	"errors"
	"menu-tree-system/backend/models"
	"menu-tree-system/backend/repositories"

	"github.com/google/uuid"
)

var (
	ErrMenuNotFound        = errors.New("menu item not found")
	ErrParentNotFound      = errors.New("parent menu item not found")
	ErrCircularDependency  = errors.New("cannot move menu under itself or one of its descendants")
	ErrInvalidSortPosition = errors.New("invalid reorder position")
)

type MenuService interface {
	GetMenuTree() ([]models.Menu, error)
	GetMenuItem(id uuid.UUID) (*models.Menu, error)
	CreateMenuItem(input models.CreateMenuInput) (*models.Menu, error)
	UpdateMenuItem(id uuid.UUID, input models.UpdateMenuInput) (*models.Menu, error)
	DeleteMenuItem(id uuid.UUID) error
	MoveMenuItem(id uuid.UUID, input models.MoveMenuInput) (*models.Menu, error)
	ReorderMenuItem(id uuid.UUID, input models.ReorderMenuInput) (*models.Menu, error)
}

type menuService struct {
	repo repositories.MenuRepository
}

func NewMenuService(repo repositories.MenuRepository) MenuService {
	return &menuService{repo: repo}
}

// GetMenuTree builds a hierarchical tree structure from flat list
func (s *menuService) GetMenuTree() ([]models.Menu, error) {
	menus, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Map untuk lookup cepat O(1)
	menuMap := make(map[uuid.UUID]*models.Menu)
	for i := range menus {
		menus[i].Children = []models.Menu{} // Pastikan empty slice, bukan nil
		menuMap[menus[i].ID] = &menus[i]
	}

	var rootMenus []models.Menu

	for i := range menus {
		item := &menus[i]
		if item.ParentID == nil {
			// Item root
			rootMenus = append(rootMenus, *item)
		} else {
			// Item anak, masukkan ke parent-nya
			if parent, exists := menuMap[*item.ParentID]; exists {
				parent.Children = append(parent.Children, *item)
			}
		}
	}

	// Rekonstruksi struktur tree dari rootMenus
	var buildTree func(items []models.Menu) []models.Menu
	buildTree = func(items []models.Menu) []models.Menu {
		var result []models.Menu
		for _, item := range items {
			if ptr, exists := menuMap[item.ID]; exists {
				if len(ptr.Children) > 0 {
					ptr.Children = buildTree(ptr.Children)
				}
				result = append(result, *ptr)
			}
		}
		return result
	}

	return buildTree(rootMenus), nil
}

func (s *menuService) GetMenuItem(id uuid.UUID) (*models.Menu, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrMenuNotFound
	}
	return menu, nil
}

func (s *menuService) CreateMenuItem(input models.CreateMenuInput) (*models.Menu, error) {
	depth := 0

	// Jika memiliki parent, validasi dan hitung depth
	if input.ParentID != nil {
		parent, err := s.repo.FindByID(*input.ParentID)
		if err != nil {
			return nil, ErrParentNotFound
		}
		depth = parent.Depth + 1
	}

	maxOrder, err := s.repo.GetMaxSortOrder(input.ParentID)
	if err != nil {
		return nil, err
	}

	menu := &models.Menu{
		Name:      input.Name,
		ParentID:  input.ParentID,
		Depth:     depth,
		SortOrder: maxOrder + 1,
	}

	if err := s.repo.Create(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

func (s *menuService) UpdateMenuItem(id uuid.UUID, input models.UpdateMenuInput) (*models.Menu, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrMenuNotFound
	}

	menu.Name = input.Name
	if err := s.repo.Update(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

func (s *menuService) DeleteMenuItem(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return ErrMenuNotFound
	}

	return s.repo.Delete(id)
}

func (s *menuService) MoveMenuItem(id uuid.UUID, input models.MoveMenuInput) (*models.Menu, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrMenuNotFound
	}

	// 1. Cek jika parent baru sama dengan parent saat ini
	if (menu.ParentID == nil && input.ParentID == nil) ||
		(menu.ParentID != nil && input.ParentID != nil && *menu.ParentID == *input.ParentID) {
		return menu, nil
	}

	// 2. Hindari self-parenting
	if input.ParentID != nil && *input.ParentID == id {
		return nil, ErrCircularDependency
	}

	targetDepth := 0
	if input.ParentID != nil {
		targetParent, err := s.repo.FindByID(*input.ParentID)
		if err != nil {
			return nil, ErrParentNotFound
		}

		// 3. Validasi circular dependency (target parent tidak boleh berada di bawah node ini)
		current := targetParent
		for current.ParentID != nil {
			if *current.ParentID == id {
				return nil, ErrCircularDependency
			}
			p, err := s.repo.FindByID(*current.ParentID)
			if err != nil {
				break
			}
			current = p
		}

		targetDepth = targetParent.Depth + 1
	}

	// 4. Update data menu
	maxOrder, err := s.repo.GetMaxSortOrder(input.ParentID)
	if err != nil {
		return nil, err
	}

	depthDiff := targetDepth - menu.Depth
	menu.ParentID = input.ParentID
	menu.Depth = targetDepth
	menu.SortOrder = maxOrder + 1

	if err := s.repo.Update(menu); err != nil {
		return nil, err
	}

	// 5. Update depth untuk semua anak secara rekursif jika depth berubah
	if depthDiff != 0 {
		s.updateDescendantsDepth(menu.ID, depthDiff)
	}

	return menu, nil
}

func (s *menuService) updateDescendantsDepth(parentID uuid.UUID, depthDiff int) {
	children, err := s.repo.FindChildren(&parentID)
	if err != nil {
		return
	}

	for _, child := range children {
		child.Depth += depthDiff
		_ = s.repo.Update(&child)
		s.updateDescendantsDepth(child.ID, depthDiff)
	}
}

func (s *menuService) ReorderMenuItem(id uuid.UUID, input models.ReorderMenuInput) (*models.Menu, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrMenuNotFound
	}

	siblings, err := s.repo.FindChildren(menu.ParentID)
	if err != nil {
		return nil, err
	}

	if input.NewPosition < 0 || input.NewPosition >= len(siblings) {
		return nil, ErrInvalidSortPosition
	}

	// Filter item yang sedang dipindahkan
	var reordered []models.Menu
	for _, item := range siblings {
		if item.ID != id {
			reordered = append(reordered, item)
		}
	}

	// Insert ke posisi baru
	reordered = append(reordered[:input.NewPosition], append([]models.Menu{*menu}, reordered[input.NewPosition:]...)...)

	// Update urutan semua siblings
	for i := range reordered {
		reordered[i].SortOrder = i
	}

	if err := s.repo.UpdateSortOrders(reordered); err != nil {
		return nil, err
	}

	menu.SortOrder = input.NewPosition
	return menu, nil
}