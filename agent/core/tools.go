package core

import (
	"github.com/clover0/issue-agent/core/functions"
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

func CommentingTools() []functions.Function {
	m := functions.FunctionsMap()

	return []functions.Function{
		m[functions.FuncOpenFile],
		m[functions.FuncPutFile],
		m[functions.FuncListFiles],
		m[functions.FuncSearchFiles],
		m[functions.FuncGetPullRequest],
		m[functions.FuncSubmitRevision],
	}
}
