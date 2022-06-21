package main

//
//import (
//	"fmt"
//	"os"
//	"text/template"
//)
//
//func main() {
//	fmt.Println("Hello, playground")
//
//	const templ = `Here is what they said
//    {{ -f .}}
//    {{.}}
//    {{end}}
//    `
//	//x := map[string]string{
//	//	"Danny": "The guy really talked about my first time out with you",
//	//	"Doug":  "Well he said I'm really amazing, I did not believe at first",
//	//}
//	type x struct {
//		Name string `json:"name"`
//		Age  int    `json:"age"`
//	}
//
//	t, err := template.New("index.html").Parse(templ)
//	if err != nil {
//		fmt.Println("Could not parse template:", err)
//		return
//
//	}
//
//	t.Execute(os.Stdout, x{
//		Name: "hahah",
//		Age:  12,
//	})
//
//}
