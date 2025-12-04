package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dream-horizon-org/datagen/codegen"
)

//go:generate goyacc -l -o grammar.go grammar.y

type stateFn func(*lex) (stateFn, int)

type MetadataEntry int

const (
	Count MetadataEntry = iota
	Tags
	MetadataEof
)

type Section int

const (
	Fields Section = iota
	Misc
	Metadata
	Gens
	GenFn
	Calls
	Serialiser
	None
)

const (
	COMMENT_MARKER = "//"
	KEYWORD_MODEL  = "model"
	GEN_FN_KEYWORD = "func"
)

const eof = -1

type lex struct {
	input         string
	curPos        int
	width         int
	parsed        *codegen.DatagenParsed
	fn            stateFn
	lval          *yySymType
	curSection    Section
	metadataEntry MetadataEntry
	err           error
}

func (l *lex) error(msg string, args ...any) (stateFn, int) {
	l.err = fmt.Errorf(msg, args...)
	return nil, -1
}

// Error implements yyLexer.
func (l *lex) Error(s string) {
	// Preserve earlier, more specific errors set via l.error
	if l.err == nil {
		l.err = errors.New(s)
	}
}

// consumeString reads a non-empty token until space or any brace/paren.
func (l *lex) consumeString() string {
	l.ditchSpacesAndComments()
	oldPos := l.curPos
Loop:
	for {
		switch b := l.nextByte(); {
		case unicode.IsSpace(b), b == '{', b == '(', b == '}', b == ')':
			l.backup()
			break Loop
		case b == eof:
			break Loop
		}
	}

	val := l.input[oldPos:l.curPos]
	return val
}

// consumeBodyTillRBrace reads a `{`-balanced body, stopping just before the
// matching right brace of the *current* `{` (nesting supported).
func (l *lex) consumeBodyTillRBrace() (string, error) {
	l.ditchSpacesAndComments()
	nesting := 1
	start := l.curPos
	for {
		switch b := l.nextByte(); b {
		case '{':
			nesting++
		case '}':
			nesting--
			if nesting == 0 {
				l.backup()
				return l.input[start:l.curPos], nil
			}
		case eof:
			return "", errors.New("incomplete body")
		}
	}
}

func lexModelName(l *lex) (stateFn, int) {
	val := l.consumeString()
	if val == "" {
		l.curSection = None
		return l.error("expected valid model name, got '%s'", val)
	}
	l.parsed.ModelName = val
	return lexLBrace, MODEL_NAME
}

func lexBody(l *lex) (stateFn, int) {
	val := l.consumeString()
	if val == "" {
		l.curSection = None
		return lexRBrace(l)
	}

	switch val {
	case "fields":
		l.curSection = Fields
		return lexLBrace, FIELDS
	case "misc":
		l.curSection = Misc
		return lexLBrace, MISC
	case "metadata":
		l.curSection = Metadata
		return lexLBrace, METADATA
	case "gens":
		l.curSection = Gens
		return lexLBrace, GEN_FNS
	case "calls":
		l.curSection = Calls
		return lexLBrace, CALLS
	case "serialiser":
		l.curSection = Serialiser
		return lexLBrace, SERIALISER_FUNC
	default:
		return l.error("expected section header: one of 'fields', 'misc', 'metadata', 'gens', 'calls', 'serialiser'; got '%s'", val)
	}
}

func lexFieldsBody(l *lex) (stateFn, int) {
	val, err := l.consumeBodyTillRBrace()
	if err != nil {
		return l.error("invalid fields body %s", err)
	}
	l.lval.str = val
	return lexRBrace, FIELDS_BODY
}

func lexMiscBody(l *lex) (stateFn, int) {
	val, err := l.consumeBodyTillRBrace()
	if err != nil {
		return l.error("invalid misc body %s", err)
	}
	l.lval.str = val
	return lexRBrace, MISC_BODY
}

func lexTagsBody(l *lex) (stateFn, int) {
	val, err := l.consumeBodyTillRBrace()
	if err != nil {
		return l.error("invalid tags body %s", err)
	}
	l.lval.str = val
	return lexRBrace, TAGS_BODY
}

func lexMetadataCount(l *lex) (stateFn, int) {
	val := l.consumeString()
	valInt, err := strconv.Atoi(val)
	if err != nil {
		return l.error("invalid count %s", err)
	}
	l.lval.count = valInt
	return lexMetadataBody, COUNT_INT
}

