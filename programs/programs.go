package programs

import (
	"io/ioutil"
	"fmt"
	"os"
)
type Program struct {
    Language    string `json:"language"`
    Name        string `json:"name"`
    PageID      string `json:"pageID"`
    Content     string `json:"content"`
}

type Edited struct {
	Program Program `json:"program"`
	Newcontent string `json:"newcontent"`
}

var program1 = Program{Language: "plain", Content:"# hola amigos", PageID: "1291", Name: "program1"}
var program2 = Program{Language: "md", Content:"# todo bien ${program1}", PageID: "1918", Name: "program2"}
var program3 = Program{Language: "md", Content:"# ${program1} Mundo", PageID: "1518", Name: "program3"}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

//GetPrograms busca los programas en el disco y los devuelve en un json
func GetPrograms() []Program {
	path := os.Getenv("GOPATH") + `\src\github.com\puedesleerlo\programs\`
	files, err := ioutil.ReadDir(path)
	check(err)
	for _, file := range files {
		fmt.Printf("%v", file.IsDir())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(path + file.Name())
			check(err)
			fmt.Print(string(dat))
		}
    }
	return []Program{program1, program2, program3}
}

func GetProgram(id string) Program {
	return program1
}

func GetEdited(id string, content string) Edited {
	program := GetProgram(id)
	return Edited{Program: program, Newcontent: content}
}

func SaveEdited(message interface{}) {
	//guardar el programa
}