package regexfinfunction

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

//RegexFinFunction ...
type RegexFinFunction struct {
	Keyword string
	V1      *regexp.Regexp
	V2      *regexp.Regexp
	V3      *regexp.Regexp
	V4      *regexp.Regexp
	V5      *regexp.Regexp

	EL *log.Logger
	LL *log.Logger
	GL *log.Logger
}

//NewRegexFinFunction ...
func NewRegexFinFunction(EL, LL, GL *log.Logger) (*RegexFinFunction, error) {
	var moduleName string = "[RegexFinFunction][NewRegexFinFunction()]"

	if EL == nil || LL == nil || GL == nil {
		return nil, fmt.Errorf("[ERROR]%+v Loggers came empty", moduleName)
	}
	return &RegexFinFunction{
		Keyword: "Fin de Funcion",
		V1:      regexp.MustCompile("^Fin de Funcion"),
		V2:      regexp.MustCompile("^(?i)Fin de Funcion"),
		V3:      regexp.MustCompile("^(?i)Fin de Func"),
		V4:      regexp.MustCompile("^(?i)Fin de Fu"),
		V5:      regexp.MustCompile("^(?i)Fin de P"),
		GL:      GL,
		EL:      EL,
		LL:      LL,
	}, nil
}

//StartsWithFinDeFuncion ...
func (r *RegexFinFunction) StartsWithFinDeFuncion(str string) bool {

	if r.V1.MatchString(str) {
		return true
	}

	if r.V2.MatchString(str) {
		strData := strings.Split(str, " ")
		wrongWord := strData[0]
		Keyword := strings.Split(r.Keyword, "")

		foundTypo := false
		for i, char := range wrongWord {
			if !foundTypo {
				if string(char) != Keyword[i] {
					foundTypo = true

					log.Printf("[ERR] Found typo in '%+v' declaration at [%+v]. Correct syntax should be '%+v'", wrongWord, i, r.Keyword)
					r.GL.Printf("[ERR] Found typo in '%+v' declaration at [%+v]. Correct syntax should be '%+v'", wrongWord, i, r.Keyword)
				}
			}
		}
		return true
	}

	if r.V3.MatchString(str) {
		strData := strings.Split(str, " ")
		wrongWord := strData[0]
		Keyword := strings.Split(r.Keyword, "")
		foundTypo := false
		for i, char := range wrongWord {
			if !foundTypo {
				if string(char) != Keyword[i] {
					foundTypo = true
					log.Printf("[ERR] Found typo in '%+v' declaration at [%+v]. Correct syntax should be '%+v'", wrongWord, i, r.Keyword)
					r.GL.Printf("[ERR] Found typo in '%+v' declaration at [%+v]. Correct syntax should be '%+v'", wrongWord, i, r.Keyword)
				}
			}
		}
		return true
	}

	if r.V4.MatchString(str) {
		strData := strings.Split(str, " ")
		wrongWord := strData[0]
		Keyword := strings.Split(r.Keyword, "")
		foundTypo := false
		for i, char := range wrongWord {
			if !foundTypo {
				if string(char) != Keyword[i] {
					foundTypo = true
					log.Printf("[ERR] Found typo in '%+v' declaration at [%+v]. Correct syntax should be '%+v'", wrongWord, i, r.Keyword)
					r.GL.Printf("[ERR] Found typo in '%+v' declaration at [%+v]. Correct syntax should be '%+v'", wrongWord, i, r.Keyword)
				}
			}
		}
		return true
	}

	if r.V5.MatchString(str) {
		strData := strings.Split(str, " ")
		wrongWord := strData[0]
		Keyword := strings.Split(r.Keyword, "")
		foundTypo := false
		for i, char := range wrongWord {
			if !foundTypo {
				if string(char) != Keyword[i] {
					foundTypo = true
					log.Printf("[ERR] Found typo in '%+v' declaration at [%+v]. Correct syntax should be '%+v'", wrongWord, i, r.Keyword)
					r.GL.Printf("[ERR] Found typo in '%+v' declaration at [%+v]. Correct syntax should be '%+v'", wrongWord, i, r.Keyword)
				}
			}
		}
		return true
	}

	return false
}

//StartsWithFinDeFuncionNoCheck ...
func (r *RegexFinFunction) StartsWithFinDeFuncionNoCheck(str string) bool {
	if r.V1.MatchString(str) {
		return true
	}

	if r.V2.MatchString(str) {
		return true
	}

	if r.V3.MatchString(str) {
		return true
	}

	return false
}
