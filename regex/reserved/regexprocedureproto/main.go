package regexprocedureproto

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

//RegexProcedureProto ...
type RegexProcedureProto struct {
	Keyword string
	V1      *regexp.Regexp
	V2      *regexp.Regexp
	V3      *regexp.Regexp
	V4End   *regexp.Regexp

	EL *log.Logger
	LL *log.Logger
	GL *log.Logger
}

//NewRegexProcedureProto ...
func NewRegexProcedureProto(EL, LL, GL *log.Logger) (*RegexProcedureProto, error) {
	var moduleName string = "[regexProcedureProto][NewRegexProcedureProto()]"

	if EL == nil || LL == nil || GL == nil {
		return nil, fmt.Errorf("[ERROR]%+v Loggers came empty", moduleName)
	}

	compiledV1 := regexp.MustCompile("^(?i)Procedimiento")
	compiledV2 := regexp.MustCompile("^(?i)Procedim")
	compiledV3 := regexp.MustCompile("^(?i)Proc")
	compiledV4End := regexp.MustCompile(";$")
	return &RegexProcedureProto{
		Keyword: "Procedimiento",
		V1:      compiledV1,
		V2:      compiledV2,
		V3:      compiledV3,
		V4End:   compiledV4End,
		EL:      EL,
		LL:      LL,
		GL:      GL,
	}, nil
}

//StartsWithProcedureProto ...
func (r *RegexProcedureProto) StartsWithProcedureProto(str string, lineIndex int64) bool {
	if r.V1.MatchString(str) && r.V4End.MatchString(str) {
		return true
	}

	if r.V2.MatchString(str) && r.V4End.MatchString(str) {
		strData := strings.Split(str, " ")
		wrongWord := strData[0]
		Keyword := strings.Split(r.Keyword, "")
		foundTypo := false
		if len(wrongWord) > len(r.Keyword) {
			r.LogError(lineIndex, 0, wrongWord, fmt.Sprintf("Found typo in '%+v' declaration. Correct syntax should be '%+v'", wrongWord, r.Keyword), str)
			return true
		}
		for i, char := range wrongWord {
			if !foundTypo {
				if string(char) != Keyword[i] {
					foundTypo = true
					r.LogError(lineIndex, i, wrongWord, fmt.Sprintf("Found typo in '%+v' declaration. Correct syntax should be '%+v'", wrongWord, r.Keyword), str)
				}
			}
		}
		return true
	}

	if r.V3.MatchString(str) && r.V4End.MatchString(str) {
		strData := strings.Split(str, " ")
		wrongWord := strData[0]
		Keyword := strings.Split(r.Keyword, "")
		foundTypo := false
		if len(wrongWord) > len(r.Keyword) {
			r.LogError(lineIndex, 0, wrongWord, fmt.Sprintf("Found typo in '%+v' declaration. Correct syntax should be '%+v'", wrongWord, r.Keyword), str)
			return true
		}
		for i, char := range wrongWord {
			if !foundTypo {
				if string(char) != Keyword[i] {
					foundTypo = true
					r.LogError(lineIndex, i, wrongWord, fmt.Sprintf("Found typo in '%+v' declaration. Correct syntax should be '%+v'", wrongWord, r.Keyword), str)
				}
			}
		}
		return true
	}

	return false
}

//StartsWithProcedureProtoNoCheck ...
func (r *RegexProcedureProto) StartsWithProcedureProtoNoCheck(str string) bool {
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

//r.LogError(lineIndex, i, wrongWord, fmt.Sprintf("Found typo in '%+v' declaration. Correct syntax should be '%+v'", wrongWord, r.Keyword), str)

//LogError ...
//"# Linea | # Columna | Error | Descripcion | Linea del Error"
func (r *RegexProcedureProto) LogError(lineIndex int64, columnIndex interface{}, err string, description string, currentLine string) {
	//log.Printf("[ERR] %+v [Line: %+v]", description, lineIndex)
	r.GL.Printf("[ERR] %+v [Line: %+v]", description, lineIndex)
	r.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, columnIndex, err, description, currentLine)
}
