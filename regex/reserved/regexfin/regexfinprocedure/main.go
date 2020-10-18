package regexfinprocedure

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

//RegexFinProcedure ...
type RegexFinProcedure struct {
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

//NewRegexFinProcedure ...
func NewRegexFinProcedure(EL, LL, GL *log.Logger) (*RegexFinProcedure, error) {
	var moduleName string = "[RegexFinProcedure][NewRegexFinProcedure()]"

	if EL == nil || LL == nil || GL == nil {
		return nil, fmt.Errorf("[ERROR]%+v Loggers came empty", moduleName)
	}
	return &RegexFinProcedure{
		Keyword: "Fin de Procedimiento",
		V1:      regexp.MustCompile("^Fin de Procedimiento"),
		V2:      regexp.MustCompile("^(?i)Fin de Procedimi"),
		V3:      regexp.MustCompile("^(?i)Fin de Proce"),
		V4:      regexp.MustCompile("^(?i)Fin de Pro"),
		V5:      regexp.MustCompile("^(?i)Fin de P"),
		GL:      GL,
		EL:      EL,
		LL:      LL,
	}, nil
}

//StartsWithFinDeProcedimiento ...
func (r *RegexFinProcedure) StartsWithFinDeProcedimiento(str string) bool {

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

//StartsWithFinDeProcedimientoNoCheck ...
func (r *RegexFinProcedure) StartsWithFinDeProcedimientoNoCheck(str string) bool {
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
