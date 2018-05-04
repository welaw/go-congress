package fdsys

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// billres or amendmentdoc
func DecodeLaw(f io.Reader) string {
	decoder := xml.NewDecoder(f)

	out := ""
	cols := 0
	var inElement string
	var tag string
	indent := 0
	stack := make([]xml.StartElement, 0)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			switch inElement {
			case "colspec":
				cols++
				//name := getAttr(se.Attr, "colname")
				//out += fmt.Sprintf("|%s", name)
			case "entry":
				out += "|"
			case "section":
				indent = 0
				out += "\n\n"
			case "subsection":
				indent++
				out += "\n"
				out = printIndent(indent, out)
				out += "- "
			case "title", "subtitle":
				out += "\n"
			case "paragraph", "subparagraph":
				indent++
				out += "\n"
				out = printIndent(indent, out)
				out += "- "
			case "enum", "header":
			case "external-xref":
			case "text":
				out += ""
				//tag = inElement
			case "header-in-text":
				//tag = inElement
				out += "__"
			case "quote":
				//tag = inElement
				out += "\""
			case "quoted-block":
				out += "\n"
			case "row":
				out += "\n"
				out = printIndent(indent, out)
			case "after-quoted-block", "clause", "quoted-block-continuation-text", "term",
				"continuation-text", "subclause", "short-title", "fraction", "chapter", "subchapter",
				"bold", "enum-in-header", "cosponsor", "part", "subpart", "pagebreak", "item",
				"toc-enum", "level-header", "target", "multi-column-toc-entry", "act-name":
			case "toc":
				indent = 0
				out += "\n"
			case "toc-entry":
				out += "\n"
			case "table":
				indent++
				cols = 0
				out += "\n"
			case "tbody":
				out += "\n"
				out = printIndent(indent, out)
				for i := 0; i < cols; i++ {
					out += "| --- "
				}
				out += "|"

			case "tgroup", "thead":
			case "toc-quoted-entry":
			default:
				fmt.Printf("unhandled start-tag: %v\n", inElement)
			}
			stack = append([]xml.StartElement{se}, stack...)

		case xml.EndElement:
			inElement = se.Name.Local
			switch inElement {
			case "entry", "colspec":
				//out += "|"
			case "section":
				indent = 0
			case "subsection":
				indent--
			case "title", "subtitle":
				out += "\n"
			case "paragraph", "subparagraph":
				indent--
			case "header-in-text":
				out += "__"
			case "quote":
				out += "\""
			case "quoted-block":
				out += "\n"
			case "text":
			case "toc-entry":
				//out += "\n"
			case "row":
				out += "|"
				//tag = section
			case "external-xref":
				//tag = "text"
			case "enum", "header", "after-quoted-block", "clause", "quoted-block-continuation-text",
				"term", "toc", "continuation-text", "subclause", "short-title", "fraction", "item",
				"bold", "enum-in-header", "cosponsor", "subpart", "part", "pagebreak", "chapter", "subchapter",
				"toc-enum", "level-header", "target", "multi-column-toc-entry", "act-name":
			case "table":
				indent--
				out += "\n"
			case "tgroup", "tbody", "thead":
			case "toc-quoted-entry":
			default:
				fmt.Printf("unhandled end-tag: %v\n", inElement)
			}
			stack = stack[1:]

		case xml.CharData:
			t := string(se)
			t = trimInnerWhitespace(t)

			if len(stack) == 0 {
				break
			}
			tag = stack[0].Name.Local
			switch tag {
			case "act-name":
				out += t
			case "enum":
				switch stack[1].Name.Local {
				case "section":
					// TODO  if top = SECTION else SEC.
					out += fmt.Sprintf("## SEC. %s", t)
				case "subsection":
					//out = printIndent(indent, out)
					out += fmt.Sprintf("_%s_ ", t)
				case "paragraph", "subparagraph":
					//out = printIndent(indent, out)
					out += fmt.Sprintf("_%s_ ", t)
				case "title":
					out += fmt.Sprintf("### TITLE %s", t)
				case "subtitle":
					out += fmt.Sprintf("#### Subtitle %s--", t)
				}
			case "header":
				switch stack[1].Name.Local {
				case "section":
					out += fmt.Sprintf(" %s\n", strings.ToUpper(t))
				case "subsection":
					out += fmt.Sprintf("%s.--", t)
				case "paragraph", "subparagraph":
					out += fmt.Sprintf("%s.--", t)
				case "title":
					out += fmt.Sprintf("--%s\n", strings.ToUpper(t))
				case "subtitle":
					out += fmt.Sprintf("%s\n", t)
				}
			case "text", "quoted-block", "paragraph", "quoted-block-continuation-text", "subparagraph", "term", "continuation-text":
				//out = printIndent(indent, out)
				out += t
			case "after-quoted-block":
				out += "\n" + t
			case "external-xref":
				cite := getAttr(stack[0].Attr, "parsable-cite")
				cites := strings.Split(cite, "/")
				if len(cites) < 3 {
					panic(fmt.Errorf("unknown external-xref: less than three parts: %v", cites))
				}
				//usc-chapter/5/45"
				if cites[0] == "usc" {
					out += fmt.Sprintf("[%s](http://uscode.house.gov/quicksearch/get.plx?title=%s&section=%s)", t, cites[1], cites[2])
				} else if cites[0] == "usc-chapter" {
					out += fmt.Sprintf("[%s](http://uscode.house.gov/view.xhtml?req=granuleid:USC-prelim-title%s-chapter%s-front&num=0&edition=prelim)", t, cites[1], cites[2])
				}
				//} else {
				//panic(fmt.Errorf("unknown external-xref: unk first part: %v", cites[0]))
				//}
			case "header-in-text", "quote":
				out += t
			case "enum-in-header":
				out += fmt.Sprintf(" __%s__ ", t)
			case "short-title":
				out += fmt.Sprintf("__%s__", t)
			case "toc-entry":
				lvl := getAttr(stack[0].Attr, "level")
				link := strings.Replace(t, " ", "-", -1)
				link = strings.ToLower(link)
				link = "#" + link
				switch lvl {
				case "section":
					out += fmt.Sprintf("* [%s](%s)", t, link)
				case "title":
					out += fmt.Sprintf("### [%s](%s)", strings.ToUpper(t), link)
				case "subtitle":
					out += fmt.Sprintf("#### [%s](%s)", t, link)
				case "chapter":
					out += fmt.Sprintf("### [%s](%s)", t, link)
				case "subchapter":
					out += fmt.Sprintf("#### [%s](%s)", t, link)
				case "part":
					out += fmt.Sprintf("#### [%s](%s)", t, link)
				case "subpart":
					out += fmt.Sprintf("##### [%s](%s)", t, link)
				case "division":
					out += fmt.Sprintf("##### [%s](%s)", t, link)
				case "item":
					out += fmt.Sprintf("##### [%s](%s)", t, link)
				default:
					panic(fmt.Errorf("toc-entry level not handled: %v", lvl))
				}
			case "toc-quoted-entry":
			case "entry":
				out += fmt.Sprintf("%s", t)
			case "bold":
				out += fmt.Sprintf("__%s__", t)
			case "cosponsor":
				out += fmt.Sprintf("%s", t)
			case "action-desc":
				out += fmt.Sprintf("%s", t)
			case "fraction":
				out += fmt.Sprintf("%s", t)
			case "level-header":
				// TODO
				out += t
			case "target":
				out += t
			case "toc-enum":
				out += t
			case "subsection", "subtitle", "section", "clause", "subclause", "toc", "title", "tbody",
				"chapter", "subchapter", "subpart", "part", "table", "thead", "tgroup", "item":
			default:
				fmt.Printf("not handled chardata: %v: '%v'\n", tag, t)
			}
		default:
			fmt.Printf("default: type=%+v\n", reflect.TypeOf(se))
		}
	}

	return out
}

func printIndent(i int, out string) string {
	for idx := 0; idx < i; idx++ {
		out += fmt.Sprintf("  ")
	}
	return out
}

func isInline(e xml.StartElement) bool {
	for _, a := range e.Attr {
		if a.Value == DisplayInline {
			return true
		}
		if a.Value == DontDisplayInline {
			return false
		}
	}
	return false
}

func isNotInline(e xml.StartElement) bool {
	for _, a := range e.Attr {
		if a.Value == DisplayInline {
			return false
		}
		if a.Value == DontDisplayInline {
			return true
		}
	}
	return false
}

func getAttr(attributes []xml.Attr, name string) string {
	for _, a := range attributes {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}
