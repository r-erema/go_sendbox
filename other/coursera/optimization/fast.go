package optimization

import (
	"encoding/json"
	"fmt"
	"github.com/thoas/go-funk"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

var pool = sync.Pool{
	New: func() interface{} {
		return new(User)
	},
}

var wg = &sync.WaitGroup{}
var m = &sync.Mutex{}

var allBrowsers []string
var usersStrings []string

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(fileContents), "\n")

	usersStrings = make([]string, len(lines))

	for i, line := range lines {
		wg.Add(1)
		go func(line string, i int, wg *sync.WaitGroup) {
			defer wg.Done()
			u := pool.Get().(*User)
			err = json.Unmarshal([]byte(line), u)
			if err != nil {
				log.Fatal(err)
			}

			isMSIE, isAndroid := false, false
			for _, browserStr := range u.Browsers {
				if strings.Contains(browserStr, "MSIE") {
					isMSIE = true
					allBrowsers = append(allBrowsers, browserStr)
				} else if strings.Contains(browserStr, "Android") {
					isAndroid = true
					allBrowsers = append(allBrowsers, browserStr)
				}
			}

			if isMSIE && isAndroid {
				m.Lock()
				email := strings.Replace(u.Email, "@", " [at] ", -1)
				usersStrings[i] = fmt.Sprintf("[%d] %s <%s>\n", i, u.Name, email)
				m.Unlock()
			}

			pool.Put(u)
		}(line, i, wg)
	}

	wg.Wait()

	_, _ = fmt.Fprintln(out, "found users:\n"+strings.Join(usersStrings, ""))
	_, _ = fmt.Fprintln(out, "Total unique browsers", len(funk.UniqString(allBrowsers)))

}
