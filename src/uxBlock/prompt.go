package uxBlock

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zeropsio/zcli/src/i18n"
	"github.com/zeropsio/zcli/src/uxBlock/styles"
)

type promptConfig struct {
}

type PromptOption = func(cfg *promptConfig)

func (b *uxBlocks) Prompt(
	ctx context.Context,
	message string,
	choices []string,
	auxOptions ...PromptOption,
) (int, error) {
	cfg := promptConfig{}
	for _, opt := range auxOptions {
		opt(&cfg)
	}

	if !b.isTerminal {
		b.PrintInfo(styles.InfoLine(message))
		return 0, errors.New(i18n.T(i18n.PromptAllowedOnlyInTerminal))
	}

	model := &promptModel{
		cfg:      cfg,
		uxBlocks: b,
		message:  message,
		choices:  choices,
	}
	p := tea.NewProgram(model, tea.WithoutSignalHandler(), tea.WithContext(ctx))

	if _, err := p.Run(); err != nil {
		return 0, err
	}

	if model.canceled {
		b.ctxCancel()
		return 0, context.Canceled
	}

	return model.cursor, nil
}

type promptModel struct {
	cfg      promptConfig
	uxBlocks *uxBlocks
	message  string
	choices  []string
	cursor   int
	quiting  bool
	canceled bool
}

func (m *promptModel) Init() tea.Cmd {
	return nil
}

func (m *promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c":
			m.canceled = true
			return m, tea.Quit

		case "left":
			if m.cursor > 0 {
				m.cursor--
			}

		case "right":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.quiting = true

			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *promptModel) View() string {
	if m.quiting {
		return ""
	}

	buttonsTexts := []string{}
	for i, choice := range m.choices {
		if i == m.cursor {
			buttonsTexts = append(buttonsTexts, styles.ActiveDialogButton().Render(choice))
		} else {
			buttonsTexts = append(buttonsTexts, styles.DialogButton().Render(choice))
		}
	}

	question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(m.message)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, buttonsTexts...)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	dialog := lipgloss.Place(0, 0,
		lipgloss.Left, lipgloss.Center,
		styles.DialogBox().Render(ui),
	)

	return dialog
}