func lexMetadataColon(l *lex) (stateFn, int) {
	val := l.consumeString()
	if val != ":" {
		return l.error("expected colon (:)")
	}

	if l.metadataEntry == Count {
		return lexMetadataCount, COLON
	}

	if l.metadataEntry == Tags {
		return lexLBrace, COLON
	}

	return l.error("invalid metadata field")
}

func lexMetadataBody(l *lex) (stateFn, int) {
	val := l.consumeString()
	if strings.HasSuffix(val, ":") {
		l.prevByte()
		val = strings.TrimSuffix(val, ":")
	}
	if val == "count" {
		l.metadataEntry = Count
		return lexMetadataColon, COUNT
	}

	if val == "tags" {
		l.metadataEntry = Tags
		return lexMetadataColon, TAGS
	}

	if val != "" {
		return l.error("invalid metadata field")
	}

	// setting the metadataEntry to MetadataEof
	l.metadataEntry = MetadataEof
	return lexRBrace(l)
}

func lexCallsBody(l *lex) (stateFn, int) {
	val, err := l.consumeBodyTillRBrace()
	if err != nil {
		return l.error("invalid calls body %s", err)
	}
	l.lval.str = val
	return lexRBrace, CALLS_BODY
}

func lexFuncBody(l *lex) (stateFn, int) {
	val, err := l.consumeBodyTillRBrace()
	if err != nil {
		return l.error("invalid fn body %s", err)
	}
	l.lval.str = val
	return lexRBrace, FN_BODY
}

func lexGenFns(l *lex) (stateFn, int) {
	val := l.consumeString()
	if val == "" {
		// we are at an end of the gen section
		return lexRBrace(l)
	}

	if val != GEN_FN_KEYWORD {
		return l.error("invalid gens body, expected func, got '%s'", val)
	}

	return lexGenFnName, FN
}

func lexGenFnName(l *lex) (stateFn, int) {
	val := l.consumeString()
	if val == "" {
		return l.error("expected gen fn name, got '%s'", val)
	}
	l.lval.str = val
	return lexLParenthesis, FN_NAME
}

func lexGenFnArgs(l *lex) (stateFn, int) {
	l.ditchSpacesAndComments()
	oldPos := l.curPos
	for {
		b := l.nextByte()
		if b == ')' {
			l.backup()
			break
		}
		if b == eof {
			return l.error("incomplete body")
		}
	}

	val := l.input[oldPos:l.curPos]
	l.lval.str = val
	return lexRightRoundBrace, FN_ARGS
}

func lexRightRoundBrace(l *lex) (stateFn, int) {
	l.ditchSpacesAndComments()
	if b := l.nextByte(); b != ')' {
		l.backup()
		return l.error("expected ')', got '%s'", string(b))
	}
	l.curSection = GenFn
	return lexLBrace, R_PARENTHESIS
}

func lexLParenthesis(l *lex) (stateFn, int) {
	l.ditchSpacesAndComments()
	if b := l.nextByte(); b != '(' {
		l.backup()
		return l.error("expected '(', got '%s'", string(b))
	}
	return lexGenFnArgs, L_PARENTHESIS
}

func lexRBrace(l *lex) (stateFn, int) {
	l.ditchSpacesAndComments()
	if b := l.nextByte(); b != '}' {
		l.backup()
		return l.error("expected '}', got '%s'", string(b))
	}

	// we are done parsing the body
	if l.curSection == None {
		return nil, R_BRACE
	}

	// if we are in gens section, keep parsing gen funcs
	if l.curSection == GenFn {
		l.curSection = Gens
		return lexGenFns, R_BRACE
	}

	// if we are in metadata section, stay there
	if l.curSection == Metadata && l.metadataEntry != MetadataEof {
		l.curSection = Metadata
		return lexMetadataBody, R_BRACE
	}

	// if we are in the body, keep parsing sections
	return lexBody, R_BRACE
}

func lexLBrace(l *lex) (stateFn, int) {
	l.ditchSpacesAndComments()
	if b := l.nextByte(); b != '{' {
		l.backup()
		return l.error("expected '{', got '%s'", string(b))
	}

	// we are in the model line
	if l.curSection == None {
		return lexBody, L_BRACE
	}

	if l.curSection == Fields {
		return lexFieldsBody, L_BRACE
	}

	if l.curSection == Misc {
		return lexMiscBody, L_BRACE
	}

	if l.curSection == Metadata && l.metadataEntry == Tags {
		return lexTagsBody, L_BRACE
	}

	if l.curSection == Metadata {
		return lexMetadataBody, L_BRACE
	}

	if l.curSection == Gens {
		return lexGenFns, L_BRACE
	}

	if l.curSection == GenFn {
		return lexFuncBody, L_BRACE
	}

	if l.curSection == Calls {
		return lexCallsBody, L_BRACE
	}

	if l.curSection == Serialiser {
		return lexFuncBody, L_BRACE
	}

	return lexRBrace, MODEL_NAME
}

