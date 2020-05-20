package mygo

import (
	"fmt"
	"os"
	"testing"
	"text/template"
)
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func CustomFn(v string) string {
	return v+"22222"
}


func handleInt(number uint) uint {
	return number + 10
}
func handleString(field string) string {
	return " string is: " + field
}


func TestOld1(t *testing.T) {
	type Inventory struct {
		Material string
		Count    uint
	}
	sweaters := Inventory{"wool", 17}

	//tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	//tmpl, err := template.ParseFiles("demo.go.tpl")
	tmpl, err := template.New("T2").Parse("{{.Count}} items are made of")
	tmpl, err = tmpl.New("test").Parse("{{.Count}} items are made of {{.Material}}")

	CheckErr(err)
	//err = tmpl.Execute(os.Stdout, sweaters)
	err = tmpl.ExecuteTemplate(os.Stdout, "T2", sweaters) //可以选取模板
	CheckErr(err)
	fmt.Println("")
	fmt.Println(tmpl.Name())
	tmpl = tmpl.Lookup("test") //切换模板，必须要有返回，否则不生效
	fmt.Println(tmpl.Name())

}

func TestOld2(t *testing.T) {
	type Inventory struct {
		Material string
		Count    uint
	}
	sweaters := Inventory{"wool", 17}

	tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	CheckErr(err)
	file, err := os.OpenFile("demo.txt", os.O_CREATE|os.O_WRONLY, 0755)

	CheckErr(err)
	err = tmpl.Execute(file, sweaters)
	CheckErr(err)
}

func TestOld3(t *testing.T) {
		type Inventory struct {
			Material string
			Count    uint
		}
		type NewInventory struct {
			Fields []Inventory
		}
		sweaters := NewInventory{
			Fields: []Inventory{
				Inventory{Material: "wool", Count: 19},
				Inventory{Material: "wooltwo", Count: 20},
			}}

		var Text = `
	{{range .Fields }}
	  Material: {{.Material}} - Count:{{.Count}}
	{{ end }}
	`
		tmpl, err := template.New("test").Parse(Text)
		CheckErr(err)
		err = tmpl.Execute(os.Stdout, sweaters)
		CheckErr(err)

}

func TestOld4(t *testing.T) {
	type Inventory struct {
		Material string
		Count    uint
	}
	type NewInventory struct {
		Fields []Inventory
	}
	sweaters := NewInventory{
		Fields: []Inventory{
			Inventory{Material: "wool", Count: 19},
			Inventory{Material: "wooltwo", Count: 20},
		}}

	var Text = `
	{{range .Fields }}
	  Material: {{.Material | handleString}} - Count:{{.Count | handleInt }}
	{{ end }}
	`
	tmpl, err := template.New("test").Funcs(template.FuncMap{"handleString": handleString, "handleInt": handleInt}).Parse(Text)
	CheckErr(err)
	err = tmpl.Execute(os.Stdout, sweaters)
	CheckErr(err)

}

func TestOld5(t *testing.T) {
	type Inventory struct {
		Material string
		Count    uint
		Gift     *string
		Islog    bool
		Islog2   bool
		XX		string
	}

	gift := new(string)
	gift = nil
	sweaters := Inventory{"wool", 17, gift, true, true, "..."}

	tmpl, err := template.New("test").Funcs(template.FuncMap{"customFn":CustomFn}).Parse(
		`{{.Count}} items are made of {{.Material}} 
				{{customFn "xxx"}}
				{{- with customFn "xxx"}}
					{{ if eq . "xxx22222" }}
						fuckxxxxxxxxxxxxxxxx
					{{ end }}
				{{end}}
				{{ if eq "1" "1" }}
					{{.Material}}_{{.XX}}
				fuck
				{{ end }}
		`)
	//tmpl := template.Must(template.ParseFiles("demo.go.tpl", "demo2.go.tpl"))
	//tmpl, err := template.New("T2").Parse("{{.Count}} items are made of")
	//tmpl, err = tmpl.New("test").Parse("{{.Count}} items are made of {{.Material}}")

	//err = tmpl.Execute(os.Stdout, sweaters)
	err = tmpl.ExecuteTemplate(os.Stdout, "test", sweaters) //可以选取模板
	CheckErr(err)
	fmt.Println("")
	fmt.Println(tmpl.Name())
	tmpl = tmpl.Lookup("test") //切换模板，必须要有返回，否则不生效
	fmt.Println(tmpl.Name())
}