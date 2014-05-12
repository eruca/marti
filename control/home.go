package control

import (
	//"html/template"
	//"log"
	"net/http"
	"strconv"
	"strings"

	"marti/models"
)

func HomeGet(w http.ResponseWriter, r *http.Request) {
	tmpl, err := parseFiles(nil, "./tpl/home.html")
	trace("home.HomeGet->", err)
	//get data from database,
	//this one is just a test.
	//d := map[string]interface{}{"name": "User1", "age": 18}

	//if get allFunc,the param must be 0
	//if err is not nil,it still can write out without data
	var fs []*models.Func

	t := r.FormValue("tag")
	if len(t) == 0 {
		fs, err = models.GetFunc(nil)
	} else {
		mtag := map[string]string{
			"tag": t,
		}
		fs, err = models.GetFunc(&mtag)
	}

	trace("home.HomeGet->", err)
	//write out
	trace("home.HomeGet->", tmpl.Execute(w, fs))
}

func AddFuncGet(w http.ResponseWriter, r *http.Request) {
	tmpl, err := parseFiles(nil, "./tpl/addFunc.html")
	trace("home.AddFuncGet->", err)

	d := map[string]interface{}{"name": "User1", "age": 18}
	//log.Println("in here")
	//write out
	trace("home.AddFuncGet->", tmpl.Execute(w, d))
}

func AddFuncPost(w http.ResponseWriter, r *http.Request) {
	trace("home.AddFuncPost->", r.ParseForm())

	keys := []string{"lang", "sign", "src", "body", "inst", "tags"}
	values := make([]string, 0, 7)
	var value string
	var fname string

	for _, v := range keys {
		value = r.Form.Get(v)
		if v == "sign" {
			//get the function name ex: AddFunc
			fnames := strings.Fields(value)
			for _, v := range fnames {
				index := strings.Index(v, "(")
				if index != -1 {
					fname = v[:index]
					break
				}
			}
		}
		values = append(values, value)
	}
	values = append(values, fname)

	err := models.AddFunc(values)
	trace("home.AddFuncPost->", err)

	r.Method = "Get"
	//r.Form = nil
	http.Redirect(w, r, "/", 302)
}

func View(w http.ResponseWriter, r *http.Request) {
	tmpl, err := parseFiles(nil, "./tpl/view.html")
	trace("home.View->", err)

	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if trace("home.View->", err) {
		http.Redirect(w, r, "/", 302)
	}

	mid := map[string]int64{
		"id": id,
	}
	fs, err := models.GetFunc(&mid)
	if trace("home.View->", err) {
		http.Redirect(w, r, "/", 302)
	}

	m := make(map[string]interface{})
	m["Func"] = fs[0]
	label := strings.Replace(strings.Replace(fs[0].Tags, "$", " ", -1), "#", "", -1)
	m["Labels"] = strings.Fields(label)

	trace("home.View->", tmpl.Execute(w, m))
}

//search autocomplete request
func Autocomplete(w http.ResponseWriter, r *http.Request) {
	p := r.FormValue("p")
	fs := models.SearchFunc(p)
	var res string

	for _, v := range fs {
		res += v.Name + "\n"
	}

	_, err := w.Write([]byte(res))
	trace("home.Search->", err)
}

//search the database not finished
func Search(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	var data = make(map[string]interface{})

	//search function with name
	f := map[string]string{"name": search}
	fs, err := models.GetFunc(&f)
	trace("home.Search->", err)
	if len(fs) != 0 {
		data["func"] = fs
	}

}
