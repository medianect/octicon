package octicon_test

import (
	"fmt"
	"log"

	"github.com/medianect/octicon"
)

func Example() {
	str := octicon.Icon("alert", 28, 28)
	if str == "" {
		log.Fatalln("Icon doesn't exist.")
	}
	fmt.Printf("%s", str)
	// Output:
	// <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 28 28" style="fill: currentColor; vertical-align: top;"><path d="M8.893 1.5c-.183-.31-.52-.5-.887-.5s-.703.19-.886.5L.138 13.499a.98.98 0 000 1.001c.193.31.53.501.886.501h13.964c.367 0 .704-.19.877-.5a1.03 1.03 0 00.01-1.002L8.893 1.5zm.133 11.497H6.987v-2.003h2.039v2.003zm0-3.004H6.987V5.987h2.039v4.006z"></path></svg>
}
