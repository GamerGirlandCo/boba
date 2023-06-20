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
	checked bool
}

func (r *rListItem) realAdd(in int, ra rListItem) {
	ra.parent = r
	var poit []rListItem = *r.children
	poit = slices.Insert(poit, utils.MinInt(len(poit) - 1, in), ra)
	r.children = &poit
	ra.parent = r
}

func (r rListItem) Add(in int, ra rListItem) {
	(&r).realAdd(in, ra)
}

func (r rListItem) AddMulti(in int, ra ...rListItem) {
  for li, val := range ra {
    r.Add(utils.MaxInt(li + in - 1, 0), val)
  }
}

func (r rListItem) Remove(tr int) rListItem {
	ret := (*r.children)[tr]
	*r.children = slices.Delete(*r.children, tr, tr + 1)
	return ret
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

func (r rListItem) GetChildren() []rListItem {
	// var c []recursivelist.ItemWrapper[rListItem]
	// for _, val := range *r.children {
	// 	c = append(c, recursivelist.NewItem[rListItem](val, rListDelegate{}).Value())
	// }
	// return c
	return *r.children
}
func (r rListItem) TotalBeneath() int {
	accum := len(*r.children)
	for _, val := range *r.children {
		accum += val.TotalBeneath()
	}
	return accum
}

func (r *rListItem) realSetParent(what *rListItem) {
	r.parent = what
} 
func (r rListItem) SetParent(what *rListItem) {
	(&r).realSetParent(what)
} 

func (r rListItem) SetChildren(ree []rListItem) {
	*&r.children = &ree
}

func (r rListItem) Value() *rListItem {
	return &r
}

func (r *rListItem) Toggle() {
	r.checked = !r.checked
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
	cbox := ""
	if i.Value().checked {
		cbox = " [✓]"
	} else {		
		cbox = " [ ]"
	}
	if i.ParentModel.Options.Expandable {
		if len(i.GetChildren()) > 0 {
			if i.Expanded() {
				s += " " + i.ParentModel.Options.OpenPrefix + " "
			} else {
				s += " " + i.ParentModel.Options.ClosedPrefix + " "
			}
		}
	}
	s += cbox + " " + i.Value().Name
	indento := i.Value().Lvl() * 3
	fn := styles.DefaultStyles.Text.Copy().
		Width(i.ParentModel.Options.Width).
		MarginLeft(indento).Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.DefaultStyles.Active.Copy().
				Width(i.ParentModel.Options.Width).
				MarginLeft(indento).
				Render(strings.Join(s, ""))
		}
	}
	fmt.Fprint(w, fn(s))
}
