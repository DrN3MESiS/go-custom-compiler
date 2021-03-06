package lexyc

import (
	"go-custom-compiler/helpers"
	"go-custom-compiler/models"
	"regexp"
	"strings"
)

//NextProcedureProto ...
func (l *LexicalAnalyzer) NextProcedureProto(currentLine string, lineIndex int64, debug bool) {
	funcName := "[NextProcedureProto()] "
	// var moduleName string = "[regexfunctionproto][NextProcedureProto()]"

	if l.CurrentBlockType == models.PROCEDUREPROTOBLOCK {
		if l.R.RegexProcedureProto.StartsWithProcedureProtoNoCheck(currentLine) {
			currentLine = strings.Join(strings.Split(currentLine, " ")[1:], " ")
		}
		currentLine = strings.TrimSpace(currentLine)

		if l.R.RegexProcedureProtoEntero.MatchProcedureEntero(currentLine) {
			procName, procParamType, procParamVars := getDataFromProcedureProto(currentLine)
			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				procName, helpers.IDENTIFICADOR,
				"(", helpers.DELIMITADOR,
			}))

			paramType := models.VarTypeToTokenType(procParamType)
			params := []*models.Token{}
			vars := strings.Split(procParamVars, ", ")
			for i, procParamVar := range vars {
				l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{procParamVar, helpers.IDENTIFICADOR}))
				params = append(params, &models.Token{Type: paramType, Key: procParamVar})
				if i != len(vars)-1 {
					l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{",", helpers.DELIMITADOR}))
				}
			}
			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				":", helpers.DELIMITADOR,
				procParamType, helpers.IDENTIFICADOR,
				")", helpers.DELIMITADOR,
				";", helpers.DELIMITADOR,
			}))
			l.FunctionStorage = append(l.FunctionStorage, &models.TokenFunc{Key: procName, Params: params})

			l.GL.Printf("%+v[PROCEDURE PROTO] Entero Found > %+v", funcName, currentLine)
			if debug {
				// //log.Printf("[PROCEDURE PROTO] Entero Found > %+v", currentLine)
			}
			return
		}
		if l.R.RegexProcedureProtoLogico.MatchProcedureLogico(currentLine) {
			procName, procParamType, procParamVars := getDataFromProcedureProto(currentLine)
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{procName, helpers.IDENTIFICADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"(", helpers.DELIMITADOR}))

			vars := strings.Split(procParamVars, ", ")
			for i, procParamVar := range vars {
				l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{procParamVar, helpers.IDENTIFICADOR}))
				if i != len(vars)-1 {
					l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{",", helpers.DELIMITADOR}))
				}
			}
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{":", helpers.DELIMITADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{procParamType, helpers.IDENTIFICADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{")", helpers.DELIMITADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{";", helpers.DELIMITADOR}))

			l.GL.Printf("%+v[PROCEDURE PROTO] Logico Found > %+v", funcName, currentLine)
			if debug {
				// //log.Printf("[PROCEDURE PROTO] Logico Found > %+v", currentLine)
			}
			return
		}
		if l.R.RegexProcedureProtoReal.MatchProcedureReal(currentLine) {
			procName, procParamType, procParamVars := getDataFromProcedureProto(currentLine)
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{procName, helpers.IDENTIFICADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"(", helpers.DELIMITADOR}))

			vars := strings.Split(procParamVars, ", ")
			for i, procParamVar := range vars {
				l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{procParamVar, helpers.IDENTIFICADOR}))
				if i != len(vars)-1 {
					l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{",", helpers.DELIMITADOR}))
				}
			}
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{":", helpers.DELIMITADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{procParamType, helpers.IDENTIFICADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{")", helpers.DELIMITADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{";", helpers.DELIMITADOR}))

			l.GL.Printf("%+v[PROCEDURE PROTO] Real Found > %+v", funcName, currentLine)
			if debug {
				// //log.Printf("[PROCEDURE PROTO] Real Found > %+v", currentLine)
			}
			return
		}
		if l.R.RegexProcedureProtoAlfabetico.MatchProcedureAlfabetico(currentLine) {
			procName, procParamType, procParamVars := getDataFromProcedureProto(currentLine)
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{procName, helpers.IDENTIFICADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"(", helpers.DELIMITADOR}))

			vars := strings.Split(procParamVars, ", ")
			for i, procParamVar := range vars {
				l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{procParamVar, helpers.IDENTIFICADOR}))
				if i != len(vars)-1 {
					l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{",", helpers.DELIMITADOR}))
				}
			}
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{":", helpers.DELIMITADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{procParamType, helpers.IDENTIFICADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{")", helpers.DELIMITADOR}))
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{";", helpers.DELIMITADOR}))

			l.GL.Printf("%+v[PROCEDURE PROTO] Alfabetico Found > %+v", funcName, currentLine)
			if debug {
				// //log.Printf("[PROCEDURE PROTO] Alfabetico Found > %+v", currentLine)
			}
			return
		}

		if l.R.RegexProcedureProtoDefault.MatchProcedureDefault(currentLine) {
			_, procParamType, _ := getDataFromProcedureProto(currentLine)

			regexEntero := regexp.MustCompile(`^(?i)entero`)
			regexReal := regexp.MustCompile(`^(?i)real`)
			regexLogico := regexp.MustCompile(`^(?i)logico`)
			regexAlfabetico := regexp.MustCompile(`^(?i)alfabetico`)

			if regexEntero.MatchString(procParamType) {
				l.GL.Printf("%+v[PROCEDURE PROTO] Entero Found > %+v", funcName, currentLine)
				if debug {
					// //log.Printf("[PROCEDURE PROTO] Entero Found > %+v", currentLine)
				}
				foundTypo := false
				keyData := strings.Split(l.R.RegexProcedureProtoEntero.Keyword, "")
				for i, char := range procParamType {
					if i < len(keyData)-1 {
						if !foundTypo {
							if string(char) != keyData[i] {
								foundTypo = true
								// //log.Printf("[ERR] Found typo in '%+v' declaration at [%+v][Line: %+v]. Correct syntax should be '%+v'", procParamType, i, lineIndex, l.R.RegexProcedureProtoEntero.Keyword)
								l.GL.Printf("[ERR] Found typo in '%+v' declaration at [%+v][Line: %+v]. Correct syntax should be '%+v'", procParamType, i, lineIndex, l.R.RegexProcedureProtoEntero.Keyword)
								//"# Linea | # Columna | Error | Descripcion | Linea del Error"
								l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, i, procParamType, l.R.RegexProcedureProtoEntero.Keyword, currentLine)
							}
						}
					}
				}
				return
			}

			if regexReal.MatchString(procParamType) {
				l.GL.Printf("%+v[PROCEDURE PROTO] Real Found > %+v", funcName, currentLine)
				if debug {
					// //log.Printf("[PROCEDURE PROTO] Real Found > %+v", currentLine)
				}
				foundTypo := false
				keyData := strings.Split(l.R.RegexProcedureProtoReal.Keyword, "")
				for i, char := range procParamType {
					if i < len(keyData)-1 {
						if !foundTypo {
							if string(char) != keyData[i] {
								foundTypo = true
								// //log.Printf("[ERR] Found typo in '%+v' declaration at [%+v][Line: %+v]. Correct syntax should be '%+v'", procParamType, i, lineIndex, l.R.RegexProcedureProtoReal.Keyword)
								l.GL.Printf("[ERR] Found typo in '%+v' declaration at [%+v][Line: %+v]. Correct syntax should be '%+v'", procParamType, i, lineIndex, l.R.RegexProcedureProtoReal.Keyword)
								//"# Linea | # Columna | Error | Descripcion | Linea del Error"
								l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, i, procParamType, l.R.RegexProcedureProtoReal.Keyword, currentLine)
							}
						}
					}
				}
				return
			}

			if regexLogico.MatchString(procParamType) {
				l.GL.Printf("%+v[PROCEDURE PROTO] Logico Found > %+v", funcName, currentLine)
				if debug {
					// //log.Printf("[PROCEDURE PROTO] Logico Found > %+v", currentLine)
				}
				foundTypo := false
				keyData := strings.Split(l.R.RegexProcedureProtoLogico.Keyword, "")
				for i, char := range procParamType {
					if i < len(keyData)-1 {
						if !foundTypo {
							if string(char) != keyData[i] {
								foundTypo = true
								// //log.Printf("[ERR] Found typo in '%+v' declaration at [%+v][Line: %+v]. Correct syntax should be '%+v'", procParamType, i, lineIndex, l.R.RegexProcedureProtoLogico.Keyword)
								l.GL.Printf("[ERR] Found typo in '%+v' declaration at [%+v][Line: %+v]. Correct syntax should be '%+v'", procParamType, i, lineIndex, l.R.RegexProcedureProtoLogico.Keyword)
								//"# Linea | # Columna | Error | Descripcion | Linea del Error"
								l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, i, procParamType, l.R.RegexProcedureProtoLogico.Keyword, currentLine)
							}
						}
					}
				}
				return
			}

			if regexAlfabetico.MatchString(procParamType) {
				l.GL.Printf("%+v[PROCEDURE PROTO] Alfabetico Found > %+v", funcName, currentLine)
				if debug {
					// //log.Printf("[PROCEDURE PROTO] Alfabetico Found > %+v", currentLine)
				}
				foundTypo := false
				keyData := strings.Split(l.R.RegexProcedureProtoAlfabetico.Keyword, "")
				for i, char := range procParamType {
					if i < len(keyData)-1 {
						if !foundTypo {
							if string(char) != keyData[i] {
								foundTypo = true
								// //log.Printf("[ERR] Found typo in '%+v' declaration at [%+v][Line: %+v]. Correct syntax should be '%+v'", procParamType, i, lineIndex, l.R.RegexProcedureProtoAlfabetico.Keyword)
								l.GL.Printf("[ERR] Found typo in '%+v' declaration at [%+v][Line: %+v]. Correct syntax should be '%+v'", procParamType, i, lineIndex, l.R.RegexProcedureProtoAlfabetico.Keyword)
								//"# Linea | # Columna | Error | Descripcion | Linea del Error"
								l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, i, procParamType, l.R.RegexProcedureProtoAlfabetico.Keyword, currentLine)
							}
						}
					}
				}
				return
			}

		}

		// l.GL.Printf("%+v Did not found any type of match on Line[%+v]! ", funcName, lineIndex)

	}
}

func getDataFromProcedureProto(currentLine string) (string, string, string) {
	currentLine = strings.TrimSuffix(currentLine, ";")
	lineData := splitAtCharRespectingQuotes(currentLine, '(')
	procedureName := lineData[0]
	procedureParamsToParse := lineData[1]
	procedureParamsToParse = strings.TrimSuffix(procedureParamsToParse, ")")
	paramsData := splitAtCharRespectingQuotes(procedureParamsToParse, ':')
	procedureParamType := paramsData[1]
	procedureParamVars := paramsData[0]

	return procedureName, procedureParamType, procedureParamVars
}
