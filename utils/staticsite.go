package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/harshitsinghai/gogist/models"
)

// Write writes the buffer to the target file.
func Write(tgt string, buf []byte) {
	os.MkdirAll(filepath.Dir(tgt), 0777)
	if err := ioutil.WriteFile(tgt, buf, 0755); err != nil {
		log.Fatalf("must write: %s: %s", tgt, err)
	}
}

func timeFormat(t time.Time) string {
	return t.Format("02 January 2006")
}

// GenerateTimeline takes data and injects it inside the html file
func GenerateTimeline(metadata []models.Repo) {
	// for _, val := range metadata {
	// 	fmt.Println(val.Name)
	// }

	funcMap := template.FuncMap{
		"timeFormat": timeFormat,
	}

	t, err := template.New("timeline.html").Funcs(funcMap).ParseFiles("./utils/timeline.html") //parse the html file homepage.html
	if err != nil {                                                                            // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	output := bytes.Buffer{}
	err = t.Execute(&output, metadata) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {
		log.Fatal("Some error occured", err)
	}
	Write("timeline.html", output.Bytes())
}
