package demo

import (

	"math/rand"

	// "reflect"

	"git.tablet.sh/tablet/boba/recursivelist"
	"github.com/charmbracelet/bubbles/list"
	"github.com/google/uuid"
)

var MyOptions recursivelist.Options = recursivelist.Options{
	ClosedPrefix: ">",
	OpenPrefix:   "⌵",
	Width:        600,
	Height:       250,
}

func genRandList(par *rListItem, model recursivelist.Model[rListItem], deleg list.ItemDelegate, maxDepth int, curDepth int) ([]recursivelist.ListItem[rListItem], []rListItem) {
	retVal := make([]recursivelist.ListItem[rListItem], 0)
	secRetVal := make([]rListItem, 0)
	for i := 0; i < rand.Intn(8) + 3; i++ {
		sts := []recursivelist.ListItem[rListItem]{}


		var oooooo recursivelist.Options
		if par != nil {
			oooooo = par.Options
		} else {
			oooooo = MyOptions
		}

		cv := recursivelist.NewItem(rListItem{
				Name: uuid.NewString(),
				parent: par,
				children: &[]rListItem{},
				Options: oooooo,
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
	MyOptions.SetExpandable(true)
	rlisto, _ := genRandList(nil, m, rListDelegate{}, 5, 0)
	m.SetItems(rlisto)
	return m
}