func lexModel(l *lex) (stateFn, int) {
	val := l.consumeString()
	if val != KEYWORD_MODEL {
		return l.error("expected 'model', got '%s'", val)
	}
	return lexModelName, MODEL
}

func (l *lex) parse_fields(s string) *ast.FieldList {
	fieldList, err := parseFieldList(s, parseWrappedExpr)
	if err != nil {
		l.error("could not parse field list: %s", err)
	}
	return fieldList
}

func (l *lex) parse_misc(s string) string {
	return s
}

func (l *lex) parse_tags(s string) map[string]string {
	tags, err := parseTags(s, parseWrappedExpr)
	if err != nil {
		l.error("could not parse tags: %s", err)
	}
	return tags
}

func (l *lex) parse_calls(s string) []*ast.CallExpr {
	calls, err := parseCallList(s, parseWrappedExpr)
	if err != nil {
		l.error("could not parse calls: %s", err)
	}
	return calls
}

func (l *lex) add_gen_fn(name, args, body string) {
	argsList, err := parseParamList(args, parseWrappedExpr)
	if err != nil {
		l.error("could not parse args list: %s", err)
	}

	funcBody, err := parseFunctionBlock(body, parseWrappedExpr)
	if err != nil {
		l.error("could not parse gen func block: %s", err)
	}

	l.parsed.GenFuns = append(l.parsed.GenFuns, &codegen.GenFn{
		Name:  name,
		Calls: argsList,
		Body:  funcBody,
	})
}

func (l *lex) add_serialiser_fn(body string) *codegen.SerialiserFunc {
	funcBody, err := parseFunctionBlock(body, parseWrappedExpr)
	if err != nil {
		l.error("could not parse serialiser func block: %s", err)
	}

	return &codegen.SerialiserFunc{
		Body: funcBody,
	}
}

// Lex implements yyLexer
func (l *lex) Lex(lval *yySymType) int {
	if l.err != nil {
		return -1
	}
	if l.fn == nil {
		return 0
	}

	l.lval = lval
	l.ditchSpacesAndComments()
	nextFn, val := l.fn(l)
	l.fn = nextFn
	return val
}

func (l *lex) ditchSpacesAndComments() {
	l.ditchSpaces()
	l.ditchComments()
}

func (l *lex) ditchSpaces() {
	for {
		b := l.nextByte()
		if !unicode.IsSpace(b) {
			l.backup()
			return
		}
	}
}

func (l *lex) ditchComments() {
	if v := l.peekn(len(COMMENT_MARKER)); v != COMMENT_MARKER {
		return
	}

	// consume COMMENT_MARKER
	l.seekn(len(COMMENT_MARKER))

	// consume the comment line now
	for {
		if v := l.nextByte(); v == '\n' {
			l.ditchSpacesAndComments()
			return
		}
	}
}

func (l *lex) seekn(n int) {
	l.curPos += n
}

func (l *lex) peekn(n int) string {
	if l.curPos+n >= len(l.input) {
		return ""
	}
	return l.input[l.curPos : l.curPos+n]
}

func (l *lex) backup() {
	l.curPos -= l.width
}

func (l *lex) nextByte() rune {
	if l.curPos >= len(l.input) {
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.curPos:])
	if r == utf8.RuneError && w == 1 {
		l.Error(fmt.Sprintf("invalid UTF-8 sequence while decoding next rune at pos %d", l.curPos))
	}
	l.curPos += w
	l.width = w
	return r
}

func (l *lex) prevByte() {
	if l.curPos == 0 {
		return
	}
	r, w := utf8.DecodeLastRuneInString(l.input[:l.curPos])
	if r == utf8.RuneError && w == 1 {
		l.Error(fmt.Sprintf("invalid UTF-8 sequence while decoding previous rune at pos %d", l.curPos))
	}
	l.curPos -= w
	l.width -= w
}

// ---- Entry point ----
func Parse(src []byte, path string) (*codegen.DatagenParsed, error) {
	l := lex{
		input:      string(src),
		fn:         lexModel,
		curSection: None,
		parsed:     &codegen.DatagenParsed{Filepath: path},
	}
	yyParse(&l)
	return l.parsed, l.err
}
