package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const dateLayout = "2006-01-02T15:04:05"

var logger *log.Logger

func init() {
	f, err := os.Create("logs")
	if err != nil {
		panic(err)
	}

	logger = log.New(f, "", log.Flags())
}

type editingFinishedMsg struct {
	err      error
	filepath string
}

type dateChangedMsg struct {
	selectedDate time.Time
}

func changeDates(newDate time.Time) func() tea.Msg {
	return func() tea.Msg {
		return dateChangedMsg{selectedDate: newDate}
	}
}

var mdRenderer *glamour.TermRenderer

//go:embed input.tpl
var templateData string

func main() {
	now := time.Now()

	initialModel := &Model{
		today:        now,
		datesOfMonth: initMonth(),
		currentMonth: now.Month(),
		currentYear:  now.Year(),

		selectedWeek: 0,
		selectedDay:  int(now.Weekday()),
		selectedDate: time.Now(),

		repo:  newTaskrepo(),
		tasks: list.NewModel([]list.Item{taskItem{}}, list.NewDefaultDelegate(), 30, 25),
	}

	mdRenderer, _ = glamour.NewTermRenderer(
		// detect background color and pick either the default dark or light theme
		glamour.WithAutoStyle(),
		// wrap output at specific width
		glamour.WithWordWrap(40),
	)

	err := tea.NewProgram(initialModel, tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Model struct {
	today time.Time

	datesOfMonth [][7]time.Time
	currentMonth time.Month
	currentYear  int

	selectedDay  int
	selectedWeek int
	selectedDate time.Time

	repo         *taskrepo
	tasks        list.Model
	selectedTask taskItem

	detailMode bool
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) toggleDetailMode() {
	m.detailMode = !m.detailMode
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.detailMode = false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyCtrlN:
			m.selectedWeek++
			if m.selectedWeek > len(m.datesOfMonth) {
				m.selectedWeek = m.selectedWeek % len(m.datesOfMonth)
			}
			return m, changeDates(m.datesOfMonth[m.selectedWeek][m.selectedDay])
		case tea.KeyCtrlP:
			if m.selectedWeek > 0 && m.selectedWeek <= 5 {
				m.selectedWeek--
			}
			return m, changeDates(m.datesOfMonth[m.selectedWeek][m.selectedDay])
		case tea.KeyCtrlF:
			m.selectedDay++
			if m.selectedDay > 6 {
				if m.selectedWeek < len(m.datesOfMonth) {
					m.selectedWeek++
				} else {
					m.selectedWeek = 0
				}

				m.selectedDay = int(time.Sunday)
			}
			return m, changeDates(m.datesOfMonth[m.selectedWeek][m.selectedDay])
		case tea.KeyCtrlB:
			m.selectedDay--
			if m.selectedDay < 0 {
				if m.selectedWeek > len(m.datesOfMonth) {
					m.selectedWeek--
				} else {
					m.selectedWeek = 0
				}

				m.selectedDay = int(time.Saturday)
			}
			return m, changeDates(m.datesOfMonth[m.selectedWeek][m.selectedDay])
		case tea.KeyCtrlT:
			tmpFile, err := newTmpFileWithTemplate(template.Template{})
			if err != nil {
				panic(err)
			}

			cmd := exec.Command(os.Getenv("EDITOR"), tmpFile.Name())
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				tmpFile.Close()
				return editingFinishedMsg{err: err, filepath: tmpFile.Name()}
			})
		case tea.KeyEnter:
			idx := m.tasks.Index()
			item := m.tasks.Items()[idx]

			task := item.(taskItem)

			m.selectedTask = task
			m.toggleDetailMode()
		case tea.KeyEsc:
			m.toggleDetailMode()
		}
	case editingFinishedMsg:
		if msg.err != nil {
			panic(msg.err)
			//return m, tea.Quit
		}

		tmpFile, err := os.Open(msg.filepath)
		if err != nil {
			panic(err)
		}

		data, err := io.ReadAll(tmpFile)
		if err != nil {
			panic(err)
		}

		ti, err := parseInput(data)
		if err != nil {
			panic(err)
		}

		err = m.repo.addTask(ti)
		if err != nil {
			panic(err)
		}
	case dateChangedMsg:
		m.selectedDate = msg.selectedDate
	}

	m.tasks, cmd = m.tasks.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	rightView := ""

	if m.detailMode {
		rightView = m.renderTaskDetail()
	} else {
		rightView = m.renderTasks()
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, m.renderMonth(), rightView)
}

