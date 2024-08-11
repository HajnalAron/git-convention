package commands

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hajnalaron/git-convention-cli/types"
	"github.com/spf13/cobra"
	"strings"
)

var (
	commit = &cobra.Command{
		Use:   "commit",
		Short: "Create a conventional commit message",
		Run:   createCommit,
	}
)

type commitModel struct {
	commitTypes  list.Model
	summary      textinput.Model
	body         textarea.Model
	selectedType *types.Commit
	step         int
	generatedMsg string
	quitting     bool
	err          error
	windowHeight int
	windowWidth  int
}

func createCommit(_ *cobra.Command, _ []string) {
	p := tea.NewProgram(initialCommitModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	finalModel := m.(commitModel)
	if finalModel.quitting {
		return
	}

	fmt.Printf("\nGenerated commit message:\n%s\n", finalModel.generatedMsg)
	fmt.Println("To create this commit, run:")
	fmt.Printf("git commit -m \"%s\"\n", strings.ReplaceAll(finalModel.generatedMsg, "\n", "\\n"))

	if err := clipboard.WriteAll(finalModel.generatedMsg); err == nil {
		fmt.Println("Commit message copied to clipboard!")
	}
}

func initialCommitModel() commitModel {
	commitItems := make([]list.Item, len(conf.CommitTypes))
	for i, ct := range conf.CommitTypes {
		commitItems[i] = ct
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("205"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("240"))
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("230"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(lipgloss.Color("245"))

	commitList := list.New(commitItems, delegate, 100, 30)
	commitList.SetHeight(20)
	commitList.SetShowTitle(false)
	commitList.SetFilteringEnabled(false)
	commitList.DisableQuitKeybindings()
	commitList.SetShowStatusBar(false)
	commitList.SetShowHelp(false)
	commitList.InfiniteScrolling = true

	summary := textinput.New()
	summary.Placeholder = "Enter commit summary (required)"
	summary.CharLimit = 72
	summary.Focus()

	body := textarea.New()
	body.Placeholder = "Enter commit body (optional)"
	body.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	body.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	body.CharLimit = 500
	body.Focus()

	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("205"))
	body.Cursor.Style = cursorStyle

	return commitModel{
		commitTypes: commitList,
		summary:     summary,
		body:        body,
		step:        0,
	}
}

func (m commitModel) Init() tea.Cmd {
	return nil
}

func (m commitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width

		listHeight := m.windowHeight - 6
		if listHeight < 3 {
			listHeight = 3
		}

		m.commitTypes.SetHeight(listHeight)
		m.commitTypes.SetWidth(m.windowWidth - 4)

		if m.step == 1 {
			m.summary.Width = m.windowWidth - 4
		} else if m.step == 2 {
			m.body.SetWidth(m.windowWidth - 4)
			m.body.SetHeight(m.windowHeight - 10)
		}

		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			return m.handleEnter()
		}
	}

	switch m.step {
	case 0:
		m.commitTypes, cmd = m.commitTypes.Update(msg)
	case 1:
		m.summary, cmd = m.summary.Update(msg)
	case 2:
		m.body, cmd = m.body.Update(msg)
	}
	return m, cmd
}

func (m commitModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case 0:
		if i, ok := m.commitTypes.SelectedItem().(types.Commit); ok {
			m.selectedType = &i
			m.step++
			m.summary.Focus()
		}
	case 1:
		if m.summary.Value() == "" {
			m.err = fmt.Errorf("summary is required")
			return m, nil
		}
		m.err = nil
		m.step++
		m.body.Focus()
	case 2:
		m.generatedMsg = formatCommitMessage(*m.selectedType, m.summary.Value(), m.body.Value())
		return m, tea.Quit
	}
	return m, nil
}

func (m commitModel) View() string {
	if m.quitting {
		return "Thanks for using Git Convention CLI!\n"
	}

	var s string
	switch m.step {
	case 0:
		s = m.commitTypes.View()
	case 1:
		s = fmt.Sprintf(
			"Commit Type: %s %s\n\n%s\n\n%s\n\n%s",
			highlightStyle.Render(m.selectedType.Type),
			m.selectedType.Emoji,
			"Enter a commit summary (required):",
			m.summary.View(),
			fmt.Sprintf("Characters: %d/%d", len(m.summary.Value()), m.summary.CharLimit),
		)
		if m.err != nil {
			s += "\n\n" + errorStyle.Render(m.err.Error())
		}
	case 2:
		s = fmt.Sprintf(
			"Commit Type: %s %s\nSummary: %s\n\n%s\n\n%s\n\n%s",
			highlightStyle.Render(m.selectedType.Type),
			m.selectedType.Emoji,
			m.summary.Value(),
			"Enter commit body (optional):",
			m.body.View(),
			fmt.Sprintf("Characters: %d/%d", len(m.body.Value()), m.body.CharLimit),
		)
	}
	return fmt.Sprintf("%s\n\n%s", s, "(press q to quit)")
}

func formatCommitMessage(commitType types.Commit, summary, body string) string {
	msg := fmt.Sprintf("%s%s: %s", commitType.Type, commitType.Emoji, summary)
	if body != "" {
		msg += "\n\n" + body
	}
	return msg
}
