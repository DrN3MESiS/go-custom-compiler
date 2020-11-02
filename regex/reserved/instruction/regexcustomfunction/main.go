package regexcustomfunction

import (
	"log"
	"regexp"
)

//RegexCustomFunction ...
type RegexCustomFunction struct {
	Keyword string

	CustomFunction *regexp.Regexp
	ENDSWITH       *regexp.Regexp

	EL *log.Logger
	LL *log.Logger
	GL *log.Logger
}

//NewRegexCustomFunction ...
func NewRegexCustomFunction(EL, LL, GL *log.Logger) (*RegexCustomFunction, error) {
	// var moduleName string = "[regexint][NewRegexCustomFunction()]"

	return &RegexCustomFunction{
		Keyword:        "CustomFunction",
		CustomFunction: regexp.MustCompile(`^([a-zA-Z]+[a-zA-Z0-9](\s*)*)\((.*)\)`),
		ENDSWITH:       regexp.MustCompile(`;$`),

		GL: GL,
		EL: EL,
		LL: LL,
	}, nil
}

//MatchPC ...
func (r *RegexCustomFunction) MatchPC(str string, lineIndex int64) bool {
	if r.ENDSWITH.MatchString(str) {
		return true
	}

	return false
}

//MatchCustomFunction ...
func (r *RegexCustomFunction) MatchCustomFunction(str string, lineIndex int64) bool {
	if r.CustomFunction.MatchString(str) {
		return true
	}

	return false
}

//LogError ...
//"# Linea | # Columna | Error | DescripcCustomFunctionn | Linea del Error"
func (r *RegexCustomFunction) LogError(lineIndex int64, columnIndex interface{}, err string, descriptCustomFunctionn string, currentLine string) {
	log.Printf("[ERR] %+v [Line: %+v]", descriptCustomFunctionn, lineIndex)
	r.GL.Printf("[ERR] %+v [Line: %+v]", descriptCustomFunctionn, lineIndex)
	r.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, columnIndex, err, descriptCustomFunctionn, currentLine)
}
