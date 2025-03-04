package core

import (
	"github.com/clover0/issue-agent/functions"
)

func PlanTools() []functions.Function {
	m := functions.FunctionsMap()

	return []functions.Function{
		m[functions.FuncOpenFile],
		m[functions.FuncListFiles],
		m[functions.FuncSearchFiles],
		m[functions.FuncGetPullRequest],
	}
}
