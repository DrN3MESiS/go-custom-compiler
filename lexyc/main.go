package lexyc

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strings"

	"go-custom-compiler/helpers"
	"go-custom-compiler/models"
	"go-custom-compiler/regex"

	"github.com/DrN3MESiS/pprnt"
)

//LexicalAnalyzer ...
type LexicalAnalyzer struct {
	File *bufio.Scanner     //File
	R    *regex.CustomRegex //Regex Handler
	EL   *log.Logger        //Error Logger
	LL   *log.Logger        //Lex Logger
	GL   *log.Logger        //General Logger
	TEST *log.Logger

	//TEST
	Status           int
	CurrentBlockType models.BlockType
	ParentBlockType  models.BlockType
	BlockQueue       []models.BlockType
	OpQueue          []models.TokenComp
	NamesQueue       []string
	ConstantStorage  []*models.Token
	VariableStorage  []*models.Token
	FunctionStorage  []*models.TokenFunc
	Context          string
	HasMain          bool
	ErrorsCount      int
}

//NewLexicalAnalyzer ...
func NewLexicalAnalyzer(file *bufio.Scanner, ErrorLogger, LexLogger, GeneralLogger, TestLogger *log.Logger) (*LexicalAnalyzer, error) {
	var moduleName string = "[Lexyc][NewLexicalAnalyzer()]"
	GeneralLogger.Printf("Started the Lexical Analyzer")

	if file == nil {
		GeneralLogger.Printf("[ERR]%+v file is not present", moduleName)
		return nil, fmt.Errorf("[ERR]%+v file is not present", moduleName)
	}
	if ErrorLogger == nil || LexLogger == nil || GeneralLogger == nil {
		GeneralLogger.Printf("[ERR]%+v Loggers are not present", moduleName)
		return nil, fmt.Errorf("[ERR]%+v Loggers are not present", moduleName)
	}
	R, err := regex.NewRegex(ErrorLogger, LexLogger, GeneralLogger)
	if err != nil {
		GeneralLogger.Printf("[ERR]%+v %+v", moduleName, err.Error())
		return nil, fmt.Errorf("[ERR]%+v %+v", moduleName, err.Error())
	}

	LexLogger.Println("--------------------------------------------------------------------------------------------")
	LexLogger.Println(helpers.IndentString(helpers.LEXINDENT, []string{"Lexema", "Token"}))
	LexLogger.Println("--------------------------------------------------------------------------------------------")
	ErrorLogger.Printf("=============================================================")
	ErrorLogger.Printf("# Linea | # Columna | Error | Descripcion | Linea del Error")
	ErrorLogger.Printf("=============================================================")

	return &LexicalAnalyzer{
		File: file,
		R:    R,
		EL:   ErrorLogger,
		LL:   LexLogger,
		GL:   GeneralLogger,
		TEST: TestLogger,

		Status:           0,
		ParentBlockType:  models.NULLBLOCK,
		BlockQueue:       []models.BlockType{},
		CurrentBlockType: models.NULLBLOCK,
		OpQueue:          []models.TokenComp{},
		NamesQueue:       []string{},
		ConstantStorage:  []*models.Token{},
		VariableStorage:  []*models.Token{},
		Context:          "Global",
	}, nil
}

