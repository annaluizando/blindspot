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

func (nd *NotificationDisplay) RenderErrorWithPrefix(gs *game.GameState, prefix string) string {
	if gs.HasError() && gs.IsErrorRecent() {
		return "\n" + errorStyle.Render(prefix+gs.GetError())
	}
	return ""
}

func (nd *NotificationDisplay) RenderErrorInFooter(gs *game.GameState) string {
	return nd.RenderError(gs)
}

func (nd *NotificationDisplay) RenderSuccessMessageInFooter(gs *game.GameState) string {
	return nd.RenderSuccessMessage(gs)
}

func (nd *NotificationDisplay) RenderErrorAboveHelp(gs *game.GameState) string {
	if gs.HasError() && gs.IsErrorRecent() {
		return errorStyle.Render("⚠️  Error: "+gs.GetError()) + "\n"
	}
	return ""
}

func (nd *NotificationDisplay) RenderSuccessMessageAboveHelp(gs *game.GameState) string {
	if gs.HasSuccessMessage() && gs.IsSuccessMessageRecent() {
		return successStyle.Render("✅ "+gs.GetSuccessMessage()) + "\n"
	}
	return ""
}

func (nd *NotificationDisplay) RenderAllNotifications(gs *game.GameState) string {
	var result string
	result += nd.RenderError(gs)
	result += nd.RenderSuccessMessage(gs)
	return result
}
