# kal

Kal is an interactive calendar with date wise note taking in CLI. The
program can be used as a todo or journaling. The date wise notes is
stored in a sqlite database. The program is still work in progress. I
add features and bugs as I use it. 

Dual-licensed under MIT or the [UNLICENSE](https://unlicense.org/).

## Installation

todo

## Usage

Kal currently shows the calendar of the current month and list of the
task associated for the selected date. A new note can be added to the
date by typing `Ctrl+t` which opens the default editor

The dates can be navigated using:

- `Ctrl+p` navigate to previous week
- `Ctrl+n` navigate to next week
- `Ctrl+f` move to next day
- `Ctrl+b` move to previoud day.
- `Enter` to view the markdown details of task

Key bindings are inspired from emacs.

## TODO

- Display description
- The calendar should display other months. 
- Add ability for user to include tags for the notes.
- Handle panics gracefully.
- easy way to export