//Analyze ...
func (l *LexicalAnalyzer) Analyze(debug bool) error {
	funcName := "[Analyze()]"
	var lineIndex int64 = 1
	for l.File.Scan() {
		currentLine := l.File.Text()
		foundSomething := false

		if len(currentLine) == 0 {
			l.GL.Printf("%+v Skipped [Line: %+v]; Reason: Empty", funcName, lineIndex)
			lineIndex++

			continue
		}
		var LastBlockState models.BlockType
		LastBlockState = l.CurrentBlockType
		/* Type Validation */
		isComment, err := l.R.StartsWith("//", currentLine)
		if err != nil {
			l.GL.Printf("%+v[APP_ERR] %+v", funcName, err.Error())
			return fmt.Errorf("%+v[APP_ERR] %+v", funcName, err.Error())
		}
		if isComment {
			l.GL.Printf("%+vSkipping Comment at line %+v", funcName, lineIndex)
			if debug {
				log.Printf("Skipping Comment at line %+v", lineIndex)
			}
			lineIndex++

			continue
		}

		// log.Printf("> Line: [%+v] CurrentBlock: %+v", lineIndex, l.CurrentBlockType)

		currentLine = strings.TrimSpace(currentLine)

		// log.Printf("BLOCK [Line:%+v]['%+v'] > %+v\n", lineIndex, currentLine, l.BlockQueue)
		// log.Printf("BLOCK [Line:%+v] > %+v\n", lineIndex, l.BlockQueue)

		/* StartsWith */

		//Contante
		if l.R.RegexConstante.StartsWithConstante(currentLine, lineIndex) {
			l.CurrentBlockType = models.CONSTANTBLOCK
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"constantes", helpers.PALABRARESERVADA}))
			foundSomething = true
		}

		//Variable
		if l.R.RegexVariable.StartsWithVariable(currentLine, lineIndex) {
			l.CurrentBlockType = models.VARIABLEBLOCK
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"variables", helpers.PALABRARESERVADA}))
			foundSomething = true
		}

		//FunctionProto
		if l.R.RegexFuncionProto.StartsWithFuncionProto(currentLine, lineIndex) && l.ParentBlockType == models.NULLBLOCK {
			l.CurrentBlockType = models.FUNCTIONPROTOBLOCK
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"funcion", helpers.PALABRARESERVADA}))
			foundSomething = true
		}

		//ProcedureProto
		if l.R.RegexProcedureProto.StartsWithProcedureProto(currentLine, lineIndex) && l.ParentBlockType == models.NULLBLOCK {
			l.CurrentBlockType = models.PROCEDUREPROTOBLOCK
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"procedimiento", helpers.PALABRARESERVADA}))
			foundSomething = true
		}

		//Procedure
		if l.R.RegexProcedure.StartsWithProcedure(currentLine, lineIndex) {
			// if l.CurrentBlockType != models.NULLBLOCK {
			// 	l.CurrentBlockType = models.NULLBLOCK
			// }
			l.GL.Println()

			if len(l.BlockQueue) > 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to create new procedure without finalizing the last Function or Procedure", currentLine)
				l.BlockQueue = []models.BlockType{}
			}
			l.BlockQueue = append(l.BlockQueue, models.PROCEDUREBLOCK)
			procedureGroups := helpers.GetGroupMatches(currentLine, helpers.PROCEDIMIENTOREGEXP)
			token := []string{
				procedureGroups[0], helpers.PALABRARESERVADA,
				procedureGroups[1], helpers.IDENTIFICADOR,
				"(", helpers.DELIMITADOR,
			}

			l.Context = procedureGroups[1]

			symbol := models.TokenFunc{Key: procedureGroups[1], IsDefined: true}

			params := strings.Join(procedureGroups[2:], "")
			groups := strings.Split(params, ";")
			for i, group := range groups {
				if i > 0 {
					token = append(token, []string{";", helpers.DELIMITADOR}...)
				}
				groupVars := strings.Split(group, ":")

				paramType := models.VarTypeToTokenType(groupVars[len(groupVars)-1])

				vars := strings.Split(groupVars[0], ",")
				if vars[0] != "" {
					token = append(token, []string{vars[0], helpers.IDENTIFICADOR}...)
					symbol.Params = append(symbol.Params, &models.Token{Type: paramType, Key: vars[0]})
				}

				for _, v := range vars[1:] {
					v = strings.TrimSpace(v)
					token = append(token, []string{
						",", helpers.DELIMITADOR,
					}...)
					token = append(token, l.AnalyzeType("", 0, v)...)

					symbol.Params = append(symbol.Params, &models.Token{Type: paramType, Key: v})
				}
				if vars[0] != "" {
					token = append(token, []string{
						":", helpers.DELIMITADOR,
						strings.TrimSpace(groupVars[len(groupVars)-1]), helpers.PALABRARESERVADA,
					}...)
				}
			}
			token = append(token, []string{
				")", helpers.DELIMITADOR,
			}...)
			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, token))

			funcProto := l.FindFunction(currentLine, lineIndex, symbol.Key)
			if funcProto == nil {
				l.FunctionStorage = append(l.FunctionStorage, &symbol)
			} else {
				l.CompareFunction(currentLine, lineIndex, funcProto, &symbol, false)
			}

			foundSomething = true
		}

		//Function
		if l.R.RegexFunction.StartsWithFunction(currentLine, lineIndex) {
			l.GL.Println()

			if len(l.BlockQueue) > 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to create new function without finalizing the last Function or Procedure", currentLine)
				l.BlockQueue = []models.BlockType{}
			}

			l.BlockQueue = append(l.BlockQueue, models.FUNCTIONBLOCK)
			funcionGroups := helpers.GetGroupMatches(currentLine, helpers.FUNCIONREGEXP)
			token := []string{
				funcionGroups[0], helpers.PALABRARESERVADA,
				funcionGroups[1], helpers.IDENTIFICADOR,
				"(", helpers.DELIMITADOR,
			}

			l.Context = funcionGroups[1]
			symbol := models.TokenFunc{Key: funcionGroups[1], IsDefined: true}

			params := strings.Join(funcionGroups[2:len(funcionGroups)-1], "")
			groups := strings.Split(params, ";")
			for i, group := range groups {
				if i > 0 {
					token = append(token, []string{";", helpers.DELIMITADOR}...)
				}
				groupVars := strings.Split(group, ":")
				vars := strings.Split(groupVars[0], ",")

				paramType := models.VarTypeToTokenType(groupVars[len(groupVars)-1])
				symbol.Params = append(symbol.Params, &models.Token{Type: paramType, Key: vars[0]})

				token = append(token, []string{vars[0], helpers.IDENTIFICADOR}...)
				for _, v := range vars[1:] {
					v = strings.TrimSpace(v)
					token = append(token, []string{
						",", helpers.DELIMITADOR,
					}...)
					token = append(token, l.AnalyzeType("", 0, v)...)

					symbol.Params = append(symbol.Params, &models.Token{Type: paramType, Key: v})
				}
				token = append(token, []string{
					":", helpers.DELIMITADOR,
					strings.TrimSpace(groupVars[len(groupVars)-1]), helpers.PALABRARESERVADA,
				}...)
			}
			token = append(token, []string{
				")", helpers.DELIMITADOR,
				":", helpers.DELIMITADOR,
				funcionGroups[len(funcionGroups)-1], helpers.PALABRARESERVADA,
			}...)
			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, token))

			funcType := models.VarTypeToTokenType(funcionGroups[len(funcionGroups)-1])
			symbol.Type = funcType

			procProto := l.FindFunction(currentLine, lineIndex, symbol.Key)
			if procProto == nil {
				l.FunctionStorage = append(l.FunctionStorage, &symbol)
			} else {
				l.CompareFunction(currentLine, lineIndex, procProto, &symbol, false)
			}

			foundSomething = true
		}

		//Inicio
		if l.R.RegexInicio.StartsWithInicio(currentLine, lineIndex) {
			if len(l.BlockQueue) == 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to initialize something outside of a Block", currentLine)
			}

			switch l.BlockQueue[len(l.BlockQueue)-1] {
			case models.INITBLOCK:
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to initialize something when already initialized", currentLine)
				break
			case models.PROCEDUREBLOCK, models.FUNCTIONBLOCK, models.CUANDOBLOCK:
				l.GL.Printf("%+v Initialized a %+v [Line: %+v]", funcName, l.BlockQueue[len(l.BlockQueue)-1], lineIndex)
				l.BlockQueue = append(l.BlockQueue, models.INITBLOCK)
				l.CurrentBlockType = models.NULLBLOCK
				break

			default:
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to initialize something non existent", currentLine)
				break
			}

			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"Inicio", helpers.PALABRARESERVADA}))

			foundSomething = true
		}

		//FinDeFuncion
		if l.R.RegexFinFunction.StartsWithFinDeFuncion(currentLine, lineIndex) {
			if len(l.BlockQueue) == 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a FUNCTIONBLOCK outside of a FUNCTIONBLOCK", currentLine)
			}

			if l.BlockQueue[len(l.BlockQueue)-1] != models.INITBLOCK {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a FUNCTIONBLOCK that wasn't initialized", currentLine)
			}

			newArr, ok := helpers.RemoveFromQueue(l.BlockQueue, models.INITBLOCK)
			if ok {
				l.BlockQueue = newArr
			} else {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a FUNCTIONBLOCK that wasn't initialized", currentLine)
			}

			newArr, ok = helpers.RemoveFromQueue(l.BlockQueue, models.FUNCTIONBLOCK)
			if ok {
				l.BlockQueue = newArr
			} else {
				if helpers.QueueContainsBlock(l.BlockQueue, models.PROCEDUREBLOCK) {
					l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a FUNCTIONBLOCK:Inicio with a PROCEDUREBLOCK as parent", currentLine)
				} else {
					l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a FUNCTIONBLOCK outside of a FUNCTIONBLOCK", currentLine)
				}

			}
			l.GL.Println()
			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				"Fin", helpers.PALABRARESERVADA,
				"de", helpers.PALABRARESERVADA,
				"funcion", helpers.PALABRARESERVADA,
				";", helpers.DELIMITADOR,
			}))

			l.Context = "Global"
			foundSomething = true
		}

		//FinDeProcedimiento
		if l.R.RegexFinProcedure.StartsWithFinDeProcedimiento(currentLine, lineIndex) {
			if len(l.BlockQueue) == 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a PROCEDUREBLOCK outside of a PROCEDUREBLOCK", currentLine)
			}

			newArr, ok := helpers.RemoveFromQueue(l.BlockQueue, models.INITBLOCK)
			if ok {
				l.BlockQueue = newArr
			} else {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a PROCEDUREBLOCK that wasn't initialized", currentLine)
			}

			newArr, ok = helpers.RemoveFromQueue(l.BlockQueue, models.PROCEDUREBLOCK)
			if ok {
				l.BlockQueue = newArr
			} else {
				if helpers.QueueContainsBlock(l.BlockQueue, models.FUNCTIONBLOCK) {
					l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a PROCEDUREBLOCK:Inicio with a FUNCTIONBLOCK as parent", currentLine)
				} else {
					l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a PROCEDUREBLOCK outside of a PROCEDUREBLOCK", currentLine)
				}
			}
			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				"Fin", helpers.PALABRARESERVADA,
				"de", helpers.PALABRARESERVADA,
				"procedimiento", helpers.PALABRARESERVADA,
				";", helpers.DELIMITADOR,
			}))

			l.Context = "Global"
			foundSomething = true
		}

		//Fin
		if l.R.RegexFin.StartsWithFin(currentLine, lineIndex) {
			if !l.R.RegexIO.MatchPC(currentLine, lineIndex) {
				l.LogError(lineIndex, len(currentLine)-1, ";", "Missing ';'", currentLine)
			}

			newArr, ok := helpers.RemoveFromQueue(l.BlockQueue, models.INITBLOCK)
			if ok {
				l.BlockQueue = newArr
			} else {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a SOMETHING that wasn't initialized", currentLine)
			}

			switch l.BlockQueue[len(l.BlockQueue)-1] {
			case models.CUANDOBLOCK:
				newArr, ok = helpers.RemoveFromQueue(l.BlockQueue, models.CUANDOBLOCK)
				if ok {
					l.BlockQueue = newArr
				}
				break
			default:
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a SOMETHING:Inicio that didn't exist", currentLine)
				break
			}
			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				"Fin", helpers.PALABRARESERVADA,
				";", helpers.DELIMITADOR,
			}))

			foundSomething = true
		}

		//Repetir
		if l.R.RegexLoopRepetir.StartsWithRepetir(currentLine, lineIndex) {
			if len(l.BlockQueue) == 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to create a REPEATBLOCK outside of a BLOCK", currentLine)
			}

			l.BlockQueue = append(l.BlockQueue, models.REPEATBLOCK)
			l.GL.Printf("%+v Initialized a REPEATBLOCK [Line: %+v]", funcName, lineIndex)

			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"repetir", helpers.PALABRARESERVADA}))

			foundSomething = true
		}

		//Hasta Que (Repetir)
		if l.R.RegexLoopHastaQue.StartsWithHastaQue(currentLine, lineIndex) {
			if len(l.BlockQueue) == 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to end a REPEATBLOCK outside of a BLOCK", currentLine)
			}

			/* Analyze Params */

			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				"hasta", helpers.PALABRARESERVADA,
				"que", helpers.PALABRARESERVADA,
				"(", helpers.DELIMITADOR,
			}))

			l.OpQueue = []models.TokenComp{}
			l.NamesQueue = []string{}

			groups := helpers.GetGroupMatches(currentLine, helpers.HASTAQUEREGEXP)
			if len(groups) > 0 {
				params := groups[0]
				l.AnalyzeParams(currentLine, lineIndex, params)
			} else {
				l.LogError(lineIndex, "N/A", "N/A", "Instruction 'Hasta que' doesn't have params", currentLine)
			}

			l.AnalyzeFuncQueue(currentLine, lineIndex)

			/* End Analyze Params*/

			if l.BlockQueue[len(l.BlockQueue)-1] == models.REPEATBLOCK {
				newArr, ok := helpers.RemoveFromQueue(l.BlockQueue, models.REPEATBLOCK)
				if ok {
					l.BlockQueue = newArr
				} else {
					l.LogErrorGeneral(lineIndex, "N/A", "N/A", "I tried to delete something that was inside the slice that I saw before trying to delete", currentLine)
				}
			} else {
				l.LogError(lineIndex, "N/A", "N/A", fmt.Sprintf("Attempted to end a REPEATBLOCK before finalizing a %+v", l.BlockQueue[len(l.BlockQueue)-1]), currentLine)
			}

			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				")", helpers.DELIMITADOR,
				";", helpers.DELIMITADOR,
			}))

			foundSomething = true
		}

		//ImprimeNL
		if l.R.RegexIO.MatchImprimenl(currentLine, lineIndex) {
			if !l.R.RegexIO.MatchPC(currentLine, lineIndex) {
				l.LogError(lineIndex, len(currentLine)-1, ";", "Missing ';'", currentLine)
			}

			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				"Imprimenl", helpers.PALABRARESERVADA,
				"(", helpers.DELIMITADOR,
			}))

			params := l.R.RegexImprime.GroupsImprime(currentLine)
			params = strings.Split(params[len(params)-1], ",")
			l.OpQueue = []models.TokenComp{}
			l.NamesQueue = []string{}
			for i, str := range params {
				str = strings.TrimSpace(str)
				token := l.AnalyzeType(currentLine, lineIndex, str)
				if i != len(params)-1 {
					token = append(token, []string{",", helpers.DELIMITADOR}...)
					l.OpQueue = append(l.OpQueue, models.DELIM)
				}
				if len(token) > 0 {
					l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, token))
				}
			}

			l.AnalyzeOpQueue(currentLine, lineIndex)
			l.AnalyzeFuncQueue(currentLine, lineIndex)
			if !l.ExpectNoNone() {
				l.LogError(lineIndex, "N/A", "N/A", "One of the parameters introduced is not valid", currentLine)
			}

			l.GL.Printf("%+v Found 'Imprimenl' instruction [Line: %+v]", funcName, lineIndex)

			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				")", helpers.DELIMITADOR,
				";", helpers.DELIMITADOR,
			}))

			//Imprime
			foundSomething = true
		} else if l.R.RegexIO.MatchImprime(currentLine, lineIndex) {
			if !l.R.RegexIO.MatchPC(currentLine, lineIndex) {
				l.LogError(lineIndex, len(currentLine)-1, ";", "Missing ';'", currentLine)
			}
			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				"Imprime", helpers.PALABRARESERVADA,
				"(", helpers.DELIMITADOR,
			}))

			params := l.R.RegexImprime.GroupsImprime(currentLine)
			params = strings.Split(params[len(params)-1], ",")
			l.OpQueue = []models.TokenComp{}
			l.NamesQueue = []string{}
			for i, str := range params {
				str = strings.TrimSpace(str)
				token := l.AnalyzeType(currentLine, lineIndex, str)
				if i != len(params)-1 {
					token = append(token, []string{",", helpers.DELIMITADOR}...)
					l.OpQueue = append(l.OpQueue, models.DELIM)
				}
				if len(token) > 0 {
					l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, token))
				}
			}

			l.AnalyzeOpQueue(currentLine, lineIndex)
			l.AnalyzeFuncQueue(currentLine, lineIndex)
			if !l.ExpectNoNone() {
				l.LogError(lineIndex, "N/A", "N/A", "One of the parameters introduced is not valid", currentLine)
			}
			l.GL.Printf("%+v Found 'Imprime' instruction [Line: %+v]", funcName, lineIndex)

			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				")", helpers.DELIMITADOR,
				";", helpers.DELIMITADOR,
			}))

			foundSomething = true
		}

		//Lee
		if l.R.RegexIO.MatchLee(currentLine, lineIndex) {
			if !l.R.RegexIO.MatchPC(currentLine, lineIndex) {
				l.LogError(lineIndex, len(currentLine)-1, ";", "Missing ';'", currentLine)
			}

			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				"Lee", helpers.PALABRARESERVADA,
				"(", helpers.DELIMITADOR,
			}))

			l.OpQueue = []models.TokenComp{}
			l.NamesQueue = []string{}
			params := l.R.RegexLee.GroupsLee(currentLine)
			l.AnalyzeParams(currentLine, lineIndex, params[0])

			l.AnalyzeFuncQueue(currentLine, lineIndex)

			if !l.ExpectIdent(currentLine, lineIndex) {
				l.LogError(lineIndex, "N/A", "N/A", "Expected <Ident> in parameters", currentLine)
			}

			l.GL.Printf("%+v Found 'Lee' instruction [Line: %+v]", funcName, lineIndex)

			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				")", helpers.DELIMITADOR,
				";", helpers.DELIMITADOR,
			}))

			foundSomething = true
		}

		//Cuando
		if l.R.RegexConditionCuando.StartsWithCuando(currentLine, lineIndex) {
			if len(l.BlockQueue) == 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to create a CUANDOBLOCK outside of a BLOCK", currentLine)
			}
			l.BlockQueue = append(l.BlockQueue, models.CUANDOBLOCK)

			//TODO: Get params

			l.GL.Printf("%+v Created a CUANDOBLOCK [Line: %+v]", funcName, lineIndex)

			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"cuando", helpers.PALABRARESERVADA}))

			foundSomething = true
		}

		//Si
		if l.R.RegexConditionSi.StartsWithSi(currentLine, lineIndex) {
			if len(l.BlockQueue) == 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to create a 'Si' condition outside of a BLOCK", currentLine)
			}

			l.R.RegexConditionSi.ValidateCondition(currentLine, lineIndex)
			l.GL.Printf("%+v Found 'Si' condition [Line: %+v]", funcName, lineIndex)

			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				"si", helpers.PALABRARESERVADA,
				"(", helpers.DELIMITADOR,
			}))

			l.OpQueue = []models.TokenComp{}
			l.NamesQueue = []string{}
			groups := helpers.GetGroupMatches(currentLine, helpers.SIREGEXP)
			params := groups[0]
			l.AnalyzeParams(currentLine, lineIndex, params)

			l.AnalyzeFuncQueue(currentLine, lineIndex)

			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				")", helpers.DELIMITADOR,
				"hacer", helpers.PALABRARESERVADA,
			}))

			foundSomething = true
		}

		//Sino
		if l.R.RegexConditionSi.StartsWithSino(currentLine, lineIndex) {
			if len(l.BlockQueue) == 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to create a 'Si' condition outside of a BLOCK", currentLine)
			}

			l.R.RegexConditionSi.ValidateCondition(currentLine, lineIndex)

			//TODO: Get Params

			l.GL.Printf("%+v Found 'Sino' condition [Line: %+v]", funcName, lineIndex)

			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"sino", helpers.PALABRARESERVADA}))

			foundSomething = true
		}

		//Switch: Sea
		if l.R.RegexConditionSwitch.StartsWithSea(currentLine, lineIndex) {
			if len(l.BlockQueue) == 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to create a 'Sea' instruction outside of a BLOCK", currentLine)
			}

			if l.BlockQueue[len(l.BlockQueue)-1] != models.INITBLOCK && l.BlockQueue[len(l.BlockQueue)-2] != models.CUANDOBLOCK {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to create a 'Sea' instruction outside of a CUANDOBLOCK", currentLine)
			}

			//TODO: Get Params

			l.GL.Printf("%+v Found 'Sea' instruction for CUANDOBLOCK [Line: %+v]", funcName, lineIndex)
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"sea", helpers.PALABRARESERVADA}))

			foundSomething = true
		}

		//Switch: Otro
		if l.R.RegexConditionSwitch.StartsWithOtro(currentLine, lineIndex) {
			if len(l.BlockQueue) == 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to create a 'Otro' instruction outside of a BLOCK", currentLine)
			}
			if l.BlockQueue[len(l.BlockQueue)-1] != models.INITBLOCK && l.BlockQueue[len(l.BlockQueue)-2] != models.CUANDOBLOCK {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to create a 'Otro' instruction outside of a CUANDOBLOCK", currentLine)
			}

			//TODO: Get Params

			l.GL.Printf("%+v Found 'Otro' instruction for CUANDOBLOCK [Line: %+v]", funcName, lineIndex)
			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"otro", helpers.PALABRARESERVADA}))

			foundSomething = true
		}

		//Regresa
		if l.R.RegexRegresa.MatchRegresa(currentLine, lineIndex) {
			if !l.R.RegexRegresa.MatchPC(currentLine, lineIndex) {
				l.LogError(lineIndex, len(currentLine)-1, ";", "Missing ';'", currentLine)
			}
			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				"Regresa", helpers.PALABRARESERVADA,
				"(", helpers.DELIMITADOR,
			}))

			l.OpQueue = []models.TokenComp{}
			l.NamesQueue = []string{}
			params := l.R.RegexRegresa.GroupsRegresa(currentLine)[0]
			l.AnalyzeParams(currentLine, lineIndex, params)

			l.AnalyzeFuncQueue(currentLine, lineIndex)
			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, []string{
				")", helpers.DELIMITADOR,
				";", helpers.DELIMITADOR,
			}))

			foundSomething = true
		}

		//Desde
		if l.R.RegexLoopDesde.StartsWithDesde(currentLine, lineIndex) {
			//TODO: Analyze
			l.GL.Printf("%+v Found 'Desde' instruction [Line: %+v]", funcName, lineIndex)

			l.LL.Println(helpers.IndentString(helpers.LEXINDENT, []string{"desde", helpers.PALABRARESERVADA}))

			foundSomething = true
		}

		//Interrumpe
		if l.R.RegexSystem.MatchInterrumpe(currentLine, lineIndex) {
			if !l.R.RegexSystem.MatchPC(currentLine, lineIndex) {
				l.LogError(lineIndex, len(currentLine)-1, ";", "Missing ';'", currentLine)
			}

			l.GL.Printf("%+v Found 'Interrumpe' instruction [Line: %+v]", funcName, lineIndex)

			foundSomething = true
		}

		//Limpia
		if l.R.RegexSystem.MatchLimpia(currentLine, lineIndex) {
			if !l.R.RegexSystem.MatchPC(currentLine, lineIndex) {
				l.LogError(lineIndex, len(currentLine)-1, ";", "Missing ';'", currentLine)
			}

			l.GL.Printf("%+v Found 'Limpia' instruction [Line: %+v]", funcName, lineIndex)

			foundSomething = true
		}

		//Asignación
		if l.R.RegexAsignacion.MatchAsignacion(currentLine, lineIndex) {
			if !l.R.RegexAsignacion.MatchPC(currentLine, lineIndex) {
				l.LogError(lineIndex, len(currentLine)-1, ";", "Missing ';'", currentLine)
			}

			currentLine = strings.TrimSpace(currentLine)
			data := strings.Split(currentLine, ":=")
			varToAssingData := data[0]
			varToAssingData = strings.TrimSpace(varToAssingData)
			assignToAnalyze := data[1]
			assignToAnalyze = strings.TrimSuffix(assignToAnalyze, ";")
			assignToAnalyze = strings.TrimSpace(assignToAnalyze)

			if l.R.RegexCustom.MatchCteLog(assignToAnalyze, lineIndex) {
				foundSomething = true
				l.GL.Printf("%+v Found 'Logica Assign' Operation [Line: %+v]", funcName, lineIndex)

				/*CHECK NO ASSIGN TO CONSTANT*/
				curToken := &models.Token{Type: models.LOGICO, Key: varToAssingData, Value: assignToAnalyze}
				if test := l.DoesTheTokenExistsInGlobalConstants(curToken); test {
					log.Printf("[ERR] Attempted to assign a value to a constant at [%+v][Line: %+v]", 0, lineIndex)
					l.GL.Printf("[ERR] Attempted to assign a value to a constant at [%+v][Line: %+v]", 0, lineIndex)
					//"# Linea | # Columna | Error | Descripcion | Linea del Error"
					l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "CONSTANT ASSIGN", "Attempted to assign a value to a constant", currentLine)
				}
				/*CHECK END*/

				/* CHECK IF ASSIGN CORRECT FOR VAR */
				if data := l.RetrieveGlobalVarIfExists(curToken); data != nil {
					if curToken.Type != data.Type {
						log.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
						l.GL.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
						//"# Linea | # Columna | Error | Descripcion | Linea del Error"
						l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "VARIABLE ASSIGN", "Attempted to assign a value of different type to a defined variable", currentLine)
					}
				}

				function := l.FindFunction(currentLine, lineIndex, l.Context)
				if function != nil {
					if data := l.RetrieveLocalVariableIfExists(curToken, function); data != nil {
						if curToken.Type != data.Type {
							log.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
							l.GL.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
							//"# Linea | # Columna | Error | Descripcion | Linea del Error"
							l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "VARIABLE ASSIGN", "Attempted to assign a value of different type to a defined variable", currentLine)
						}
					}
				}
				/*CHECK END*/
			} else if l.R.RegexCustom.MatchCteEnt(assignToAnalyze) {
				foundSomething = true
				l.GL.Printf("%+v Found 'Entera Assign' Operation [Line: %+v]", funcName, lineIndex)

				/*CHECK NO ASSIGN TO CONSTANT*/
				curToken := &models.Token{Type: models.ENTERO, Key: varToAssingData, Value: assignToAnalyze}
				if test := l.DoesTheTokenExistsInGlobalConstants(curToken); test {
					log.Printf("[ERR] Attempted to assign a value to a constant at [%+v][Line: %+v]", 0, lineIndex)
					l.GL.Printf("[ERR] Attempted to assign a value to a constant at [%+v][Line: %+v]", 0, lineIndex)
					//"# Linea | # Columna | Error | Descripcion | Linea del Error"
					l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "CONSTANT ASSIGN", "Attempted to assign a value to a constant", currentLine)
				}
				/*CHECK END*/
				/* CHECK IF ASSIGN CORRECT FOR VAR */
				if data := l.RetrieveGlobalVarIfExists(curToken); data != nil {
					if curToken.Type != data.Type {
						log.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
						l.GL.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
						//"# Linea | # Columna | Error | Descripcion | Linea del Error"
						l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "VARIABLE ASSIGN", "Attempted to assign a value of different type to a defined variable", currentLine)
					}
				}

				function := l.FindFunction(currentLine, lineIndex, l.Context)
				if function != nil {
					if data := l.RetrieveLocalVariableIfExists(curToken, function); data != nil {
						if curToken.Type != data.Type {
							log.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
							l.GL.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
							//"# Linea | # Columna | Error | Descripcion | Linea del Error"
							l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "VARIABLE ASSIGN", "Attempted to assign a value of different type to a defined variable", currentLine)
						}
					}
				}
				/*CHECK END*/
			} else if l.R.RegexCustom.MatchCteAlfa(assignToAnalyze) {
				foundSomething = true
				l.GL.Printf("%+v Found 'Alfabetica Assign' Operation [Line: %+v]", funcName, lineIndex)

				/*CHECK NO ASSIGN TO CONSTANT*/
				curToken := &models.Token{Type: models.ALFABETICO, Key: varToAssingData, Value: assignToAnalyze}
				if test := l.DoesTheTokenExistsInGlobalConstants(curToken); test {
					log.Printf("[ERR] Attempted to assign a value to a constant at [%+v][Line: %+v]", 0, lineIndex)
					l.GL.Printf("[ERR] Attempted to assign a value to a constant at [%+v][Line: %+v]", 0, lineIndex)
					//"# Linea | # Columna | Error | Descripcion | Linea del Error"
					l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "CONSTANT ASSIGN", "Attempted to assign a value to a constant", currentLine)
				}
				/*CHECK END*/
				/* CHECK IF ASSIGN CORRECT FOR VAR */
				if data := l.RetrieveGlobalVarIfExists(curToken); data != nil {
					if curToken.Type != data.Type {
						log.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
						l.GL.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
						//"# Linea | # Columna | Error | Descripcion | Linea del Error"
						l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "VARIABLE ASSIGN", "Attempted to assign a value of different type to a defined variable", currentLine)
					}
				}

				function := l.FindFunction(currentLine, lineIndex, l.Context)
				if function != nil {
					if data := l.RetrieveLocalVariableIfExists(curToken, function); data != nil {
						if curToken.Type != data.Type {
							log.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
							l.GL.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
							//"# Linea | # Columna | Error | Descripcion | Linea del Error"
							l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "VARIABLE ASSIGN", "Attempted to assign a value of different type to a defined variable", currentLine)
						}
					}
				}
				/*CHECK END*/
			} else if l.R.RegexCustom.MatchCteReal(assignToAnalyze) {
				foundSomething = true
				l.GL.Printf("%+v Found 'Real Assign' Operation [Line: %+v]", funcName, lineIndex)

				/*CHECK NO ASSIGN TO CONSTANT*/
				curToken := &models.Token{Type: models.REAL, Key: varToAssingData, Value: assignToAnalyze}
				if test := l.DoesTheTokenExistsInGlobalConstants(curToken); test {
					log.Printf("[ERR] Attempted to assign a value to a constant at [%+v][Line: %+v]", 0, lineIndex)
					l.GL.Printf("[ERR] Attempted to assign a value to a constant at [%+v][Line: %+v]", 0, lineIndex)
					//"# Linea | # Columna | Error | Descripcion | Linea del Error"
					l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "CONSTANT ASSIGN", "Attempted to assign a value to a constant", currentLine)
				}
				/*CHECK END*/
				/* CHECK IF ASSIGN CORRECT FOR VAR */
				if data := l.RetrieveGlobalVarIfExists(curToken); data != nil {
					if curToken.Type != data.Type {
						log.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
						l.GL.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
						//"# Linea | # Columna | Error | Descripcion | Linea del Error"
						l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "VARIABLE ASSIGN", "Attempted to assign a value of different type to a defined variable", currentLine)
					}
				}

				function := l.FindFunction(currentLine, lineIndex, l.Context)
				if function != nil {
					if data := l.RetrieveLocalVariableIfExists(curToken, function); data != nil {
						if curToken.Type != data.Type {
							log.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
							l.GL.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
							//"# Linea | # Columna | Error | Descripcion | Linea del Error"
							l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "VARIABLE ASSIGN", "Attempted to assign a value of different type to a defined variable", currentLine)
						}
					}
				}
				/*CHECK END*/
			} else {
				//TEMP CASE TO UNCHECKED EXPRESSION
				typeOfAssigment := l.GetOperationTypeFromAssignment(assignToAnalyze, lineIndex)

				/*CHECK NO ASSIGN TO CONSTANT*/
				if typeOfAssigment != models.INDEFINIDO {
					curToken := &models.Token{Type: typeOfAssigment, Key: varToAssingData, Value: assignToAnalyze}

					if test := l.DoesTheTokenExistsInGlobalConstants(curToken); test {
						log.Printf("[ERR] Attempted to assign a value to a constant at [%+v][Line: %+v]", 0, lineIndex)
						l.GL.Printf("[ERR] Attempted to assign a value to a constant at [%+v][Line: %+v]", 0, lineIndex)
						//"# Linea | # Columna | Error | Descripcion | Linea del Error"
						l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "CONSTANT ASSIGN", "Attempted to assign a value to a constant", currentLine)
					}
					/*CHECK END*/
					/* CHECK IF ASSIGN CORRECT FOR VAR */
					if data := l.RetrieveGlobalVarIfExists(curToken); data != nil {
						if curToken.Type != data.Type {
							log.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
							l.GL.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
							//"# Linea | # Columna | Error | Descripcion | Linea del Error"
							l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "VARIABLE ASSIGN", "Attempted to assign a value of different type to a defined variable", currentLine)
						}
					}

					function := l.FindFunction(currentLine, lineIndex, l.Context)
					if function != nil {
						if data := l.RetrieveLocalVariableIfExists(curToken, function); data != nil {
							if curToken.Type != data.Type {
								log.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
								l.GL.Printf("[ERR] Attempted to assign a %+v to a defined variable of type %+v at [%+v][Line: %+v]", curToken.Type, data.Type, 0, lineIndex)
								//"# Linea | # Columna | Error | Descripcion | Linea del Error"
								l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "VARIABLE ASSIGN", "Attempted to assign a value of different type to a defined variable", currentLine)
							}
						}
					}
				} else {
					log.Printf("[ERR] Attempted to assign an invalid expression to a defined variable at [%+v][Line: %+v]", 0, lineIndex)
					l.GL.Printf("[ERR] Attempted to assign an invalid expression to a defined variable at [%+v][Line: %+v]", 0, lineIndex)
					//"# Linea | # Columna | Error | Descripcion | Linea del Error"
					l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, 0, "VARIABLE ASSIGN", "Attempted to assign an invalid expression to a defined variable", currentLine)
				}

				/*CHECK END*/
				foundSomething = true

			}

			if !foundSomething {
				l.GL.Printf("%+v Found 'Unknown Assign [`%+v`]' instruction [Line: %+v] ", funcName, assignToAnalyze, lineIndex)
			}

			foundSomething = true
		}

		//Programa
		if l.R.RegexPrograma.StartsWithPrograma(currentLine, lineIndex) {
			l.GL.Println()

			if len(l.BlockQueue) > 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to create new Programa without finalizing the last Function or Procedure", currentLine)
				l.BlockQueue = []models.BlockType{}
			}
			l.BlockQueue = append(l.BlockQueue, models.PROGRAMBLOCK)

			l.Context = "Programa"
			l.HasMain = true
			foundSomething = true
		}

		//Fin de Programa
		if l.R.RegexFinPrograma.StartsWithFinPrograma(currentLine, lineIndex) {
			if len(l.BlockQueue) == 0 {
				l.LogError(lineIndex, "N/A", "N/A", "Attempted to END a PROGRAMBLOCK outside of a PROGRAMBLOCK", currentLine)
			}

			newArr, ok := helpers.RemoveFromQueue(l.BlockQueue, models.PROGRAMBLOCK)
			if ok {
				l.BlockQueue = newArr
			}

			l.Context = "Global"
			foundSomething = true
		}

		//Custom Functions
		if l.R.RegexCustomFunction.MatchCustomFunction(currentLine) && !foundSomething {
			l.GL.Printf("%+v Found 'Custom Function' instruction [Line: %+v]", funcName, lineIndex)
			currentLine = strings.TrimSuffix(currentLine, ";")
			groupsFunction := l.R.RegexCustomFunction.GroupsCustomFunction(currentLine)

			l.OpQueue = []models.TokenComp{models.ID, models.BRACK}
			l.NamesQueue = []string{groupsFunction[0]}
			token := []string{
				groupsFunction[0], helpers.IDENTIFICADOR,
				"(", helpers.DELIMITADOR,
			}
			if len(groupsFunction) > 0 {
				params := strings.Split(groupsFunction[1], ",")
				for i, param := range params {
					param = strings.TrimSpace(param)
					token = append(token, l.AnalyzeType("", 0, param)...)
					if i < len(params)-1 {
						l.OpQueue = append(l.OpQueue, models.DELIM)
						token = append(token, []string{",", helpers.DELIMITADOR}...)
					}
				}
			}
			l.OpQueue = append(l.OpQueue, models.BRACK)
			token = append(token, []string{
				")", helpers.DELIMITADOR,
				";", helpers.DELIMITADOR,
			}...)

			l.AnalyzeFuncQueue(currentLine, lineIndex)
			l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, token))
			foundSomething = true
		}

		//Logger
		l.RegisterBlockChange(LastBlockState, debug, funcName, lineIndex)

		/* Data Segregator */
		if l.CurrentBlockType == models.CONSTANTBLOCK {
			l.NextConstant(currentLine, lineIndex, debug)
		}

		if l.CurrentBlockType == models.VARIABLEBLOCK {
			l.NextVariable(currentLine, lineIndex, debug)
			// l.Print()
		}

		if l.CurrentBlockType == models.FUNCTIONPROTOBLOCK {
			l.NextFuncionProto(currentLine, lineIndex, debug)
		}

		if l.CurrentBlockType == models.PROCEDUREPROTOBLOCK {
			l.NextProcedureProto(currentLine, lineIndex, debug)
		}

		if l.ErrorsCount >= 20 {
			l.LogError(lineIndex, "N/A", "Compilation Stop", "Too many errors...", "")
			l.Status = -1
			return nil
		}

		if !foundSomething {
			switch l.CurrentBlockType {
			case models.NULLBLOCK:
				l.LogTest(lineIndex, "", "", "Didn't find anything", currentLine)
			}
		}
		lineIndex++
	}

	l.VerifyFunctions()
	// l.Print()
	return nil
}

