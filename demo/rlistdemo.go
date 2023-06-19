package demo

import (
	"math/rand"

	// "reflect"

	"git.tablet.sh/tablet/boba/recursivelist"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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
	Check key.Binding
}

var dkm km = km {
	Check: key.NewBinding(
		key.WithKeys("."),
		key.WithHelp(".", "check/uncheck line"),
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, dkm.Check):
			selit := w.InnerValue.List.SelectedItem().(recursivelist.ListItem[rListItem])
			selit.Value().Toggle()
		}
	}
	something, cmd := w.InnerValue.Update(msg)
	cmds = append(cmds, cmd)
	*w.InnerValue = something.(recursivelist.Model[rListItem])
	return w, tea.Batch(cmds...)
}

func (w WrapperModel) View() string {
	return w.InnerValue.View()
}

var lor loremipsum.LoremIpsum = *loremipsum.NewWithSeed(1234)

func genRandList(par *rListItem, deleg list.ItemDelegate, maxDepth int, curDepth int, re recursivelist.Model[rListItem]) ([]recursivelist.ListItem[rListItem], []rListItem) {
	retVal := make([]recursivelist.ListItem[rListItem], 0)
	secRetVal := make([]rListItem, 0)
	for i := 0; i < rand.Intn(10)+1; i++ {
		sts := []recursivelist.ListItem[rListItem]{}
		fuckery := make([]rListItem, 0)
		cri := rListItem{
			Name:     lor.Word(),
			parent:   par,
			children: &[]rListItem{},
			checked: rand.Intn(2) == 1,
		}

		cv := recursivelist.NewItem(cri, deleg, MyOptions, re)
		if curDepth <= maxDepth {
			curDepth++
			sts, fuckery = genRandList(cv.Value(), deleg, maxDepth, curDepth, re)
		}
		cri.AddMulti(fuckery...)
			cv.SetChildren(sts)
		for ii, b := range sts {
			pari := cv.Value()
			(*cri.children)[ii] = *b.Value()
			for _, c := range *cri.children {
				c.parent = pari
			}
		}
		secRetVal = append(secRetVal, *cv.Value())
		// secRetVal = append(secRetVal, *cv.Value().children...)
		// cv.AddMulti(cv.Value.children...)
		retVal = append(retVal, cv)
	}

	return retVal, secRetVal
}

func initRlistModel() WrapperModel {
	MyOptions.SetExpandable(true)
	nu := recursivelist.New[rListItem]([]rListItem{}, rListDelegate{}, MyOptions)
	m := WrapperModel{
		InnerValue: &nu,
	}
	rlisto, _ := genRandList(nil, rListDelegate{}, 6, 0, nu)
	m.InnerValue.SetItems(rlisto)
	return m
}
