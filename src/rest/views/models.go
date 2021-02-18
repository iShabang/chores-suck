package views

import (
	"chores-suck/types"
)

type DashboardModel struct {
	User   *types.User
	Chores []types.ChoreListItem
}