//ValidateOperation ...
func (l *LexicalAnalyzer) isAValidOperation(assignStr string) bool {
	regCheck := regexp.MustCompile(`(\")?([a-zA-Z0-9.]+){1}(\")?((\*|\+|\/|\-){1}(\")?[a-zA-Z0-9.]+(\")?)*$`)
	return regCheck.MatchString(assignStr)
}

//GetOperationTypeFromAssignment ...
func (l *LexicalAnalyzer) GetOperationTypeFromAssignment(assignStr string, lineindex int64) models.TokenType {
	if l.isAValidOperation(assignStr) {
		curStr := assignStr
		test := regexp.MustCompile(`((\*){1}|(\+){1}|(\/){1}|(\-){1})`)
		operationParameters := test.Split(curStr, -1)
		log.Printf("TEST PARAMS> %+v", operationParameters)

		paramTypes := []models.TokenType{}
		for _, eachParam := range operationParameters {
			match := false
			if !match && l.R.RegexCustom.MatchCteAlfa(eachParam) {
				paramTypes = append(paramTypes, models.ALFABETICO)
				match = true
			}
			if !match && l.R.RegexCustom.MatchCteEnt(eachParam) {
				paramTypes = append(paramTypes, models.ENTERO)
				match = true
			}
			if !match && l.R.RegexCustom.MatchCteReal(eachParam) {
				paramTypes = append(paramTypes, models.REAL)
				match = true
			}
			if !match && l.R.RegexCustom.MatchCteLog(eachParam, lineindex) {
				paramTypes = append(paramTypes, models.LOGICO)
				match = true
			}
			if !match && l.R.RegexCustom.MatchIdent(eachParam) {

				//checar si variable o funcion o constante

				// paramTypes = append(paramTypes, models.LOGICO)
				match = true
			}
		}
		log.Printf("TEST TYPES> %+v", paramTypes)
	}

	return models.INDEFINIDO
}

