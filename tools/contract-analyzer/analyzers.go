/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright 2019-2022 Dapper Labs, Inc.
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

package main

import (
	"fmt"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/sema"
	"github.com/onflow/cadence/tools/analysis"
)

var referenceToOptionalElements = []ast.Element{
	(*ast.ReferenceExpression)(nil),
}

var referenceToOptionalAnalyzer = &analysis.Analyzer{
	Requires: []*analysis.Analyzer{
		analysis.InspectorAnalyzer,
	},
	Run: func(pass *analysis.Pass) interface{} {
		inspector := pass.ResultOf[analysis.InspectorAnalyzer].(*ast.Inspector)

		location := pass.Program.Location
		elaboration := pass.Program.Elaboration
		report := pass.Report

		inspector.Preorder(
			referenceToOptionalElements,
			func(element ast.Element) {
				referenceExpression, ok := element.(*ast.ReferenceExpression)
				if !ok {
					return
				}

				indexExpression, ok := referenceExpression.Expression.(*ast.IndexExpression)
				if !ok {
					return
				}

				indexedType := elaboration.IndexExpressionIndexedTypes[indexExpression]
				resultType := indexedType.ElementType(false)
				_, ok = resultType.(*sema.OptionalType)
				if !ok {
					return
				}

				report(
					analysis.Diagnostic{
						Location: location,
						Range:    ast.NewRangeFromPositioned(indexExpression),
						Message:  "reference to optional",
					},
				)
			},
		)

		return nil
	},
}

func init() {

	registerAnalyzer(
		"reference-to-optional",
		referenceToOptionalAnalyzer,
	)
}

var deprecatedKeyFunctionsElements = []ast.Element{
	(*ast.InvocationExpression)(nil),
}

var deprecatedKeyFunctionsAnalyzer = &analysis.Analyzer{
	Requires: []*analysis.Analyzer{
		analysis.InspectorAnalyzer,
	},
	Run: func(pass *analysis.Pass) interface{} {
		inspector := pass.ResultOf[analysis.InspectorAnalyzer].(*ast.Inspector)

		location := pass.Program.Location
		elaboration := pass.Program.Elaboration
		report := pass.Report

		inspector.Preorder(
			deprecatedKeyFunctionsElements,
			func(element ast.Element) {
				invocationExpression, ok := element.(*ast.InvocationExpression)
				if !ok {
					return
				}

				memberExpression, ok := invocationExpression.InvokedExpression.(*ast.MemberExpression)
				if !ok {
					return
				}

				memberInfo := elaboration.MemberExpressionMemberInfos[memberExpression]
				member := memberInfo.Member
				if member == nil {
					return
				}

				if member.ContainerType != sema.AuthAccountType {
					return
				}

				var details string
				switch member.Identifier.Identifier {
				case sema.AuthAccountAddPublicKeyField:
					details = fmt.Sprintf(
						"replace '%s' with '%s'",
						sema.AuthAccountAddPublicKeyField,
						"keys.add",
					)
				case sema.AuthAccountRemovePublicKeyField:
					details = fmt.Sprintf(
						"replace '%s' with '%s'",
						sema.AuthAccountRemovePublicKeyField,
						"keys.revoke",
					)
				default:
					return
				}

				report(
					analysis.Diagnostic{
						Location: location,
						Range:    ast.NewRangeFromPositioned(element),
						Message:  fmt.Sprintf("use of deprecated key management API: %s", details),
					},
				)
			},
		)

		return nil
	},
}

func init() {
	registerAnalyzer(
		"deprecated-key-functions",
		deprecatedKeyFunctionsAnalyzer,
	)
}
