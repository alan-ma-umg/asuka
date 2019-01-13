package database

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type DataTestTable struct {
	Id        int                    `xorm:"pk autoincr"`
	Url       string                 `xorm:"varchar(1024)"`
	Data      map[string]interface{} `xorm:"json"`
	Version   int                    `xorm:"version"`
	UpdatedAt int                    `xorm:"updated"`
	CreatedAt int                    `xorm:"created"`
}

func TestMysql(t *testing.T) {
	//err := Mysql().DropTables(&DataWww{})
	//if err != nil {
	//	t.Error(err)
	//	panic(err)
	//}

	//Mysql().IsTableExist()
	err := Mysql().CreateTables(&DataTestTable{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	data := &DataTestTable{
		Url: "dsg98d7g9d7sfg987sdfg9ds7",
		Data: map[string]interface{}{
			"test":  "test",
			"test2": "test2",
			"test3": "test2",
			"test4": "test2",
			"test5": "test2",
			"test6": "test2",
		},
	}

	_, err = Mysql().Insert(data)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var j DataTestTable
	has, err := Mysql().Id(data.Id).Get(&j)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println(has)

	i, err := Mysql().Id(data.Id).Update(&DataTestTable{
		Url:     "122222222222222222",
		Version: data.Version,
	})
	fmt.Println(i)

}

type JsonProfile struct {
	Name string
	Age  int
}

type JsonInt int

type JsonField struct {
	Id       int64
	Name     map[string]string `xorm:"json"`
	Indexes  []int             `xorm:"json"`
	Indexes3 []JsonInt         `xorm:"json"`
	Profile  JsonProfile       `xorm:"json"`
	Profile2 *JsonProfile      `xorm:"json"`
	Name2    map[string]string
	Indexes2 []int
	//Profile3 JsonProfile
	//Profile4 *JsonProfile
}

func TestJsonField(t *testing.T) {
	//err := Mysql().DropTables(&JsonField{})
	//if err != nil {
	//	t.Error(err)
	//	panic(err)
	//}
	//
	//err = Mysql().CreateTables(&JsonField{})
	//if err != nil {
	//	t.Error(err)
	//	panic(err)
	//}

	js := &JsonField{
		Name: map[string]string{
			"test":  "test",
			"test2": "test2",
		},
		Indexes:  []int{1, 3, 5},
		Indexes3: []JsonInt{2, 4},
		Profile: JsonProfile{
			Name: "lll",
			Age:  12,
		},
		Profile2: &JsonProfile{
			Name: "lll",
			Age:  12,
		},
		Name2: map[string]string{
			"test":  "test",
			"test2": "test2",
		},
		Indexes2: []int{1, 3, 5},
		/*Profile3: JsonProfile{
			Name: "lll",
			Age:  12,
		},
		Profile4: &JsonProfile{
			Name: "lll",
			Age:  12,
		},*/
	}

	_, err := Mysql().Insert(js)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var j JsonField
	has, err := Mysql().Id(js.Id).Get(&j)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println("j:", j)

	if !has {
		err = errors.New("not exist")
		t.Error(err)
		panic(err)
	}

	if !reflect.DeepEqual(js, &j) {
		err = fmt.Errorf("%v is not equal %v", *js, j)
		t.Error(err)
		panic(err)
	}

	var jss = make([]JsonField, 0)
	err = Mysql().Find(&jss)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(jss) != 1 {
		err = errors.New("not exist")
		t.Error(err)
		panic(err)
	}

	if !reflect.DeepEqual(js, &jss[0]) {
		err = fmt.Errorf("%v is not equal %v", *js, j)
		t.Error(err)
		panic(err)
	}

	_, err = Mysql().Id(js.Id).Update(&JsonField{
		Profile: JsonProfile{
			Name: "---",
			Age:  100,
		},
	})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var newField JsonField
	has, err = Mysql().Id(js.Id).Get(&newField)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		err = errors.New("not exist")
		t.Error(err)
		panic(err)
	}
	newField.Profile = JsonProfile{
		Name: "lll",
		Age:  12,
	}

	if !reflect.DeepEqual(js, &newField) {
		err = fmt.Errorf("%v is not equal %v", *js, newField)
		t.Error(err)
		panic(err)
	}

	_, err = Mysql().Id(js.Id).Update(&JsonField{
		Profile2: &JsonProfile{
			Name: "---",
			Age:  100,
		},
	})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var newField2 JsonField
	has, err = Mysql().Id(js.Id).Get(&newField2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		err = errors.New("not exist")
		t.Error(err)
		panic(err)
	}
	newField2.Profile2 = &JsonProfile{
		Name: "lll",
		Age:  12,
	}
	newField2.Profile = JsonProfile{
		Name: "lll",
		Age:  12,
	}

	if !reflect.DeepEqual(js, &newField2) {
		err = fmt.Errorf("%v is not equal %v", *js, newField2)
		t.Error(err)
		panic(err)
	}

	_, err = Mysql().Id(js.Id).Update(&JsonField{
		Indexes3: []JsonInt{4, 8},
	})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var newField3 JsonField
	has, err = Mysql().Id(js.Id).Get(&newField3)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		err = errors.New("not exist")
		t.Error(err)
		panic(err)
	}

	newField3.Profile2 = &JsonProfile{
		Name: "lll",
		Age:  12,
	}
	newField3.Profile = JsonProfile{
		Name: "lll",
		Age:  12,
	}
	newField3.Indexes3 = []JsonInt{2, 4}

	if !reflect.DeepEqual(js, &newField3) {
		err = fmt.Errorf("%v is not equal %v", *js, newField3)
		t.Error(err)
		panic(err)
	}
}