//AnalyzeParams ...
func (l *LexicalAnalyzer) AnalyzeParams(currentLine string, lineIndex int64, params string) {
	condiciones := l.R.RegexOperatorLogico.V1.Split(params, -1)
	condicionadores := l.R.RegexOperatorLogico.GroupsOpLogico(params)
	for i, condicion := range condiciones {
		relaciones := l.R.RegexOperatorRelacional.V1.Split(condicion, -1)
		relacionadores := l.R.RegexOperatorRelacional.GroupsOpRelacional(condicion)
		for j, relacion := range relaciones {
			aritmeticos := l.R.RegexOperatorAritmetico.V1.Split(relacion, -1)
			aritmeticadores := l.R.RegexOperatorAritmetico.GroupsOpAritmetico(relacion)
			for k, aritmetico := range aritmeticos {
				aritmetico = strings.TrimPrefix(aritmetico, " ")
				aritmetico = strings.TrimSuffix(aritmetico, " ")
				token := []string{
					aritmetico,
				}
				token = l.AnalyzeType(currentLine, lineIndex, aritmetico)

				if len(token) > 0 {
					l.LL.Print(helpers.IndentStringInLines(helpers.LEXINDENT, 2, token))
				}
				if k < len(aritmeticadores) {
					l.OpQueue = append(l.OpQueue, models.OPARIT)
					l.LL.Print(helpers.IndentString(helpers.LEXINDENT, []string{aritmeticadores[k], helpers.OPERADORARITMETICO}))
				}
			}
			if j < len(relacionadores) {
				l.OpQueue = append(l.OpQueue, models.OPREL)
				l.LL.Print(helpers.IndentString(helpers.LEXINDENT, []string{relacionadores[j], helpers.OPERADORRELACIONAL}))
			}
		}
		if i < len(condicionadores) {
			l.OpQueue = append(l.OpQueue, models.OPLOG)
			l.LL.Print(helpers.IndentString(helpers.LEXINDENT, []string{condicionadores[i], helpers.OPERADORLOGICO}))
		}
	}
}

