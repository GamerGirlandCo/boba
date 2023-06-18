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
	children []rListItem
	parent   *rListItem
	Options  recursivelist.Options
}

func (r rListItem) FilterValue() string {
	return r.Name
}

func (r rListItem) Value() *rListItem {
	return &r
}

func (r rListItem) Flatten() []rListItem {
	accum := make([]rListItem, 0)
	accum = append(accum, r)
	for _, ite := range r.children {
		accum = append(accum, ite.Flatten()...)
	}
	return accum
}

func (r rListItem) RModify(fnn func(rListItem)) {
	for _, val := range r.children {
		// val.RModify(fnn)
		fnn(val)
	}
}

func (r rListItem) Find(a rListItem) int {
	for u := range r.children {
		if (r.children)[u].Name == a.Name {
			return u
		}
	}
	return len(r.children) - 1
}

func (r rListItem) Children() []recursivelist.Indentable[rListItem] {
	ret := make([]recursivelist.Indentable[rListItem], 0)
	for _, i := range r.children {
		ret = append(ret, i)
	}
	return ret
}

func (r rListItem) Parent() *rListItem {
	return r.parent
}

func (r rListItem) TotalBeneath() int {
	accum := len(r.children)
	for _, val := range r.children {
		accum += val.TotalBeneath()
	}
	return accum
}

func (r rListItem) IndexWithinParent() int {
	if r.parent != nil {
		v := r.parent.Find(r)
		return v
	}

	return 0
}

func (r rListItem) Add(ra rListItem) {
	vee := &r
	vee.RealAdd(ra)
}

func (r *rListItem) RealAdd(ra rListItem) {
	ra.parent = r
	r.children = append(r.children, ra)
	ra.parent = r
}

func (r rListItem) AddMulti(ra ...rListItem) {
	for _, val := range ra {
		r.Add(val)
	}
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

func (r rListItem) ParentOptions() recursivelist.Options {
	return r.Options
}

func (r rListItem) SetOptions(o recursivelist.Options) {
	*&r.Options = o
	for _, rack := range r.children {
		rack.Options = o
		rack.SetOptions(o)
	}
}

type rListDelegate struct{}

func (d rListDelegate) Height() int                               { return 1 }
func (d rListDelegate) Spacing() int                              { return 1 }
func (d rListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d rListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i := recursivelist.NewItem[rListItem](listItem.(rListItem), d)

	// if !ok {
	// 	return
	// }
	s := ""
	if i.Value.ParentOptions().Expandable {
		if i.Expanded() {
			s += i.Value.ParentOptions().OpenPrefix + " "
		} else {
			s += i.Value.ParentOptions().ClosedPrefix + " "
		}
	}
	s += "" + i.Value.Name
	indento := i.Value.Lvl() * 2
	fn := styles.DefaultStyles.Text.Copy().Padding(0, 0, 0, indento).Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.DefaultStyles.Active.Copy().Padding(0, 0, 0, indento).Render(strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(s))
}