package edit

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/breadysimon/goless/logging"
	"gopkg.in/yaml.v2"
)

var log *logging.Logger = logging.GetLogger()

func applyChanges(inFile, outFile string, scripts [][]string) (err error) {
	var lines []string

	if lines, err = ReadLines(inFile); err == nil {
		for _, s := range scripts {
			var r *regexp.Regexp
			r, err = regexp.Compile(s[1])
			if err == nil {
				switch s[0] {
				case "S":
					for i, _ := range lines {
						lines[i] = r.ReplaceAllString(lines[i], s[2])
					}
				case "D":
					var lines_new []string
					for i, _ := range lines {
						if !r.MatchString(lines[i]) {
							lines_new = append(lines_new, lines[i])
						}
					}
					lines = lines_new
				}
			} else {
				// stop script if any err occurs
				break
			}
		}

		// write back to file
		if err == nil {
			if err = WriteLines(outFile, lines); err == nil {
				return
			}
		}
	}

	return
}

func ReadLines(inFile string) (lines []string, err error) {
	var fin *os.File
	lines = []string{}

	if fin, err = os.Open(inFile); err == nil {
		defer fin.Close()
		scanner := bufio.NewScanner(fin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		err = scanner.Err()
	}
	return
}

func WriteLines(outFile string, lines []string) (err error) {
	var fout *os.File
	if outFile == "" {
		fout = os.Stdout
	} else {
		fout, err = os.Create(outFile)
	}
	if err == nil {
		defer fout.Close()
		delim := ""
		for _, v := range lines {
			if _, err = fout.WriteString(delim + v); err != nil {
				break
			}
			delim = "\n"
		}
		if err == nil {
			return
		}
	}
	return
}

func RunScript(srcFile, inFile, outFile string) (err error) {
	var src []byte

	if srcFile == "-" {
		log.Debug("read yaml scripts from stdin")
		src, err = ioutil.ReadAll(os.Stdin)
	} else {
		log.Debug("read yaml scripts from file:", srcFile)
		src, err = ioutil.ReadFile(srcFile)
	}

	if err == nil {
		err = RunYaml(src, inFile, outFile)
		if err == nil {
			return
		}
	}
	log.Error(err)
	return
}

func RunYaml(src []byte, inFile, outFile string) (err error) {
	var s [][]string
	err = yaml.Unmarshal(src, &s)
	log.Debug("src:", s)
	if err == nil {
		if err = applyChanges(inFile, outFile, s); err == nil {
			return
		}
	}
	log.Error(err)
	return
}
