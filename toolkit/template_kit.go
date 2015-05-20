package toolkit

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type Template struct {
	Template string
	Comps    map[string]string
}

func NewTemplate(template string, comps map[string]string) *Template {
	return &Template{template, comps}
}

func (this *Template) Do() string {
	result := this.Template

	temp := bytes.NewBuffer([]byte{})
	status := 0
	inRender := false

	for _, s := range this.Template {

		if !inRender {
			if s == '{' {
				status += 1
			} else {
				status = 0
			}

			if status == 2 {
				inRender = true
				temp.Reset()
				continue
			}
		}

		if inRender {
			temp.WriteRune(s)

			if s == '}' {
				status -= 1
			} else {
				status = 2
			}

			if status == 0 {
				inRender = false
				tmp := temp.String()
				fmt.Println(tmp)
				tmp = tmp[0 : len(tmp)-2]
				fmt.Println(this.doRender(tmp))
				result = strings.Replace(
					result,
					"{{"+tmp+"}}",
					this.doRender(tmp),
					-1,
				)
			}
		}

	}

	return result
}

func (this *Template) doRender(code string) string {

	if len(code) == 0 {
		return ""
	}

	comps := strings.Split(code, "|")
	result := this.Comps[comps[0]]

	t := reflect.Indirect(reflect.ValueOf(TemplateFunc{})).Type()
	v := reflect.New(t)

	for _, c := range comps[1:] {
		funcV := v.MethodByName(c)
		if funcV.IsValid() {
			reVals := funcV.Call(
				[]reflect.Value{
					reflect.ValueOf(result),
				},
			)
			if len(reVals) > 0 {
				result = reVals[0].String()
			}
		} else {
			fmt.Println("unknown template func", c)
		}
	}

	return result
}

type TemplateFunc struct {
}

func (this *TemplateFunc) Decode(code string) string {
	if len(code) == 32 {
		return strings.ToUpper(fmt.Sprintf("%s-%s-%s-%s-%s",
			code[0:8],
			code[8:12],
			code[12:16],
			code[16:20],
			code[20:32]))
	}
	return code

}

func (this *TemplateFunc) Lower(code string) string {
	return strings.ToLower(code)
}

func (this *TemplateFunc) Timestamp(code string) string {
	return "123"
}
