package recursivelist

// import tea "github.com/charmbracelet/bubbletea"

type ItemDeleted struct {
	DeletedIndex int
}

type ItemAdded struct {
	InsertedAt int
}

type ExpandCollapse[T ItemWrapper[T]] struct {
	Affected ListItem[T]
	Collapsed bool
}