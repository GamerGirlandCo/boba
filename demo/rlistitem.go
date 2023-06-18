package demo

import (
	"fmt"
	"io"
	"strings"

	"git.tablet.sh/tablet/boba/recursivelist"
	"git.tablet.sh/tablet/boba/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type rListItem struct {
	Name     string
	children *[]rListItem
	parent   *rListItem
	Options  recursivelist.Options
}

func (r rListItem) Find(a rListItem) int {
	for u := range *r.children {
		if (*r.children)[u].Name == a.Name {
			return u
		}
	}
	return len(*r.children) - 1
}

func (r rListItem) FilterValue() string {
	return r.Name
}

func (r rListItem) IndexWithinParent() int {
	if r.parent != nil {
		v := r.parent
		return (v.Find(r))
	}

	return 0
}


func (r rListItem) Lvl() int {
	base := 0
	par := r.parent
	for par != nil {
		base++
		par = par.parent
	}
	return base
}

func (r rListItem) GetParent() *rListItem {
	return r.parent
}

func (r rListItem) GetChildren() []recursivelist.ItemWrapper[rListItem] {
	var c []recursivelist.ItemWrapper[rListItem]
	for _, val := range *r.children {
		c = append(c, recursivelist.NewItem[rListItem](val, rListDelegate{}).Value())
	}
	return c
}
func (r rListItem) TotalBeneath() int {
  accum := len(*r.children)
  for _, val := range *r.children {
    accum += val.TotalBeneath()
  }
  return accum
}

func (r rListItem) SetChildren(ree []recursivelist.ItemWrapper[rListItem]) {
	choild := make([]rListItem, 0)
	for _, val := range ree {
		choild = append(choild, *val.Value())
	}
	reeo := &r
	*reeo.children = choild
}

func (r rListItem) Value() *rListItem {
	return &r
}

type rListDelegate struct{}

func (d rListDelegate) Height() int                               { return 1 }
func (d rListDelegate) Spacing() int                              { return 1 }
func (d rListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d rListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i := listItem.(recursivelist.ListItem[rListItem])

	// if !ok {
	// 	return
	// }
	s := ""
	if i.ParentModel.Options.Expandable {
		if i.Expanded() {
			s += i.ParentModel.Options.OpenPrefix + " "
		} else {
			s += i.ParentModel.Options.ClosedPrefix + " "
		}
	}
	s += "" + i.Value().Name
	indento := i.Value().Lvl() * 1
	fn := styles.DefaultStyles.Text.Copy().Padding(0, 0, 0, indento).Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.DefaultStyles.Active.Copy().Padding(0, 0, 0, indento).Render(strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(s))
}