//AnalyzeType ...
func (l *LexicalAnalyzer) AnalyzeType(currentLine string, lineIndex int64, line string) []string {
	token := []string{line}
	if l.R.RegexCustom.MatchCteAlfa(line) {
		token = append(token, helpers.CONSTANTEALFABETICA)
		l.OpQueue = append(l.OpQueue, models.CTEALFA)
		l.NamesQueue = append(l.NamesQueue, line)
	} else if l.R.RegexFunction.MatchFunctionCallEnd(line) {
		token = l.AnalyzeType("", 0, line[:len(line)-1])
		token = append(token, []string{")", helpers.DELIMITADOR}...)
		l.OpQueue = append(l.OpQueue, models.BRACK)
	} else if l.R.RegexConstanteReal.MatchRealConstant(line) {
		token = append(token, helpers.CONSTANTEREAL)
		l.OpQueue = append(l.OpQueue, models.CTEREAL)
		l.NamesQueue = append(l.NamesQueue, line)
	} else if l.R.RegexConstanteEntera.MatchEnteraConstant(line) {
		token = append(token, helpers.CONSTANTEENTERA)
		l.OpQueue = append(l.OpQueue, models.CTEENT)
		l.NamesQueue = append(l.NamesQueue, line)
	} else if l.R.RegexFunction.MatchFunctionCall(line) {
		groups := strings.Split(line, "(")
		l.OpQueue = append(l.OpQueue, models.ID)
		l.OpQueue = append(l.OpQueue, models.BRACK)
		token = []string{
			groups[0], helpers.IDENTIFICADOR,
			"(", helpers.DELIMITADOR,
		}
		l.NamesQueue = append(l.NamesQueue, groups[0])
		if function := l.FindFunction(currentLine, lineIndex, groups[0]); function == nil {
			l.LogError(lineIndex, "N/A", "Undeclared function", "Could not find any reference for function: "+groups[0], currentLine)
		}
		if len(groups) > 1 {
			token = append(token, l.AnalyzeType("", 0, line[1:])...)
		}
		token = append(token, []string{
			")", helpers.DELIMITADOR,
		}...)
		l.OpQueue = append(l.OpQueue, models.BRACK)
	} else if l.R.RegexFunction.MatchFunctionCall2(line) {
		groups := strings.Split(line, "(")
		l.OpQueue = append(l.OpQueue, models.ID)
		l.OpQueue = append(l.OpQueue, models.BRACK)
		token = []string{
			groups[0], helpers.IDENTIFICADOR,
			"(", helpers.DELIMITADOR,
		}
		l.NamesQueue = append(l.NamesQueue, groups[0])
		if function := l.FindFunction(currentLine, lineIndex, groups[0]); function == nil {
			l.LogError(lineIndex, "N/A", "Undeclared function", "Could not find any reference for function: "+groups[0], currentLine)
		}
		if len(groups) > 1 {
			token = append(token, l.AnalyzeType("", 0, groups[1])...)
		}
	} else {
		groups := l.R.RegexVar.GroupsVar(line)
		token = []string{groups[0], helpers.IDENTIFICADOR}
		l.OpQueue = append(l.OpQueue, models.ID)
		l.NamesQueue = append(l.NamesQueue, groups[0])
		if lineIndex != 0 {
			l.FindSymbol(currentLine, lineIndex, groups[0])
		}
		if len(groups) > 1 {
			for _, group := range groups[1:] {
				if len(group) > 2 {
					token = append(token, []string{
						"[", helpers.DELIMITADOR,
						group[1 : len(group)-1], helpers.IDENTIFICADOR,
						"]", helpers.DELIMITADOR,
					}...)
					l.OpQueue = append(l.OpQueue, models.BRACK)
					l.OpQueue = append(l.OpQueue, models.ID)
					l.OpQueue = append(l.OpQueue, models.BRACK)
					l.NamesQueue = append(l.NamesQueue, group[1:len(group)-1])
					if lineIndex != 0 {
						l.FindSymbol(currentLine, lineIndex, group[1:len(group)-1])
					}
				}
			}
		}
	}

	return token
}

