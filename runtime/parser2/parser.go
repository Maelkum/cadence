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

package parser2

import (
	"errors"
	"fmt"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/parser2/lexer"
)

type parser struct {
	tokens  []lexer.Token
	current lexer.Token
	pos     int
	errors  []error
	atEnd   bool
}

func Parse(input string) (ast.Expression, []error) {
	tokens := lexer.Lex(input)
	p := &parser{
		tokens:  tokens,
		current: tokens[0],
	}

	expr := parseExpression(p, 0)

	if !p.current.Is(lexer.TokenEOF) {
		p.report(fmt.Errorf("unexpected token: %v", p.current))
	}

	return expr, p.errors
}

func (p *parser) report(err error) {
	p.errors = append(p.errors, err)
}

func (p *parser) next() {
	p.pos++
	p.atEnd = p.pos >= len(p.tokens)
	if p.atEnd {
		p.report(errors.New("unexpected end of expression"))
	}
	p.current = p.tokens[p.pos]
}

func (p *parser) skipZeroOrOne(tokenType lexer.TokenType) {
	for p.current.Type == tokenType {
		p.next()
	}
}
