// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var oFlag = flag.String("o", "", "write output to `file` (default standard output)")

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

// {"alert":
//  {"name":"alert",
//   "keywords":["warning","triangle","exclamation","point"],
//   "heights":
//    {"16":
//     {"width":16,
//      "path":"<path fill-rule=\"evenodd\" d=\"M8.22...\"></path>"},
type octicon struct {
	Name     string
	Keywords []string
	Heights  map[string]octicon_path
}

type octicon_path struct {
	Width int
	Path  string
}

func run() error {
	f, err := os.Open(filepath.Join("_data", "data.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	var octicons map[string]octicon
	err = json.NewDecoder(f).Decode(&octicons)
	if err != nil {
		return err
	}

	var names []string
	for name := range octicons {
		names = append(names, name)
	}
	sort.Strings(names)

	var buf bytes.Buffer
	fmt.Fprint(&buf, `package octicon

import (
	"fmt"
)

type icon struct {
	Width   int
	Height  int
	SVGFmt string
}

// IconMap is a map of octigons as fmt.Printf format strings,
// sorted by their heights.
var IconMap = map[string][]icon{
`)
	// Write all individual Octicon functions.
	for _, name := range names {
		generateAndWriteOcticon(&buf, name, octicons[name])
	}

	fmt.Fprint(&buf, `}

// Icon returns the named Octicon SVG node.
// It returns nil if name is not a valid Octicon symbol name.
func Icon(name string, width int, height int) string {
	icons, found := IconMap[name]
	if !found {
		return ""
	}
	var icon icon
	for _, i := range(icons) {
		icon = i
		if icon.Height >= height {
			break
		}
	}
	return fmt.Sprintf(icon.SVGFmt, width, height)
}
`)

	var w io.Writer
	switch *oFlag {
	case "":
		w = os.Stdout
	default:
		f, err := os.Create(*oFlag)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	_, err = w.Write(buf.Bytes())
	return err
}

func generateAndWriteOcticon(w io.Writer, name string, icon octicon) error {
	fmt.Fprintln(w)
	fmt.Fprintf(w, "	\"%s\": []icon{\n", name)
	// Browse paths in height order, so that use-time icon height selection
	// will just have to walk the path list once 
	var heights []string
	for heightStr := range(icon.Heights) {
		heights = append(heights, heightStr)
	}
	sort.Strings(heights)
	for _, key := range(heights) {
		height, err := strconv.Atoi(key)
		if err != nil {
			return err
		}
		width := icon.Heights[key].Width
		path := icon.Heights[key].Path
		// fmt.Printf("generateOcticon: path is:\n%s\n", path)
		if strings.HasPrefix(path, `<path fill-rule="evenodd" `) {
			// Skip fill-rule, if present. It has no effect on displayed SVG, but takes up space.
			path = `<path ` + path[len(`<path fill-rule="evenodd" `):]
		}
		// Note, SetSize relies on the absolute position of the width, height attributes.
		// Keep them in sync with widthAttrIndex and heightAttrIndex.
		svgXML := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%%d" height="%%d" viewBox="0 0 %v %v">%s</svg>`,
			width, height, path)
		fmt.Printf("Dumping %s(%dx%d)...\n", name, width, height)
		fmt.Fprintf(w, "		icon{Width:%d, Height:%d, SVGFmt:`%s`,},\n", width, height, svgXML)
	}
	fmt.Fprintf(w, "	},\n")
	return nil
}
