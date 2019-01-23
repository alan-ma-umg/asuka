package main

import (
	"asuka/database"
	"asuka/helper"
	"asuka/project"
	"asuka/web"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(os.Args) < 2 {
		log.Fatal("Example:/path/to/asuka /path/to/env.json")
	}
	helper.PathToEnvFile = os.Args[1]
}

func main() {
	mainStart := time.Now()
	helper.ExitHandle()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()

	//type D struct {
	//	A string
	//	B []string
	//	C int
	//}
	//
	//data := &D{
	//	A: "Title",
	//	B: []string{"a", "b"},
	//	C: 100,
	//}
	//
	//payload := &D{}

	p := project.New(&project.Test{})
	p.Run()

	//for _,s:=range p.GetSpiders()

	encBuf := &bytes.Buffer{}
	gob.Register(p.Project)
	gob.Register(&http.Transport{})
	enc := gob.NewEncoder(encBuf)
	err := enc.Encode(p)
	if err != nil {
		fmt.Println(err)
	}
	database.Redis().Set("gob", encBuf.String(), time.Hour)
	//fmt.Println(encBuf.Len())
	encGob, _ := database.Redis().Get("gob").Result()

	decBuf := &bytes.Buffer{}
	decBuf.WriteString(encGob)

	dec := gob.NewDecoder(decBuf)

	pp := project.New(&project.Test{})
	err = dec.Decode(pp)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pp)

	//fmt.Println(payload)

	//reflect.ValueOf(data).Elem().FieldByName("C").SetInt(7)

	//fmt.Println(data)

	//var x float64 = 3.4
	//p := reflect.ValueOf(&x) // Note: take the address of x.
	//fmt.Println("settability of p:", p.Elem().CanSet())
	//p.Elem().SetFloat(11.2)
	//fmt.Println(p.Elem().Interface())
	//fmt.Println(x)

	//type T struct {
	//A int
	//B string
	//}
	//t := T{23, "skidoo"}

	//p := project.New(&project.Test{})
	//
	//s := reflect.ValueOf(p).Elem()
	//typeOfT := s.Type()
	//for i := 0; i < s.NumField(); i++ {
	//	f := s.Field(i)
	//	if f.CanSet() {
	//		fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
	//
	//		if f.Kind() != reflect.Slice {
	//			ss := reflect.ValueOf(f.Interface()).Elem()
	//			typeOfTT := ss.Type()
	//			for ii := 0; ii < ss.NumField(); ii++ {
	//				f := ss.Field(i)
	//				if f.CanSet() {
	//					fmt.Printf("iinner %d: %s %s = %v\n", ii, typeOfTT.Field(i).Name, f.Type(), f.Interface())
	//				}
	//			}
	//		}
	//	}
	//}

	//asuka()
}

func asuka() {
	p := project.New(&project.Test{})
	p.Run()

	//cleanUp(p) //todo !!!!!!!!!

	z := project.New(&project.ZhiHu{})
	z.Run()
	fmt.Println("Monitor: http://127.0.0.1:666")
	projects := []*project.Dispatcher{p, z}

	web.Server(projects, ":666") // http://127.0.0.1:666
}

func cleanUp(p *project.Dispatcher) {
	for i := 0; i < 10; i++ {
		os.Remove(helper.Env().BloomFilterPath + p.GetProjectName() + "_enqueue_retry_" + strconv.Itoa(i) + ".db")
	}
	//database.Mysql().Exec("truncate asuka_dou_ban")
	database.Bl().ClearAll()
	database.Redis().Del(p.GetQueueKey())
}
