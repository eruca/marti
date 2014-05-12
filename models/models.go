package models

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/xorm"
)

type Func struct {
	Id       int64
	Lang     string `xorm:"index not null varchar(40) unique(lang) "`
	Sign     string `xorm:"index not null unique(lang)"`
	Name     string `xorm:"index"`
	Source   string
	Body     string    `xorm:"text"`
	Instance string    `xorm:"text"`
	Tags     string    `xorm:"index varchar(100)"`
	Created  time.Time `xorm:"index updated"`
}

type Tags struct {
	Id        int64
	Name      string    `xorm:"index varchar(20) not null unique"`
	Relatives int64     `xorm:"index"` //该标签相关数目
	Querys    int64     `xorm:"index"` //查询次数
	Desc      string    `xorm:"varchar(512)"`
	Created   time.Time `xorm:"updated"`
}

type TagRecord struct {
	Id       int64
	TagsName string    `xorm:"index varchar(20)"`
	Lang     string    `xorm:"index varchar(40)"`
	FuncSign string    `xorm:"index"`
	Created  time.Time `xorm:"created"`
}

var o *xorm.Engine

func InitDB() {
	var err error
	o, err = xorm.NewEngine("mysql", "root:nickwill@/Xor?charset=utf8")
	if trace(err) {
		os.Exit(2)
	}

	o.ShowErr = true
	o.ShowDebug = true

	o.Sync(&Func{}, &Tags{}, &TagRecord{})
}

func AddFunc(s []string) error {
	if len(s) != 7 {
		return errors.New("the params num is not 6")
	}

	//dischange the tags in search mode
	tagslice := strings.Fields(s[5])
	if (len(tagslice)) > 5 {
		return errors.New("the tags is more than 5!")
	}
	tag := "$" + strings.Join(tagslice, "#$") + "#"

	//write the query data for js

	//log.Println("new Func")
	nf := &Func{Lang: s[0], Sign: s[1], Name: s[6]}
	count, err := o.Count(nf)
	if trace(err) {
		return err
	}
	//log.Printf("count:%d, err:%v", count, err)

	//log.Printf("tagslice:%v,tag:%s", tagslice, tag)

	if count == 0 {
		nf.Lang = s[0]
		nf.Source = s[2]
		nf.Body = s[3]
		nf.Instance = s[4]
		nf.Tags = tag
	}

	//log.Println("insert Func")
	_, err = o.InsertOne(nf)
	if err != nil {
		trace(err)
	}

	for _, v := range tagslice {
		nt := Tags{}

		//如果找到两个或以上的key，那就需要作出选择。并能对该选择作出处理
		//比如 func、function，如果选择func，那就必须删除function,
		//并把function的key，加到func下面
		ok, err := o.Where("`name` = ?", v).Get(&nt)
		trace(err)

		if ok {
			nt.Relatives++
			//log.Println("update Tags")
			_, err = o.Id(nt.Id).Update(&Tags{Relatives: nt.Relatives})
		} else {
			//log.Println("insert Tags")
			_, err = o.InsertOne(&Tags{Name: v, Relatives: 1})
		}
		trace(err)

		//log.Println("insert tagRecord")
		_, err = o.InsertOne(&TagRecord{
			TagsName: v,
			Lang:     s[0],
			FuncSign: s[1],
		})
		trace(err)
	}

	return err
}

//if Get all if data,the 'id' must to be 0.
//modify now.get used to 'id' 'tag' 'name'
func GetFunc(i interface{}) (fs []*Func, err error) {
	//if the map 'i' is ni return all the functions
	if i == nil {
		err = o.Sql("SELECT * FROM `FUNC`").Find(&fs)
		if err != nil {
			return nil, err
		}
		return fs, nil
	}

	v := reflect.ValueOf(i)
	for !v.CanSet() {
		v = v.Elem()
	}

	key := v.MapKeys()[0]
	mv := v.MapIndex(key)

	switch key.String() {
	case "id":
		id := mv.Int()
		f := &Func{Id: id}
		_, err = o.Get(f)
		if err != nil {
			return nil, err
		}
		fs = append(fs, f)
	case "tag":
		err = o.Sql(`SELECT * FROM func WHERE sign IN 
			(SELECT tag_record.func_sign FROM tag_record WHERE tag_record.tags_name = ?)`, mv.String()).Find(&fs)
		if err != nil {
			return nil, err
		}
	case "name":
		name := mv.String()
		fn := &Func{Name: name}
		_, err = o.Get(fn)
		if err != nil {
			return nil, err
		}
		fs = append(fs, fn)
	default:
		return nil, errors.New("models.GetFunc-> the map is not in consideration!!!!")
	}
	return fs, nil
}

func SearchFunc(key string) (pfuncs []Func) {
	err := o.Sql(`SELECT NAME FROM FUNC WHERE name like '%` + key + `%'`).Asc("NAME").Find(&pfuncs)
	if err != nil {
		log.Println("models.SearchFunc->", err.Error())
		return nil
	}
	return pfuncs
}

// func SearchTag(key string)(tags []string){
// 	err := o.Sql(`SELECT `)
// }

//util functions............................
func Close() {
	o.Close()
	log.Println("engine is closing now!")
}

func IsExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func trace(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
