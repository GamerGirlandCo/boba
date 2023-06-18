package demo

import (
	"math/rand"

	// "reflect"

	"git.tablet.sh/tablet/boba/recursivelist"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

var MyOptions recursivelist.Options = recursivelist.Options{
	ClosedPrefix: ">",
	OpenPrefix:   "⌵",
	Width:        600,
	Height:       250,
	Expandable:   false,
	Keymap:       recursivelist.DefaultKeys,
}

type WrapperModel struct {
	InnerValue recursivelist.Model[rListItem]
}

func (w WrapperModel) Init() tea.Cmd {
	return nil
}

func (w WrapperModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	something, cmd := w.InnerValue.Update(msg)
	cmds = append(cmds, cmd)
	w.InnerValue = something.(recursivelist.Model[rListItem])
	return w, tea.Batch(cmds...)
}

func (w WrapperModel) View() string {
	return w.InnerValue.View()
}

func genRandList(par *rListItem, model recursivelist.Model[rListItem], deleg list.ItemDelegate, maxDepth int, curDepth int) ([]recursivelist.ListItem[rListItem], []rListItem) {
	retVal := make([]recursivelist.ListItem[rListItem], 0)
	secRetVal := make([]rListItem, 0)
	for i := 0; i < rand.Intn(10); i++ {
		// sts := []rListItem{}
		fuckery := []recursivelist.ListItem[rListItem]{}

		var oooooo recursivelist.Options
		if par != nil {
			oooooo = par.Options
		} else {
			oooooo = MyOptions
		}

		cv := recursivelist.NewItem(rListItem{
			Name:     uuid.NewString(),
			parent:   par,
			children: &[]rListItem{},
			Options:  oooooo,
		}, deleg)
		// fmt.Printf("%d || %+v\n", i + 1, mo)
		if curDepth < maxDepth {
			curDepth++
			pari := cv.Value()
			fuckery, _ = genRandList(&pari, model, deleg, maxDepth, curDepth)
			cv.SetChildren(fuckery)
			// fmt.Printf("%+v\n", sts)
			// fmt.Println("CVCVCVCVCCVCVCVCVCVCVCV")
			// fmt.Printf("%+v\n", cv.Component)
		}
		// for ii, b := range sts {
		// 	cv.Add(b, ii)
		// }
		secRetVal = append(secRetVal, cv.Value())
		// secRetVal = append(secRetVal, *cv.Value().children...)
		// cv.AddMulti(cv.Value.children...)
		retVal = append(retVal, cv)
	}

	return retVal, secRetVal
}

func initRlistModel() WrapperModel {
	m := WrapperModel{
		InnerValue: recursivelist.New[rListItem]([]rListItem{}, rListDelegate{}, 500, 200, MyOptions),
	}
	MyOptions.SetExpandable(true)
	rlisto, _ := genRandList(nil, m.InnerValue, rListDelegate{}, 6, 0)
	m.InnerValue.SetItems(rlisto)
	return m
}
