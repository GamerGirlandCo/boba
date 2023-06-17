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
	children *[]rListItem
	parent *rListItem
}

func (r rListItem) FilterValue() string {
	return r.Name
}

func (r rListItem) Flatten() []rListItem {
	accum := make([]rListItem, 0)
	for _, ite := range *r.children {
		accum = append(accum, ite)
		accum = append(accum, ite.Flatten()...)
	}
 return accum
}

func (r rListItem) RModify(fnn func(rListItem)) {
	for _, val := range *r.children {
		// val.RModify(fnn)
		fnn(val)
	}
}

func (r rListItem) Find(a rListItem) int {
	for u := range *r.children {
		if (*r.children)[u].Name == a.Name {
			return u
		}
	}
	return -1
}

func (r rListItem) Children() []recursivelist.Indentable[rListItem] {
	ret := make([]recursivelist.Indentable[rListItem], 0)
	for _, i := range *r.children {
		ret = append(ret, i)
	}
	return ret
}

func (r rListItem) Parent() *rListItem {
	return r.parent
}

func (r rListItem) IndexWithinParent() int {
	if r.Parent() != nil {
		return r.Parent().Find(r)
	}
	return 0
}

func (r rListItem) Add(ra rListItem) {
	vee := &r
	vee.RealAdd(ra)
}

func (r *rListItem) RealAdd(ra rListItem) {
	*r.children = append(*r.children, ra)
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
	s += " " + i.Value.Name
	fn := styles.DefaultStyles.Text.Copy().Padding(0, 0, 0, i.Value.Lvl()).Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.DefaultStyles.Active.Copy().Padding(0, 0, 0, i.Value.Lvl()).Render(strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(s))
}


func genRandList(par *rListItem, model recursivelist.Model[rListItem], deleg list.ItemDelegate, maxDepth int, curDepth int) ([]recursivelist.ListItem[rListItem], []rListItem) {
	retVal := make([]recursivelist.ListItem[rListItem], 0)
	secRetVal := make([]rListItem, 0)
	for i := 0; i < rand.Intn(18) + 3; i++ {
		sts := []recursivelist.ListItem[rListItem]{}
		cv := model.NewItem(rListItem{
				Name: uuid.NewString(),
				parent: par,
				children: &[]rListItem{},
			}, deleg)
		// fmt.Printf("%d || %+v\n", i + 1, mo)
		if curDepth < maxDepth {
			curDepth++
			sts, *(cv.Value).children = genRandList(&cv.Value, model, deleg, maxDepth, curDepth)
			
			// fmt.Printf("%+v\n", sts)
			// fmt.Println("CVCVCVCVCCVCVCVCVCVCVCV")
			// fmt.Printf("%+v\n", cv.Component)
		}
		secRetVal = append(secRetVal, *cv.Value.children...)
		// cv.AddMulti(cv.Value.children...)
		bees := len(sts)
		for b, it := range sts {
			// fmt.Printf("dadadadadadadaada %+v\n\n", mo.Delegate)
			cv.Add(it.Value, b + bees)
		}
		retVal = append(retVal, cv)
		secRetVal = append(secRetVal, cv.Value)
	}

	return retVal, secRetVal
}

func initRlistModel() recursivelist.Model[rListItem] {
	m := recursivelist.New[rListItem]([]rListItem{}, rListDelegate{}, 500, 200)
	m.SetExpandable(true)
	
	rlisto, _ := genRandList(nil, m, rListDelegate{}, 3, 0)
	m.SetItems(rlisto)
	
	m.SetExpandable(true)

	return m
}