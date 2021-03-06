package micro

import (
	"strconv"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
scheme code:

	(define (reify-name n)
		(string->symbol
			(string-append "_" "." (number->string n))
		)
	)
*/
func reifyName(n int) *ast.SExpr {
	return ast.NewSymbol("_" + strconv.Itoa(n))
}

// reifys expects a value and initially an empty substitution.
func reifyS(v *ast.SExpr) Substitutions {
	return reifys(v, nil)
}

/*
scheme code:

	(define (reify-s v s)
		(let
			((v (walk v s)))
			(cond
				(
					(var? v)
					(let
						((n (reify-name (length s))))
						(cons `(,v . ,n) s)
					)
				)
				(
					(pair? v)
					(reify-s (cdr v) (reify-s (car v) s))
				)
				(else s)
			)
		)
	)
*/
func reifys(v *ast.SExpr, s Substitutions) Substitutions {
	vv := v
	if v.IsVariable() {
		vv = walk(v.Atom.Var, s)
	}
	if vv.IsVariable() {
		n := reifyName(len(s))
		return append([]*Substitution{&Substitution{Var: vv.Atom.Var.Name, Value: n}}, s...)
	}
	if vv.IsPair() {
		car := vv.Car()
		cars := reifys(car, s)
		cdr := vv.Cdr()
		return reifys(cdr, cars)
	}
	return s
}

// ReifyVarFromState is a curried function that reifies the input variable for the given input state.
func ReifyVarFromState(v string) func(s *State) *ast.SExpr {
	return func(s *State) *ast.SExpr {
		vv := walkStar(&ast.SExpr{Atom: &ast.Atom{Var: &ast.Variable{Name: v}}}, s.Substitutions)
		r := reifyS(vv)
		return walkStar(vv, r)
	}
}

// Reify reifies the input variable for the given input states.
func Reify(v string, ss []*State) []*ast.SExpr {
	return deriveFmapRs(ReifyVarFromState(v), ss)
}