//AnalyzeOpQueue ...
func (l *LexicalAnalyzer) AnalyzeOpQueue(currentLine string, lineIndex int64) {
	noParentheses := 0
	last := l.OpQueue[0]
	token := ""
	for _, item := range l.OpQueue[1:] {
		switch item {
		case models.BRACK:
			token = "brackets"
			noParentheses++
			if last == models.PALRES {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.PALABRARESERVADA+" before "+token, currentLine)
			} else if last == models.DELIM {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.DELIMITADOR+" before "+token, currentLine)
			}
			break
		case models.CTEALFA, models.CTEENT, models.CTELOG, models.CTEREAL:
			token = helpers.CONSTANTE
			if last == models.PALRES {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.PALABRARESERVADA+" before "+token, currentLine)
			} else if last == models.ID {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.IDENTIFICADOR+" before "+token, currentLine)
			} else if last == models.CTEALFA || last == models.CTEENT || last == models.CTELOG || last == models.CTEREAL {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.CONSTANTE+" before "+token, currentLine)
			}
			break
		case models.DELIM:
			token = helpers.DELIMITADOR
			if last == models.OPARIT || last == models.OPASIG || last == models.OPLOG || last == models.OPREL {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.OPERADOR+" before "+token, currentLine)
			} else if last == models.DELIM {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.DELIMITADOR+" before "+token, currentLine)
			}
			break
		case models.ID:
			token = helpers.IDENTIFICADOR
			if last == models.CTEALFA || last == models.CTEENT || last == models.CTELOG || last == models.CTEREAL {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.CONSTANTE+" before "+token, currentLine)
			} else if last == models.PALRES {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.PALABRARESERVADA+" before "+token, currentLine)
			} else if last == models.ID {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.IDENTIFICADOR+" before "+token, currentLine)
			}
			break
		case models.OPARIT, models.OPASIG, models.OPLOG, models.OPREL:
			token = helpers.OPERADOR
			if last == models.DELIM {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.DELIMITADOR+" before "+token, currentLine)
			} else if last == models.OPARIT || last == models.OPASIG || last == models.OPLOG || last == models.OPREL {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.OPERADOR+" before "+token, currentLine)
			} else if last == models.PALRES {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.PALABRARESERVADA+" before "+token, currentLine)
			}
			break
		case models.PALRES:
			token = helpers.PALABRARESERVADA
			if last == models.PALRES {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.PALABRARESERVADA+" before "+token, currentLine)
			} else if last == models.ID {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.IDENTIFICADOR+" before "+token, currentLine)
			} else if last == models.CTEALFA || last == models.CTEENT || last == models.CTELOG || last == models.CTEREAL {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.CONSTANTE+" before "+token, currentLine)
			} else if last == models.BRACK {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected brackets before "+token, currentLine)
			} else if last == models.OPARIT || last == models.OPASIG || last == models.OPLOG || last == models.OPREL {
				l.LogError(lineIndex, "N/A", "UNEXPECTED", "Unexpected "+helpers.OPERADOR+" before "+token, currentLine)
			}
			break
		case models.NONE:
			l.LogError(lineIndex, "N/A", "NONE", "Couldn't find reference", currentLine)
			break
		}
		last = item
	}

	if noParentheses%2 != 0 {
		l.LogError(lineIndex, "N/A", "Brackets", "Missing brackets", currentLine)
	}
}

