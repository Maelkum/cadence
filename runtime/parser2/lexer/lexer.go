/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lexer

import (
	"strings"
	"unicode/utf8"
)

type Token struct {
	Type  TokenType
	Value interface{}
}

func (t Token) Is(ty TokenType) bool {
	return t.Type == ty
}

type lexer struct {
	input  string
	start  int // start position of this item.
	pos    int // current position in the input.
	width  int
	tokens []Token
}

func Lex(input string) []Token {
	l := &lexer{input: input}
	l.run(rootState)
	return l.tokens
}

func (l *lexer) run(state stateFn) {
	for state != nil {
		state = state(l)
	}
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += w
	return r
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) word() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) acceptOne(r rune) bool {
	if l.next() == r {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptAny(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptZeroOrMore(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) acceptOneOrMore(valid string) bool {
	if !l.acceptAny(valid) {
		return false
	}
	l.acceptZeroOrMore(valid)
	return true
}

// emitValue passes an item back to the client.
func (l *lexer) emit(ty TokenType, val interface{}) {
	token := Token{ty, val}
	l.tokens = append(l.tokens, token)
	l.start = l.pos
}

func (l *lexer) emitType(ty TokenType) {
	l.emit(ty, nil)
}

func (l *lexer) emitValue(ty TokenType) {
	l.emit(ty, l.word())
}

func (l *lexer) emitError(err error) {
	l.emit(TokenError, err)
}

func (l *lexer) scanNumber() {
	// lookahead is already lexed.
	// parse more, if any
	l.acceptZeroOrMore("0123456789")
}

func (l *lexer) scanSpace() {
	// lookahead is already lexed.
	// parse more, if any
	l.acceptZeroOrMore(" \t")
}
