package programs

import (
	"io/ioutil"
	"strings"
	"fmt"
	"os"
	// "reflect"
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

// var program1 = Program{Language: "plain", Content:"# hola amigos", PageID: "1291", Name: "program1"}
// var program2 = Program{Language: "md", Content:"# todo bien ${program1}", PageID: "1918", Name: "program2"}
// var program3 = Program{Language: "md", Content:"# ${program1} Mundo", PageID: "1518", Name: "program3"}
var path = os.Getenv("GOPATH") + `\src\github.com\puedesleerlo\programs\`
func check(e error) {
    if e != nil {
        panic(e)
    }
}

//GetPrograms busca los programas en el disco y los devuelve en un json
func GetPrograms() interface{} {
	
	programs := make([]Program, 0)
	files, err := ioutil.ReadDir(path)
	check(err)
	for _, file := range files {
		if !file.IsDir() {
			program := Program{}
			dat, err := ioutil.ReadFile(path + file.Name())
			check(err)
			// fmt.Println(string(dat))
			program.Content = string(dat)
			parts := strings.Split(file.Name(), "-")
			if len(parts) != 2 {
				panic("Should have had two parts")
			}
			program.Name = parts[0]
			parts = strings.Split(parts[1], ".")
			if len(parts) != 2 {
				panic("Should have had two parts")
			}
			program.PageID = parts[0]
			program.Language = parts[1]
			programs = append(programs, program)
		}
	}
	fmt.Print(programs)
	return programs
}

func GetProgram(id string) Program {
	return Program{}
}

func GetEdited(id string, content string) Edited {
	program := GetProgram(id)
	return Edited{Program: program, Newcontent: content}
}

func SaveEdited(message map[string]interface{}) {
	//guardar el programa
	//La interface tiene que ser el id y el contenido
	fmt.Println("Se va a guardar el archivo")
	nameOfFile := path + messageProcessing(message["program"])
	fmt.Println(nameOfFile)
	newcontent :=  []byte(messageProcessing(message["newcontent"]))
	err := ioutil.WriteFile(nameOfFile, newcontent, 0644)
	check(err)
}

func messageProcessing(program interface{}) string {
	switch v := program.(type) {
	case map[string]interface{}:
		name:= fmt.Sprintf("%v-%v.%v", v["name"], v["pageID"], v["language"])
	   return name
	case string:
		return v
	 default:
	   // handle unknown type
		return ""
	 }
  }