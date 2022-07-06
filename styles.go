package main

import "github.com/charmbracelet/lipgloss"

const (
	// colors
	yellow = "228"
	purple = "63"
)

var commonStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	MarginBottom(1).
	MarginLeft(1).
	PaddingRight(1)

var boundaryStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(purple)).
	BorderTop(true).
	BorderLeft(true).
	BorderBottom(true).
	BorderRight(true).
	MarginTop(2).
	MarginRight(4).
	MarginBottom(2).
	MarginLeft(4)

var weekDayHeaderStyle = commonStyle.Copy().
	Background(lipgloss.Color("#0000FF"))

var inactiveDateStyle = commonStyle.Copy().
	Background(lipgloss.Color("#3C3C3C"))

var activeDateStyle = commonStyle.Copy().
	Background(lipgloss.Color("63"))

var selectedDateStyle = commonStyle.Copy().
	Background(lipgloss.Color(yellow))

var weekdayHeaders = lipgloss.JoinHorizontal(
	lipgloss.Center,
	weekDayHeaderStyle.Render(" Sun"),
	weekDayHeaderStyle.Render(" Mon"),
	weekDayHeaderStyle.Render(" Tue"),
	weekDayHeaderStyle.Render(" Wed"),
	weekDayHeaderStyle.Render(" Thu"),
	weekDayHeaderStyle.Render(" Fri"),
	weekDayHeaderStyle.Render(" Sat"),
)
