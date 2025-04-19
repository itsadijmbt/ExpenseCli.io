package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type Expense struct {
	Name   string
	Amount float64
	Date   string
}


type expenseItem struct{ exp Expense }

func (e expenseItem) Title() string { return e.exp.Name }
func (e expenseItem) Description() string {
	return fmt.Sprintf("%.2f – %s", e.exp.Amount, e.exp.Date)
}
func (e expenseItem) FilterValue() string {
	return strings.Join([]string{
		e.exp.Name,
		strconv.FormatFloat(e.exp.Amount, 'f', -1, 64),
		e.exp.Date,
	}, " ")
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)


type model struct {
	expenseList list.Model      // main list component
	prompting   bool            // true when adding a new expense
	animating   bool            // true while progress bar is animating
	textInput   textinput.Model // text input bubble
	progress    progress.Model  // progress bar bubble
}


func loadInitialData(path string) ([]list.Item, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var items []list.Item
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}
		amt, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}
		items = append(items, expenseItem{Expense{parts[0], amt, parts[2]}})
	}
	return items, scanner.Err()
}

func inputValidator(s string) bool {
	parts := strings.Fields(s)
	if len(parts) != 3 {
		return false
	}

	if _, err := strconv.ParseFloat(parts[1], 64); err != nil {
		return false
	}

	dateParts := strings.Split(parts[2], "-")
	if len(dateParts) != 3 {
		return false
	}
	day, err1 := strconv.Atoi(dateParts[0])
	mon, err2 := strconv.Atoi(dateParts[1])
	yr, err3 := strconv.Atoi(dateParts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}
	if mon < 1 || mon > 12 || day < 1 || day > 31 || yr < 1900 {
		return false
	}
	return true
}


func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// cmd up front so that, if we delegate to sub‑models later, we have a slot to capture their returned command.
	// ── PROMPT MODE ───────────────────────────────────────
	if m.prompting {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case msg.Type == tea.KeyEnter:

				line := m.textInput.Value()
				if !inputValidator(line) {

					m.prompting = false
					return m, nil
				}
				parts := strings.Fields(line)
				amt, _ := strconv.ParseFloat(parts[1], 64)
				newExp := Expense{parts[0], amt, parts[2]}

				items := m.expenseList.Items()
				items = append(items, expenseItem{newExp})
				m.expenseList.SetItems(items)

				if f, err := os.OpenFile("expenses.txt",
					os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644); err == nil {
					fmt.Fprintf(f, "%s %f %s\n",
						newExp.Name, newExp.Amount, newExp.Date)
					f.Close()
				}

				m.prompting = false
				m.animating = true
				return m, m.progress.SetPercent(1.0)

			case msg.String() == "q":
				
				m.prompting = false
				return m, nil
			}

		}

		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}


	if m.animating {
		generic, cmd := m.progress.Update(msg)

		m.progress = generic.(progress.Model)
		if m.progress.IsAnimating() {
			return m, cmd
		}
		m.animating = false
		return m, nil
	}


	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+n":
	
			m.prompting = true
			m.textInput.SetValue("")
			
			return m, func() tea.Msg { return textinput.Blink() }
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.expenseList.SetSize(msg.Width-h, msg.Height-v)
	}

	m.expenseList, cmd = m.expenseList.Update(msg)
	return m, cmd
}


func (m model) View() string {
	if m.prompting {
		return docStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				"Enter new expense (Name Amount Date):",
				m.textInput.View(),
			),
		)
	}
	if m.animating {

		bar := lipgloss.NewStyle().Align(lipgloss.Center).Render(m.progress.View())
		return docStyle.Render(bar)
	}

	return docStyle.Render(m.expenseList.View() +
		"\n\nPress ctrl+n to add an expense, q or ctrl+c to quit.")
}

func main() {
	
	items, err := loadInitialData("expenses.txt")
	if err != nil {
		fmt.Println("Error loading data:", err)
		os.Exit(1)
	}


	expList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	expList.Title = "Expenses"

	
	ti := textinput.New()
	ti.Placeholder = "e.g., Coffee 2.50 12-04-2023"
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 40


	pb := progress.New(progress.WithDefaultGradient())


	m := model{
		expenseList: expList,
		textInput:   ti,
		progress:    pb,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
