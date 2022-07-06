package main

import "time"

type taskItem struct {
	title, desc string
	ts          time.Time
}

func (i taskItem) Title() string       { return i.ts.Format("15:04 MST") + " " + i.title }
func (i taskItem) Description() string { return i.desc }
func (i taskItem) FilterValue() string { return i.title }
