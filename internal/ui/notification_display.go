package ui

import (
	"blindspot/internal/game"
)

type NotificationDisplay struct{}

func NewNotificationDisplay() *NotificationDisplay {
	return &NotificationDisplay{}
}

func (nd *NotificationDisplay) RenderError(gs *game.GameState) string {
	if gs.HasError() && gs.IsErrorRecent() {
		return "\n" + errorStyle.Render("⚠️  Error: "+gs.GetError())
	}
	return ""
}

func (nd *NotificationDisplay) RenderSuccessMessage(gs *game.GameState) string {
	if gs.HasSuccessMessage() && gs.IsSuccessMessageRecent() {
		return "\n" + successStyle.Render("✅ "+gs.GetSuccessMessage())
	}
	return ""
}

func (nd *NotificationDisplay) RenderAllNotifications(gs *game.GameState) string {
	var result string
	result += nd.RenderError(gs)
	result += nd.RenderSuccessMessage(gs)
	return result
}
