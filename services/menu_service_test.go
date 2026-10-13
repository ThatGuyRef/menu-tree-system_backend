package services

import (
	"testing"

	"menu-tree-system/backend/models"

	"github.com/google/uuid"
)

// Mock repository for unit testing
type mockMenuRepository struct {
	menus []models.Menu
}

func (m *mockMenuRepository) FindAll() ([]models.Menu, error) {
	return m.menus, nil
}
func (m *mockMenuRepository) FindByID(id uuid.UUID) (*models.Menu, error) {
	for _, menu := range m.menus {
		if menu.ID == id {
			return &menu, nil
		}
	}
	return nil, ErrMenuNotFound
}
func (m *mockMenuRepository) FindChildren(parentID *uuid.UUID) ([]models.Menu, error) {
	var children []models.Menu
	for _, menu := range m.menus {
		if (parentID == nil && menu.ParentID == nil) ||
			(parentID != nil && menu.ParentID != nil && *menu.ParentID == *parentID) {
			children = append(children, menu)
		}
	}
	return children, nil
}
func (m *mockMenuRepository) Create(menu *models.Menu) error {
	m.menus = append(m.menus, *menu)
	return nil
}
func (m *mockMenuRepository) Update(menu *models.Menu) error {
	for i, existing := range m.menus {
		if existing.ID == menu.ID {
			m.menus[i] = *menu
			return nil
		}
	}
	return nil
}
func (m *mockMenuRepository) Delete(id uuid.UUID) error {
	var remaining []models.Menu
	for _, menu := range m.menus {
		if menu.ID != id {
			remaining = append(remaining, menu)
		}
	}
	m.menus = remaining
	return nil
}
func (m *mockMenuRepository) GetMaxSortOrder(parentID *uuid.UUID) (int, error) {
	max := -1
	for _, menu := range m.menus {
		if (parentID == nil && menu.ParentID == nil) ||
			(parentID != nil && menu.ParentID != nil && *menu.ParentID == *parentID) {
			if menu.SortOrder > max {
				max = menu.SortOrder
			}
		}
	}
	return max, nil
}
func (m *mockMenuRepository) UpdateSortOrders(menus []models.Menu) error {
	return nil
}

func TestGetMenuTree_Hierarchy(t *testing.T) {
	rootID := uuid.New()
	childID := uuid.New()
	grandchildID := uuid.New()

	flatMenus := []models.Menu{
		{ID: rootID, Name: "System Management", ParentID: nil, Depth: 0, SortOrder: 0},
		{ID: childID, Name: "Systems", ParentID: &rootID, Depth: 1, SortOrder: 0},
		{ID: grandchildID, Name: "System Code", ParentID: &childID, Depth: 2, SortOrder: 0},
	}

	repo := &mockMenuRepository{menus: flatMenus}
	service := NewMenuService(repo)

	tree, err := service.GetMenuTree()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(tree) != 1 {
		t.Fatalf("Expected 1 root item, got: %d", len(tree))
	}

	root := tree[0]
	if root.Name != "System Management" {
		t.Errorf("Expected root name 'System Management', got: %s", root.Name)
	}

	if len(root.Children) != 1 {
		t.Fatalf("Expected 1 child under root, got: %d", len(root.Children))
	}

	child := root.Children[0]
	if child.Name != "Systems" {
		t.Errorf("Expected child name 'Systems', got: %s", child.Name)
	}

	if len(child.Children) != 1 {
		t.Fatalf("Expected 1 grandchild, got: %d", len(child.Children))
	}

	grandchild := child.Children[0]
	if grandchild.Name != "System Code" {
		t.Errorf("Expected grandchild name 'System Code', got: %s", grandchild.Name)
	}
}

func TestCreateMenuItem_DepthCalculation(t *testing.T) {
	rootID := uuid.New()
	flatMenus := []models.Menu{
		{ID: rootID, Name: "Root", ParentID: nil, Depth: 0, SortOrder: 0},
	}

	repo := &mockMenuRepository{menus: flatMenus}
	service := NewMenuService(repo)

	// Create child under root
	childInput := models.CreateMenuInput{
		Name:     "Child Level 1",
		ParentID: &rootID,
	}

	child, err := service.CreateMenuItem(childInput)
	if err != nil {
		t.Fatalf("Expected no error creating child, got: %v", err)
	}

	if child.Depth != 1 {
		t.Errorf("Expected child depth to be 1, got: %d", child.Depth)
	}
}

func TestMoveMenuItem_CircularDependencyGuard(t *testing.T) {
	nodeA := uuid.New()
	nodeB := uuid.New()

	flatMenus := []models.Menu{
		{ID: nodeA, Name: "Node A", ParentID: nil, Depth: 0},
		{ID: nodeB, Name: "Node B", ParentID: &nodeA, Depth: 1},
	}

	repo := &mockMenuRepository{menus: flatMenus}
	service := NewMenuService(repo)

	// Attempt to move parent Node A under its own child Node B
	moveInput := models.MoveMenuInput{
		ParentID: &nodeB,
	}

	_, err := service.MoveMenuItem(nodeA, moveInput)
	if err != ErrCircularDependency {
		t.Errorf("Expected ErrCircularDependency, got: %v", err)
	}
}