func (m *Model) renderTaskDetail() string {
	out, err := mdRenderer.Render(m.selectedTask.desc)
	if err != nil {
		panic(err)
	}

	return out
}

func (m *Model) renderMonth() string {
	out := ""
	weekRows := make([]string, 0)

	for wi, week := range m.datesOfMonth {
		weekDays := make([]string, 7)
		for di, day := range week {
			if day.Month() == m.currentMonth {
				weekDays[di] = lipgloss.JoinHorizontal(lipgloss.Right, out, activeDateStyle.Render(fmt.Sprintf(" %2d ", day.Day())))
			} else {
				weekDays[di] = lipgloss.JoinHorizontal(lipgloss.Right, out, inactiveDateStyle.Render(fmt.Sprintf(" %2d ", day.Day())))
			}

			if m.selectedDay == di && m.selectedWeek == wi {
				weekDays[di] = lipgloss.JoinHorizontal(lipgloss.Right, out, selectedDateStyle.Render(fmt.Sprintf(" %2d ", day.Day())))
			}
		}

		weekRows = append(weekRows, lipgloss.JoinHorizontal(lipgloss.Center, weekDays...))
	}

	out = boundaryStyle.Render(lipgloss.JoinVertical(lipgloss.Left, append([]string{weekdayHeaders}, weekRows...)...))

	return out
}

func (m *Model) getTaskList() ([]list.Item, error) {
	tasks, err := m.repo.getTasksOfDay(m.selectedDate)
	if err != nil {
		return nil, err
	}

	taskList := make([]list.Item, 0, len(tasks))
	for _, t := range tasks {
		taskList = append(taskList, list.Item(t))
	}

	return taskList, nil
}

func (m *Model) renderTasks() string {
	taskList, err := m.getTaskList()
	if err != nil {
		return err.Error()
	}

	m.tasks.SetItems(taskList)
	return m.tasks.View()
}

func newTmpFileWithTemplate(tpl template.Template) (*os.File, error) {
	t, err := template.ParseFiles("input.tpl")
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp("", "newtask_*.md")
	if err != nil {
		return nil, err
	}

	if err := t.Execute(tmpFile, struct{ Ts time.Time }{time.Now()}); err != nil {
		return nil, err
	}

	tmpFile.Sync()

	return tmpFile, nil
}

type taskInput struct {
	Timestamp time.Time
	Title     string
	Desc      string
}

func parseInput(data []byte) (taskInput, error) {
	headers := make(map[string]string)
	inFrontmatterBlock := false
	var desc strings.Builder

	r := bufio.NewReader(bytes.NewReader(data))
	for line, _, err := r.ReadLine(); err != io.EOF; line, _, err = r.ReadLine() {
		l := string(line)
		if strings.Trim(string(line), "\n") == "---" {
			if inFrontmatterBlock { // if already inside front matter block end it.
				inFrontmatterBlock = false
				continue
			} else {
				inFrontmatterBlock = true
				continue
			}
		}

		if !inFrontmatterBlock {
			desc.WriteString(l)
			continue
		}

		kv := strings.SplitN(l, ":", 2)
		headers[strings.TrimSpace(kv[0])] = kv[1]
	}

	ti := taskInput{Desc: desc.String()}
	if date, ok := headers["date"]; ok {
		t, err := time.Parse(dateLayout, trim(date))
		if err != nil {
			return ti, err
		}

		ti.Timestamp = t
	}

	if title, ok := headers["title"]; ok {
		ti.Title = trim(title)
	}

	logger.Printf("%+v\n", ti)
	return ti, nil
}

func trim(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"")
	s = strings.Trim(s, "'")

	return s
}
