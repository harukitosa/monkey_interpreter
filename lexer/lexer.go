package lexer

import "github.com/harukitosa/monkey/token"

type Lexer struct {
	input        string // 入力されるmonkeyのコード
	position     int    // 入力における現在の位置
	readPosition int    // これから読み込む位置（現在の文字の次)
	ch           byte   // 現在検査中の文字
}

// 初期化
// input(コード)を受け取ってlexer構造体を返す
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 文字列を読み込む
func (l *Lexer) readChar() {
	// これから読み込む位置が入力された文字列より長いなら0
	// そうでないならば入力行からStringのreadPosition番目
	// の文字をとってきてchに代入
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	// 現在の位置を変更
	// 次に読み込む位置を一つ進める
	l.position = l.readPosition
	l.readPosition++
}

// lexerを受け取ってtoken構造体を返している
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// 空白行等をから読み込みして文字を読んでいく
	l.skipWhitespace()

	// 読み込んだ現在の文字について分岐
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			// 一足先に読み込んだ文字と一緒に考える
			// 次の文字も=ならばEQになる
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			// = だけならばtoken.ASSIGN
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		// 次の文字が=ならばNOT_EQになる
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			// リテラルは連結
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		// 0であれば終了コード
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// 文字であるならば
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
			// 数値であるならば
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			// 不正文字
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// 文字列であるならばisLetterの間だけ文字を読み込む
// そうして文字列の分だけ値を返す
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 空白行の続く限り読み込み
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

// 数字が続く限り文字を読み込む
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// これから読み込む文字を一足先に取得する関数
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// 文字列である
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

//　数字である
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