// OperationTokenType ...
func (l *LexicalAnalyzer) OperationTokenType(first, operator, second models.TokenComp, noVars int, currentLine string, lineIndex int64) models.Token {
	if operator == models.OPLOG || operator == models.OPREL {
		return models.Token{Type: models.LOGICO}
	}

	result := models.Token{}
	switch first {
	case models.ID:
		symbol := l.FindSymbol(currentLine, lineIndex, l.NamesQueue[noVars])
		if symbol != nil {
			result = models.Token{Type: symbol.Type}
		}
		break
	case models.CTEALFA, models.CTEENT, models.CTELOG, models.CTEREAL:
		result = models.Token{Type: models.ConstTypeToTokenType(first)}
		break
	}
	return result
}

//AnalyzeFuncQueue ...
func (l *LexicalAnalyzer) AnalyzeFuncQueue(currentLine string, lineIndex int64) {
	noVars := 0
	noBracks := 0
	var function *models.TokenFunc
	currentFunction := models.TokenFunc{}
	for i := 0; i < len(l.OpQueue); i++ {
		item := l.OpQueue[i]
		if function == nil && item != models.ID {
			if item == models.CTEALFA || item == models.CTEENT || item == models.CTELOG || item == models.CTEREAL {
				noVars++
			}
			continue
		}
		if function != nil && i+1 < len(l.OpQueue) {
			next := l.OpQueue[i+1]
			if next == models.OPARIT || next == models.OPASIG || next == models.OPLOG || next == models.OPREL {
				result := l.OperationTokenType(item, next, l.OpQueue[i+2], noVars, currentLine, lineIndex)
				if result.Type != "" {
					currentFunction.Params = append(currentFunction.Params, &result)
					i += 2
					continue
				}
			}
		}
		switch item {
		case models.BRACK:
			if function != nil {
				noBracks++
			}
			break
		case models.CTEALFA, models.CTEENT, models.CTELOG, models.CTEREAL:
			if function != nil {
				if noBracks < 2 {
					currentFunction.Params = append(currentFunction.Params, &models.Token{
						Type: models.ConstTypeToTokenType(item),
					})
				}
			}
			noVars++
			break
		case models.ID:
			if function == nil {
				function = l.FindFunction(currentLine, lineIndex, l.NamesQueue[noVars])
				if function != nil {
					currentFunction = *function
					currentFunction.Params = []*models.Token{}
				}
			} else {
				if noBracks < 2 {
					symbol := l.FindSymbol(currentLine, lineIndex, l.NamesQueue[noVars])
					if symbol != nil {
						currentFunction.Params = append(currentFunction.Params, symbol)
					}
				}
			}
			noVars++
			break
		}
	}

	if function != nil {
		l.CompareFunction(currentLine, lineIndex, function, &currentFunction, true)
	}
}

