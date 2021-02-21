package views

import (
	"chores-suck/core/types"
)

type DashboardModel struct {
	User   *types.User
	Chores []types.ChoreListItem
}
