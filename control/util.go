package control

import (
	"bufio"
	"html/template"
	"io"
	"log"
	"marti/models"
	"os"
	"path/filepath"
	"strings"
)

const (
	queryfile = "./public/js/query.js"
)

func trace(des string, err error) bool {
	if err != nil {
		log.Println(des, err.Error())
		return true
	}
	return false
}

func parseFiles(t *template.Template, fileName string) (*template.Template, error) {
	var err error
	if t == nil {
		t, err = template.ParseFiles(fileName)
		if trace("util->parseFiles", err) {
			return nil, err
		}
	} else {
		t, err = t.ParseFiles(fileName)
		if trace("util->parseFiles", err) {
			return nil, err
		}
	}

	file, err := os.Open(fileName)
	if trace("util->parseFiles", err) {
		return nil, err
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	var line string
	for {
		line, err = buf.ReadString('\n')

		left := strings.Index(line, "{{")
		right := strings.LastIndex(line, "}}")
		if left != -1 && right != -1 && left < right &&
			strings.Contains(line[left+2:right], "template") {
			//log.Println("find template: ", line[left:right+2])
			s := line[left:right]
			left = strings.Index(s, `"`)
			right = strings.LastIndex(s, `"`)

			//log.Println("[template] ", s[left+1:right])

			newfile := filepath.Join("./tpl/", s[left+1:right])
			if !models.IsExist(newfile) {
				log.Println("the file is not exist", newfile)
				continue
			}
			t, err = parseFiles(t, newfile)

			line = ""
		}

		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return t, err
		}
	}
	return t, err
}

//可以使用goroutine，不过需要考虑同步问题，添加waitgroup
//暂时还是使用查询数据库吧
func addquery(query string) error {
	var bfirst = false
	if !models.IsExist(query) {
		os.MkdirAll(filepath.Dir(queryfile), os.ModePerm)
		os.Create(queryfile)
		bfirst = true
	}
	//第一次的时候需要添加
	f, err := os.OpenFile(queryfile, os.O_WRONLY, 0777)
	trace("util.addquery->", err)
	defer f.Close()
	if bfirst {
		f.WriteString("var query = [];")
	}

	fi, err := os.Stat(queryfile)
	//貌似下面一句好像没有必要
	f.Truncate(fi.Size() - 2)

	query = ",\"" + query + "\"];"
	_, err = f.WriteAt([]byte(query), fi.Size()-2)
	trace("util.addquery->", err)
	return err
}
