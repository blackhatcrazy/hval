package preprocess

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
)

var keyIsTmpl = regexp.MustCompile(`({{.*}}.*:.*)`)
var valIsTmpl = regexp.MustCompile(`(\w*: )({{.*}}.*)`)

func NewValues(input, output string) (Values, error) {
	if output == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return Values{}, err
		}
		output = fmt.Sprintf("%s/debug_valuesSanitized.yaml", cwd)
	}
	return Values{
		InputPath:       input,
		DebugOutputPath: output,
	}, nil

}

type Values struct {
	InputPath       string
	DebugOutputPath string
	file            *os.File
	Parsed          map[string]interface{}
}

func (v *Values) Open() error {
	var err error
	v.file, err = os.Open(v.InputPath)
	return err
}

func createDebug(outputPath string) (*os.File, error) {
	err := os.MkdirAll(path.Dir(outputPath), os.ModePerm)
	if err != nil {
		return nil, err
	}
	return os.Create(outputPath)
}

func (v *Values) Sanitize(debug bool) error {
	var out *os.File
	var err error
	defer v.file.Close()
	defer out.Close()

	if debug {
		out, err = createDebug(v.DebugOutputPath)
		if err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(v.file)
	for scanner.Scan() {
		line := scanner.Bytes()

		if len(keyIsTmpl.FindAll(line, -1)) > 0 {
			return fmt.Errorf(
				"illegal key found in \"%s\", line \"%s\"",
				v.InputPath, line,
			)
		}
		// stringify all occurances of template inputs
		matchTmpl := valIsTmpl.ReplaceAll(line, []byte("${1}'${2}'"))
		fmt.Println(string(matchTmpl))
		if debug {
			if _, err := fmt.Fprintln(out, string(matchTmpl)); err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// func copy(src, dst string) error {
// 	sourceFileStat, err := os.Stat(src)
// 	if err != nil {
// 		return err
// 	}

// 	if !sourceFileStat.Mode().IsRegular() {
// 		return fmt.Errorf("%s is not a regular file", src)
// 	}

// 	source, err := os.Open(src)
// 	if err != nil {
// 		return err
// 	}
// 	defer source.Close()

// 	destination, err := os.Create(dst)
// 	if err != nil {
// 		return err
// 	}
// 	defer destination.Close()
// 	_, err = io.Copy(destination, source)
// 	return err
// }
