package demo

import (
	"math/rand"

	"git.tablet.sh/tablet/boba/recursivelist"
	"git.tablet.sh/tablet/boba/utils"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-loremipsum/loremipsum"
)

var MyOptions recursivelist.Options = recursivelist.Options{
	ClosedPrefix: ">",
	OpenPrefix:   "⌵",
	Width:        600,
	Height:       250,
	Expandable:   true,
	Keymap:       recursivelist.DefaultKeys,
}

type km struct {
	Check  key.Binding
	Delete key.Binding
	Insert key.Binding
}

var dkm km = km{
	Check: key.NewBinding(
		key.WithKeys("."),
		key.WithHelp(".", "check/uncheck line"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d", tea.KeyDelete.String()),
		key.WithHelp(
			"d/del",
			"delete item",
		),
	),
	Insert: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "insert item"),
	),
}

type WrapperModel struct {
	InnerValue *recursivelist.Model[rListItem]
}

func (w WrapperModel) Init() tea.Cmd {
	return w.InnerValue.Init()
}

func (w WrapperModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	selit := w.InnerValue.List.SelectedItem().(recursivelist.ListItem[rListItem])
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, dkm.Check):
			selit.Value().Toggle()
		case key.Matches(msg, dkm.Insert):
			itit := rListItem{
				Name:     lor.Word(),
				parent:   selit.Value(),
				children: &[]rListItem{},
				checked:  false,
			}

			selit.Add(itit, len(selit.GetChildren())-1)
		}
	}
	something, cmd := w.InnerValue.Update(msg)
	cmds = append(cmds, cmd)
	*w.InnerValue = something.(recursivelist.Model[rListItem])

	return w, tea.Sequence(cmds...)
}

func (w WrapperModel) View() string {
	return w.InnerValue.View()
}

var lor *loremipsum.LoremIpsum = loremipsum.New()

func genRandList(par *rListItem, maxDepth int, curDepth int, re recursivelist.Model[rListItem]) []rListItem {

	retVal := make([]rListItem, 0)
	for i := 0; i < rand.Intn(8)+4; i++ {
		sts := []rListItem{}
		cri := rListItem{
			Name:     (*lor).Word(),
			parent:   par,
			children: &[]rListItem{},
			checked:  rand.Intn(2) == 1,
		}

		cv := re.NewItem(cri)
		if curDepth < maxDepth {
			curDepth++

			sts = genRandList(&cri, maxDepth, curDepth, re)
		}
		//cri.AddMulti(i, sts...)
		cv.AddMulti(i, sts...)

		retVal = append(retVal, cri)
	}

	return retVal
}

func initRlistModel() WrapperModel {
	MyOptions.SetExpandable(true)
	nu := recursivelist.New[rListItem]([]rListItem{}, rListDelegate{}, MyOptions)
	m := WrapperModel{}
	kiksi := genRandList(nil, 6, 0, nu)
	nu.AddToRoot(kiksi...)
	m.InnerValue = &nu
	var iv []key.Binding = m.InnerValue.List.AdditionalShortHelpKeys()
	iv = append(iv, utils.IterKeybindings(dkm)...)
	m.InnerValue.List.AdditionalShortHelpKeys = func() []key.Binding {
		return iv
	}
	return m
}
