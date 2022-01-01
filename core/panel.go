package core

type panel interface {
	name() string
	entries(*UI)
	setEntries(*UI)
	updateEntries(*UI)
	setKeybinding(*UI)
	focus(*UI)
	unfocus()
	setFilterWord(string)
}
