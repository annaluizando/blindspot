package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
)

// common viewport functionality
type ViewportHelper struct {
	viewport   viewport.Model
	contentStr string
	width      int
	height     int
	showHelp   bool
}

func NewViewportHelper(width, height int) *ViewportHelper {
	return &ViewportHelper{
		width:  width,
		height: height,
	}
}

func (vh *ViewportHelper) UpdateDimensions(width, height int) {
	vh.width = width
	vh.height = height
}

func (vh *ViewportHelper) SetContent(content string) {
	vh.contentStr = content
	vh.updateViewport()
}

func (vh *ViewportHelper) SetHelpVisibility(showHelp bool) {
	vh.showHelp = showHelp
	vh.updateViewport()
}

func (vh *ViewportHelper) updateViewport() {
	helpHeight := ShortHelpHeight
	if vh.showHelp {
		helpHeight = FullHelpHeight
	}

	contentHeight := strings.Count(vh.contentStr, "\n") + 1
	viewportHeight := min(contentHeight, vh.height-helpHeight-1)

	vh.viewport = viewport.New(vh.width, viewportHeight)
	vh.viewport.SetContent(vh.contentStr)
}
