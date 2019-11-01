package edit

import (
	"io/ioutil"
	"os"
	"regexp"

	"github.com/breadysimon/goless/logging"
)

var log *logging.Logger = logging.GetLogger()

type EditStep struct {
	cmd, exp, sub string
}

/*
EditFile
  RegExp replacement in file.
  src format:
  {
	  cmd: "S" - global replace, "D" - global delete
	  exp: register expression
	  sub: subsitution
  }
*/
func EditFile(filename string, src []EditStep) (out string, err error) {
	if txt, err := ioutil.ReadFile(filename); err == nil {
		in := string(txt)

		for _, s := range src {

			var exp string

			// pre-process reg expression for each command
			switch s.cmd {
			case "S":
				exp = s.exp
			case "D":
				exp = `(?m)[\r\n]+^.*` + s.exp + ".*$"
				s.sub = ""
			}

			// do replacement
			if r, err := regexp.Compile(exp); err == nil {
				out = r.ReplaceAllString(in, s.sub)
			} else {
				break
			}
		}

		// write back to file
		if err == nil {
			ioutil.WriteFile(filename, []byte(out), os.ModePerm)
		}
		return out, err
	}
	log.Error(err)
	return
}
