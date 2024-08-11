package commands

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hajnalaron/git-convention-cli/types"
	"github.com/spf13/cobra"
	"strings"
)

var (
	branch = &cobra.Command{
		Use:   "branch",
		Short: "Create a conventional branch name",
		Run:   createBranch,
	}

	highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

type branchModel struct {
	branchTypes   list.Model
	description   textinput.Model
	issueNumber   textinput.Model
	selectedType  *types.Branch
	step          int
	generatedName string
	quitting      bool
	err           error
}

func createBranch(_ *cobra.Command, _ []string) {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	finalModel := m.(branchModel)
	if finalModel.quitting {
		return
	}

	fmt.Printf("\nGenerated branch name: %s\n", finalModel.generatedName)
	fmt.Println("To create this branch, run:")
	fmt.Printf("git switch -c %s\n", finalModel.generatedName)

	if err := clipboard.WriteAll(finalModel.generatedName); err == nil {
		fmt.Println("Branch name message copied to clipboard!")
	}
}

func initialModel() branchModel {
	branchItems := make([]list.Item, len(conf.BranchTypes))
	for i, bt := range conf.BranchTypes {
		branchItems[i] = bt
	}

	branchList := list.New(branchItems, list.NewDefaultDelegate(), 300, 0)
	branchList.Title = "Select Branch Type"

	for i, item := range branchItems {
		if item.(types.Branch).Type == conf.DefaultBranchPrefix {
			branchList.Select(i)
			break
		}
	}

	description := textinput.New()
	description.Placeholder = "Enter brief description (required)"
	description.CharLimit = 50

	issueNumber := textinput.New()
	issueNumber.Placeholder = "Enter issue number (optional)"
	issueNumber.CharLimit = 20

	return branchModel{
		branchTypes: branchList,
		description: description,
		issueNumber: issueNumber,
		step:        0,
	}
}

func (m branchModel) Init() tea.Cmd {
	return nil
}

func (m branchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			return m.handleEnter()
		}
	}

	var cmd tea.Cmd
	switch m.step {
	case 0:
		m.branchTypes, cmd = m.branchTypes.Update(msg)
	case 1:
		m.description, cmd = m.description.Update(msg)
	case 2:
		m.issueNumber, cmd = m.issueNumber.Update(msg)
	}
	return m, cmd
}

func (m branchModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case 0:
		if i, ok := m.branchTypes.SelectedItem().(types.Branch); ok {
			m.selectedType = &i
			m.step++
			m.description.Focus()
		}
	case 1:
		if m.description.Value() == "" {
			m.err = fmt.Errorf("description is required")
			return m, nil
		}
		m.err = nil
		m.step++
		m.issueNumber.Focus()
	case 2:
		m.generatedName = formatBranchName(*m.selectedType, m.description.Value(), m.issueNumber.Value())
		return m, tea.Quit
	}
	return m, nil
}

func (m branchModel) View() string {
	if m.quitting {
		return "Thanks for using Git Convention CLI!\n"
	}

	var s string
	switch m.step {
	case 0:
		s = m.branchTypes.View()
	case 1:
		s = fmt.Sprintf(
			"Branch Type: %s\n\n%s\n\n%s\n\n%s",
			highlightStyle.Render(m.selectedType.Type),
			"Enter a brief description (required):",
			m.description.View(),
			fmt.Sprintf("Characters: %d/%d", len(m.description.Value()), m.description.CharLimit),
		)
		if m.err != nil {
			s += "\n\n" + errorStyle.Render(m.err.Error())
		}
	case 2:
		s = fmt.Sprintf(
			"Branch Type: %s\nDescription: %s\n\n%s\n\n%s\n\n%s",
			highlightStyle.Render(m.selectedType.Type),
			m.description.Value(),
			"Enter issue number (optional):",
			m.issueNumber.View(),
			fmt.Sprintf("Characters: %d/%d", len(m.issueNumber.Value()), m.issueNumber.CharLimit),
		)
	}
	return fmt.Sprintf("%s\n\n%s", s, "(press q to quit)")
}

func formatBranchName(branchType types.Branch, description, issueNumber string) string {
	description = strings.ToLower(description)
	description = strings.ReplaceAll(description, " ", "-")

	if issueNumber != "" {
		return fmt.Sprintf("%s/%s-%s", branchType.Type, issueNumber, description)
	}
	return fmt.Sprintf("%s/%s", branchType.Type, description)
}
