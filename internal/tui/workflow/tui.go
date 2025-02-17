package workflow

import (
	"github.com/cloudposse/atmos/pkg/schema"

	tea "github.com/charmbracelet/bubbletea"
	mouseZone "github.com/lrstanley/bubblezone"
)

// Execute starts the TUI app and returns the selected items from the views
func Execute(workflows map[string]schema.WorkflowConfig) (*App, error) {
	mouseZone.NewGlobal()
	mouseZone.SetEnabled(true)

	app := NewApp(workflows)
	p := tea.NewProgram(app, tea.WithMouseCellMotion())

	_, err := p.Run()
	if err != nil {
		return nil, err
	}

	return app, nil
}