//FindSymbol Returns value for given key if found in symbol table
func (l *LexicalAnalyzer) FindSymbol(currentLine string, lineIndex int64, key string) *models.Token {
	if l.R.RegexReserved.IsReserved(key) {
		l.LogError(lineIndex, "N/A", "INVALID", fmt.Sprintf("Can not use %v (reserved) as name", key), currentLine)
		return nil
	}

	for _, symbol := range l.ConstantStorage {
		if symbol.Key == key {
			return symbol
		}
	}
	for _, symbol := range l.VariableStorage {
		if symbol.Key == key {
			return symbol
		}
	}

	if l.Context != "Global" {
		// fmt.Println("Finding " + key + " in " + l.Context)
		function := l.FindFunction(currentLine, lineIndex, l.Context)
		if function != nil {
			for _, symbol := range function.Params {
				if symbol.Key == key {
					return symbol
				}
			}
			for _, symbol := range function.Vars {
				if symbol.Key == key {
					return symbol
				}
			}
		}
	}
	l.LogError(lineIndex, "N/A", "Undeclared name", "Could not find any reference for name: "+key, currentLine)
	return nil
}

//FindFunction Returns value for given key if found in symbol table
func (l *LexicalAnalyzer) FindFunction(currentLine string, lineIndex int64, key string) *models.TokenFunc {
	for _, symbol := range l.FunctionStorage {
		if symbol.Key == key {
			return symbol
		}
	}
	// l.LogError(lineIndex, "N/A", "Undeclared name", "Could not find any reference for function name: "+key, currentLine)
	return nil
}

//CompareFunction Compares both given TokenFunc params
func (l *LexicalAnalyzer) CompareFunction(currentLine string, lineIndex int64, model, current *models.TokenFunc, isCall bool) {
	if isCall {
		model.Calls = append(model.Calls, &models.Line{CurrentLine: currentLine, LineIndex: lineIndex})
	} else {
		model.IsDefined = true
	}
	modelParams := []string{}
	modelSignature := "("
	for i, param := range model.Params {
		modelParams = append(modelParams, string(param.Type))
		modelSignature += string(param.Type)
		if i < len(model.Params)-1 {
			modelSignature += ","
		}
	}
	modelSignature += ")"
	match := true
	currentParams := []string{}
	currentSignature := "("
	for i, param := range current.Params {
		currentParams = append(currentParams, string(param.Type))
		currentSignature += string(param.Type)
		if i < len(current.Params)-1 {
			currentSignature += ","
		}
		if i < len(modelParams) && modelParams[i] != currentParams[i] {
			match = false
		}
	}
	currentSignature += ")"
	if len(modelParams) != len(currentParams) {
		match = false
	}
	if match {
		return
	}

	l.LogError(lineIndex, "N/A", "UNEXPECTED", "Mismatch in "+model.Key+" Want: "+modelSignature+" Have: "+currentSignature, currentLine)
}

//VerifyFunctions ...
func (l *LexicalAnalyzer) VerifyFunctions() {
	if !l.HasMain {
		l.LogError(0, "N/A", "UNDEFINED", "Could not find any main definition", "")
	}
	for _, function := range l.FunctionStorage {
		if len(function.Calls) > 0 {
			if !function.IsDefined {
				for _, call := range function.Calls {
					l.LogError(call.LineIndex, "N/A", "UNDEFINED", "Trying to use a function that only has prototype and is not defined", call.CurrentLine)
				}
			}
		} else {
			l.LogError(0, "N/A", "WARNING", fmt.Sprintf("Function %v was declared but never used", function.Key), "")
		}
	}
}

//LogError ...
//"# Linea | # Columna | Error | Descripcion | Linea del Error"
func (l *LexicalAnalyzer) LogError(lineIndex int64, columnIndex interface{}, err string, description string, currentLine string) {
	log.Printf("[ERR] %+v [Line: %+v]", description, lineIndex)
	l.GL.Printf("[ERR] %+v [Line: %+v]", description, lineIndex)
	l.EL.Printf("%+v\t|\t%+v\t|\t%+v\t|\t%+v\t|\t%+v", lineIndex, columnIndex, err, description, currentLine)
	l.ErrorsCount++
}

//LogErrorGeneral ...
//"# Linea | # Columna | Error | Descripcion | Linea del Error"
func (l *LexicalAnalyzer) LogErrorGeneral(lineIndex int64, columnIndex interface{}, err string, description string, currentLine string) {
	log.Printf("[ERR] %+v [Line: %+v] | '%+v'", description, lineIndex, currentLine)
	l.GL.Printf("[ERR] %+v [Line: %+v] | '%+v'", description, lineIndex, currentLine)
}

//LogTest ...
//"# Linea | # Columna | Error | Descripcion | Linea del Error"
func (l *LexicalAnalyzer) LogTest(lineIndex int64, columnIndex interface{}, err string, description string, currentLine string) {
	log.Printf("[ERR] %+v [Line: %+v] | '%+v'", description, lineIndex, currentLine)
	l.TEST.Printf("[ERR] %+v [Line: %+v] | '%+v'", description, lineIndex, currentLine)
}

//RegisterBlockChange ...
func (l *LexicalAnalyzer) RegisterBlockChange(LastBlockState models.BlockType, debug bool, funcName string, lineIndex int64) {
	if LastBlockState != l.CurrentBlockType {
		l.GL.Printf("%+v Switched to %+v [%+v]", funcName, l.CurrentBlockType, lineIndex)
		if debug {
			log.Printf("Switched to %+v [%+v]", l.CurrentBlockType, lineIndex)
		}
	}
}

//Print ...
func (l *LexicalAnalyzer) Print() {
	log.SetFlags(0)
	if len(l.ConstantStorage) > 0 {
		log.Print("Constants: ")
		for _, each := range l.ConstantStorage {
			pprnt.Print(*each)
		}
		log.Print("\n")
	} else {
		log.Println("Constants: 0")
	}

	if len(l.VariableStorage) > 0 {
		log.Print("Variables: ")
		for _, each := range l.VariableStorage {
			pprnt.Print(*each)
		}
		log.Print("\n")
	} else {
		log.Println("Variables: 0")
	}

	if len(l.FunctionStorage) > 0 {
		log.Print("Functions: ")
		for _, each := range l.FunctionStorage {
			pprnt.Print(*each)
		}
		log.Print("\n")
	} else {
		log.Println("Functions: 0")
	}

	log.SetFlags(log.LstdFlags)
}
