package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jmhobbs/little-machines/chip-8/chip8"
)

type model struct {
	run     bool
	step    bool
	file    string
	lastErr error
	machine chip8.Machine

	ctr uint64
}

func (m *model) Run() {
	t := time.NewTicker(2 * time.Millisecond) // 500hz
	defer t.Stop()
	for range t.C {
		m.ctr += 1
		if m.lastErr == nil && (m.run || m.step) {
			m.lastErr = m.machine.Step()
			m.step = false
		}
	}
}

type TickMsg time.Time

func main() {
	var (
		file string
		pgm  []byte
		err  error
	)

	if len(os.Args) >= 2 {
		pgm, err = os.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		file = filepath.Base(os.Args[1])
	}

	m, err := chip8.New(pgm)
	if err != nil {
		panic(err)
	}
	root := &model{
		file:    file,
		machine: m,
		run:     false,
	}

	go root.Run()

	p := tea.NewProgram(root, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Every(16*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case TickMsg:
		return m, tea.Every(16*time.Millisecond, func(t time.Time) tea.Msg {
			return TickMsg(t)
		})

	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "s":
			m.step = true

		case "r":
			m.run = true

		case "p":
			m.run = false
		}
	}

	return m, nil
}

func (m *model) View() string {
	leftColumn := []string{
		m.renderScreen(),
		m.renderRunbar(),
		"R - Run  S - Step  P - Pause",
	}

	if m.lastErr != nil {
		leftColumn = append(
			leftColumn,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#AA0000")).Render(m.lastErr.Error()),
		)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			leftColumn...,
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderState(),
			m.renderRegisters(),
		),
	)
}

var box BoxWithLabel = NewDefaultBoxWithLabel()

// Build the screen box, 32 high by 64 wide
func (m *model) renderScreen() string {
	s := m.machine.Screen()

	// Box value for when the pixel is on
	on := lipgloss.NewStyle().
		Background(lipgloss.Color("#FAFAFA")).
		Foreground(lipgloss.Color("#FAFAFA")).
		Width(1).
		Height(1).
		Render("X")

	var b strings.Builder

	// Iterate over the screen buffer and build the screen
	for y := 0; y < 32; y++ {
		for x := 0; x < 8; x++ {
			b.WriteString(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%08b", s[y*8+x]), "0", " "), "1", on))
		}
		if y < 31 {
			b.WriteRune('\n')
		}
	}

	return lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Render(b.String())
}

// Render the current machine state
func (m *model) renderState() string {
	state := m.machine.State()

	var b strings.Builder

	b.WriteString(fmt.Sprintf("PC: 0x%04x | %d\n", state.PC, state.PC))
	b.WriteString(fmt.Sprintf("SP:   0x%02x | %d\n", state.SP, state.SP))
	b.WriteString(fmt.Sprintf(" I: 0x%04x | %d\n", state.I, state.I))
	b.WriteString(fmt.Sprintf("DT:   0x%02x | %d\n", state.DT, state.DT))
	b.WriteString(fmt.Sprintf("ST:   0x%02x | %d", state.ST, state.ST))

	return box.Render("State", b.String(), 18)
}

// Render the current register values
func (m *model) renderRegisters() string {
	state := m.machine.State()

	var b strings.Builder

	for i, r := range state.V {
		if i < 15 {
			b.WriteString(fmt.Sprintf("%01x: 0x%02x\n", i, r))
		} else {
			b.WriteString(fmt.Sprintf("%01x: 0x%02x", i, r))
		}
	}

	return box.Render("Registers", b.String(), 12)
}

// Render the run state bar
func (m *model) renderRunbar() string {
	var b strings.Builder

	if m.lastErr != nil {
		b.WriteString("!")
	} else {
		if m.run {
			b.WriteString("▶")
		} else {
			b.WriteString("⏸")
		}
	}
	b.WriteString(" | ")
	if m.file != "" {
		b.WriteString(m.file)
	} else {
		b.WriteString("<none>")
	}

	return b.String()
}
