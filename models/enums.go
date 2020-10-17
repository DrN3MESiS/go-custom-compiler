package models

//BlockType ...
type BlockType int

const (
	//DEFAULTBLOCK ...
	DEFAULTBLOCK BlockType = iota
	//CONSTANTBLOCK ...
	CONSTANTBLOCK
	//VARIABLEBLOCK ...
	VARIABLEBLOCK
	//FUNCTIONPROTOBLOCK ...
	FUNCTIONPROTOBLOCK
	//FUNCTIONBLOCK ...
	FUNCTIONBLOCK
	//PROCEDUREPROTOBLOCK ...
	PROCEDUREPROTOBLOCK
	//PROCEDUREBLOCK ...
	PROCEDUREBLOCK
)

//TokenType ...
type TokenType string

const (
	//FLOTANTE ...
	FLOTANTE TokenType = "FLOTANTE"
	//ENTERO ...
	ENTERO = "ENTERO"
	//CADENA ...
	CADENA = "CADENA"
	//ALFABETICO ...
	ALFABETICO = "ALFABETICO"
	//LOGICO ...
	LOGICO = "LOGICO"
	//REAL ...
	REAL = "REAL"
)