package system

import (
	"errors"
	"sort"
	"time"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
)

type MenuService struct {
	menuRepo *repository.MenuRepository
}

func NewMenuService() *MenuService {
	return &MenuService{
		menuRepo: repository.NewMenuRepository(),
	}
}

// GetMenuList 获取用户菜单树
func (s *MenuService) GetMenuList(userID int64) ([]*response.MenuResponse, error) {
	// 查询用户有权限的所有菜单
	menus, err := s.menuRepo.GetMenusByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	menuTree := s.buildMenuTree(menus, 0, userID)

	return menuTree, nil
}

// GetAllMenus 获取所有菜单（管理用）
func (s *MenuService) GetAllMenus() ([]*response.MenuResponse, error) {
	menus, err := s.menuRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 构建菜单树（不过滤用户权限）
	menuTree := s.buildMenuTreeAdmin(menus, 0)

	return menuTree, nil
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(req *request.CreateMenuRequest) error {
	// 如果有父菜单，检查父菜单是否存在
	if req.ParentID > 0 {
		_, err := s.menuRepo.FindByID(req.ParentID)
		if err != nil {
			return errors.New("父菜单不存在")
		}
	}

	menu := &model.Menu{
		ParentID:      req.ParentID,
		Name:          req.Label, // 前端label -> 后端Name（路由名称）
		Title:         req.Name,  // 前端name -> 后端Title（显示名称）
		Path:          req.Path,
		Component:     req.Component,
		Redirect:      req.Redirect,
		Icon:          req.Icon,
		IsEnable:      req.IsEnable,
		Sort:          req.Sort,
		IsMenu:        req.IsMenu,
		KeepAlive:     req.KeepAlive,
		IsHide:        req.IsHide,
		IsHideTab:     req.IsHideTab,
		Link:          req.Link,
		IsIframe:      req.IsIframe,
		ShowBadge:     req.ShowBadge,
		ShowTextBadge: req.ShowTextBadge,
		FixedTab:      req.FixedTab,
		ActivePath:    req.ActivePath,
		IsFullPage:    req.IsFullPage,
		CreateTime:    time.Now(),
	}

	return s.menuRepo.Create(menu)
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(req *request.UpdateMenuRequest) error {
	// 查询菜单是否存在
	menu, err := s.menuRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 不能将父菜单设置为自己或自己的子菜单
	if req.ParentID == req.ID {
		return errors.New("不能将父菜单设置为自己")
	}

	// 如果有父菜单，检查父菜单是否存在
	if req.ParentID > 0 {
		_, err := s.menuRepo.FindByID(req.ParentID)
		if err != nil {
			return errors.New("父菜单不存在")
		}
	}

	// 更新菜单信息
	menu.ParentID = req.ParentID
	menu.Name = req.Label // 前端label -> 后端Name（路由名称）
	menu.Title = req.Name // 前端name -> 后端Title（显示名称）
	menu.Path = req.Path
	menu.Component = req.Component
	menu.Redirect = req.Redirect
	menu.Icon = req.Icon
	menu.IsEnable = req.IsEnable
	menu.Sort = req.Sort
	menu.IsMenu = req.IsMenu
	menu.KeepAlive = req.KeepAlive
	menu.IsHide = req.IsHide
	menu.IsHideTab = req.IsHideTab
	menu.Link = req.Link
	menu.IsIframe = req.IsIframe
	menu.ShowBadge = req.ShowBadge
	menu.ShowTextBadge = req.ShowTextBadge
	menu.FixedTab = req.FixedTab
	menu.ActivePath = req.ActivePath
	menu.IsFullPage = req.IsFullPage

	return s.menuRepo.Update(menu)
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id int64) error {
	// 查询菜单是否存在
	_, err := s.menuRepo.FindByID(id)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 检查是否有子菜单
	hasChildren, err := s.menuRepo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		return errors.New("请先删除子菜单")
	}

	return s.menuRepo.Delete(id)
}

// buildMenuTree 递归构建菜单树（用户权限过滤）
func (s *MenuService) buildMenuTree(menus []model.Menu, parentID int64, userID int64) []*response.MenuResponse {
	tree := make([]*response.MenuResponse, 0)

	for _, menu := range menus {
		if menu.ParentID == parentID {
			// 查询菜单关联的按钮权限
			buttons, _ := s.menuRepo.GetButtonsByMenuID(menu.ID, userID)
			authList := make([]response.AuthButton, 0, len(buttons))
			for _, btn := range buttons {
				authList = append(authList, response.AuthButton{
					Title:    btn.AuthName,  // Button模型的AuthName映射到响应的Title
					AuthMark: btn.AuthLabel, // Button模型的AuthLabel映射到响应的AuthMark
				})
			}

			// 构建菜单节点
			node := &response.MenuResponse{
				ID:        menu.ID,
				ParentID:  menu.ParentID, // 添加 ParentID
				Name:      menu.Name,
				Path:      menu.Path,
				Component: menu.Component,
				Redirect:  menu.Redirect,
				Meta: response.MenuMeta{
					Title:      menu.Title,
					Icon:       menu.Icon,
					IsHide:     menu.IsHide,
					IsHideTab:  menu.IsHideTab,
					KeepAlive:  menu.KeepAlive,
					FixedTab:   menu.FixedTab,
					IsFullPage: menu.IsFullPage,
					AuthList:   authList,
				},
			}

			// 递归查找子菜单
			children := s.buildMenuTree(menus, menu.ID, userID)
			if len(children) > 0 {
				node.Children = children
			}

			tree = append(tree, node)
		}
	}

	// 按 sort 排序
	sort.Slice(tree, func(i, j int) bool {
		var iSort, jSort int
		for _, m := range menus {
			if m.ID == tree[i].ID {
				iSort = m.Sort
			}
			if m.ID == tree[j].ID {
				jSort = m.Sort
			}
		}
		return iSort < jSort
	})

	return tree
}

// buildMenuTreeAdmin 递归构建菜单树（管理用，不过滤权限）
func (s *MenuService) buildMenuTreeAdmin(menus []model.Menu, parentID int64) []*response.MenuResponse {
	tree := make([]*response.MenuResponse, 0)

	for _, menu := range menus {
		if menu.ParentID == parentID {
			// 构建菜单节点
			node := &response.MenuResponse{
				ID:        menu.ID,
				ParentID:  menu.ParentID, // 添加 ParentID
				Name:      menu.Name,
				Path:      menu.Path,
				Component: menu.Component,
				Redirect:  menu.Redirect,
				Meta: response.MenuMeta{
					Title:      menu.Title,
					Icon:       menu.Icon,
					IsHide:     menu.IsHide,
					IsHideTab:  menu.IsHideTab,
					KeepAlive:  menu.KeepAlive,
					FixedTab:   menu.FixedTab,
					IsFullPage: menu.IsFullPage,
				},
			}

			// 递归查找子菜单
			children := s.buildMenuTreeAdmin(menus, menu.ID)
			if len(children) > 0 {
				node.Children = children
			}

			tree = append(tree, node)
		}
	}

	// 按 sort 排序
	sort.Slice(tree, func(i, j int) bool {
		var iSort, jSort int
		for _, m := range menus {
			if m.ID == tree[i].ID {
				iSort = m.Sort
			}
			if m.ID == tree[j].ID {
				jSort = m.Sort
			}
		}
		return iSort < jSort
	})

	return tree
}
