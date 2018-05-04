package fdsys

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	basename = "https://www.gpo.gov/fdsys/pkg/"
)

func ModsURL(url string) string {
	s := strings.Split(url, "/")
	j := strings.Join(s[:len(s)-1], "/")
	return fmt.Sprintf("%s/mods.xml", j)
}

func BillURL(b string) string {
	return fmt.Sprintf("%s%s/xml/%s.xml", basename, b, b)
}

// BillPublicURL is for law metadata.
func BillPublicURL(b string) string {
	return fmt.Sprintf("%s%s/html/%s.html", basename, b, b)
}

func BillModsURL(b string) string {
	return fmt.Sprintf("%s%s/xml/mods.xml", basename, b)
}

func ParseLastMod(lastMod string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, lastMod)
}

/*ParseLoc
	 c - Congress
	 t - Type
	 n - Number
	 v - Version
	 Example:
 	BILLS-<congress><type><number><version>
 	BILLS-114hrres12ih
*/

func ParseLoc(loc string) (c, t, n, v string) {
	s := `\bBILLS-\b(\d+)([a-zA-Z]+)(\d+)([a-zA-Z0-9]+)`
	r := regexp.MustCompile(s)
	match := r.FindAllStringSubmatch(loc, -1)
	if len(match) == 0 {
		return "", "", "", ""
	}
	if len(match[0]) < 5 {
		return "", "", "", ""
	}
	m := match[0]
	c = m[1]
	t = m[2]
	n = m[3]
	v = m[4]
	return
}

func ToLawLoc(c, t, n string) string {
	return fmt.Sprintf("BILLS-%s%s%s", c, t, n)
}

func ToIdent(c, t, n string) string {
	return fmt.Sprintf("%s%s%s", c, t, n)
}

func ToLongIdent(c, t, n, v string) string {
	return fmt.Sprintf("BILLS-%s%s%s%s", c, t, n, v)
}

func LocToIdent(loc string) string {
	c, t, n, v := ParseLoc(loc)
	return ToLongIdent(c, t, n, v)
}

func LocToShortIdent(loc string) string {
	c, t, n, _ := ParseLoc(loc)
	return ToIdent(c, t, n)
}
