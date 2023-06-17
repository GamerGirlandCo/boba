package demo

import (
	"io"
	"math/rand"
	// "reflect"

	"fmt"
	"strings"

	"git.tablet.sh/tablet/boba/recursivelist"
	"git.tablet.sh/tablet/boba/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type rListItem struct {
	Name string
	children []rListItem
	Parent *rListItem
}

func (r rListItem) FilterValue() string {
	return r.Name
}

func (r rListItem) Lvl() int {
	base := 0
	par := r.Parent
	for par != nil {
		base++
		par = par.Parent
	}
	return base
}

type rListDelegate struct{}

func (d rListDelegate) Height() int                               { return 1 }
func (d rListDelegate) Spacing() int                              { return 1 }
func (d rListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d rListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(recursivelist.ListItem[rListItem])
	if !ok {
		return
	}
	s := ""
	if i.ListOptions().Expandable {
		if i.Expanded() {
			s += i.ListOptions().OpenPrefix + " "
		} else {
			s += i.ListOptions().ClosedPrefix + " "
		}
	}
	s += " " + i.Value.Name
	fn := styles.DefaultStyles.Text.Copy().Padding(0, 0, 0, i.Value.Lvl() * 2).Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.DefaultStyles.Active.Copy().Padding(0, 0, 0, i.Value.Lvl() * 2).Render(strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(s))
}


func genRandList(par *rListItem, deleg list.ItemDelegate, maxDepth *int, curDepth int) ([]recursivelist.ListItem[rListItem], []rListItem) {
	retVal := make([]recursivelist.ListItem[rListItem], 0)
	secRetVal := make([]rListItem, 0)
	for i := 0; i < rand.Intn(18) + 3; i++ {
		sts := []recursivelist.ListItem[rListItem]{}
		cv := recursivelist.NewItem[rListItem](rListItem{
				Name: uuid.NewString(),
				Parent: par,
			}, deleg)
		// fmt.Printf("%d || %+v\n", i + 1, mo)
		if curDepth < *maxDepth {
			curDepth++
			sts, (cv.Value).children = genRandList(&cv.Value, deleg, maxDepth, curDepth)
			// fmt.Printf("%+v\n", sts)
			// fmt.Println("CVCVCVCVCCVCVCVCVCVCVCV")
			// fmt.Printf("%+v\n", cv.Component)
		}
		// cv.AddMulti(cv.Value.children...)
		bees := len(sts)
		for b, it := range sts {
			// fmt.Printf("dadadadadadadaada %+v\n\n", mo.Delegate)
			cv.AddItem(it.Value, b + bees)
		}
		retVal = append(retVal, cv)
		secRetVal = append(secRetVal, cv.Value)
	}

	return retVal
}

func initRlistModel() recursivelist.Model[rListItem] {
	m := recursivelist.New[rListItem]([]rListItem{}, rListDelegate{}, 500, 200)
	m.SetExpandable(true)
	rlisto, _ := genRandList(nil, rListDelegate{}, &bla, 0)
	m.SetItems(rlisto)
	
	m.SetExpandable(false)

	return m
}