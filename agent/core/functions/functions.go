package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/openai/openai-go"

	corestore "github.com/clover0/issue-agent/core/store"
	"github.com/clover0/issue-agent/logger"
)

func InitializeFunctions(
	noSubmit bool,
	repoService GitHubService,
	submitFilesService SubmitFilesService,
	submitRevisionService SubmitRevisionService,
	allowFunctions []string,
) {
	if allowFunction(allowFunctions, FuncOpenFile) {
		InitOpenFileFunction()
	}
	if allowFunction(allowFunctions, FuncListFiles) {
		InitListFilesFunction()
	}
	if allowFunction(allowFunctions, FuncPutFile) {
		InitPutFileFunction()
	}
	if allowFunction(allowFunctions, FuncModifyFile) {
		InitModifyFileFunction()
	}
	// TODO:
	if !noSubmit && allowFunction(allowFunctions, FuncSubmitFiles) {
		InitSubmitFilesGitHubFunction(submitFilesService)
	}
	if allowFunction(allowFunctions, FuncGetWebSearchResult) {
		InitGetWebSearchResult()
	}
	if allowFunction(allowFunctions, FuncGetWebPageFromURL) {
		InitFuncGetWebPageFromURLFunction()
	}
	if allowFunction(allowFunctions, FuncGetPullRequest) {
		InitGetPullRequestFunction(repoService)
	}
	if allowFunction(allowFunctions, FuncSearchFiles) {
		InitSearchFilesFunction()
	}
	if allowFunction(allowFunctions, FuncRemoveFile) {
		InitRemoveFileFunction()
	}
	if allowFunction(allowFunctions, FuncSwitchBranch) {
		InitSwitchBranchFunction()
	}
	if allowFunction(allowFunctions, FuncSubmitRevision) {
		InitSubmitRevisionFunction(submitRevisionService)
	}
	if allowFunction(allowFunctions, FuncGetIssue) {
		InitGetIssueFunction(repoService)
	}
}

func allowFunction(allowFunctions []string, name string) bool {
	return slices.Contains(allowFunctions, name)
}

type FuncName string

func (f FuncName) String() string {
	return string(f)
}

type Function struct {
	Name        FuncName
	Description string
	Func        any
	Parameters  map[string]interface{}
}

var functionsMap = map[string]Function{}

// TODO: no dependent on openai-go
func (f Function) ToFunctionCalling() openai.FunctionDefinitionParam {
	return openai.FunctionDefinitionParam{
		Name:        openai.F(f.Name.String()),
		Description: openai.F(f.Name.String()),
		Parameters:  openai.F(openai.FunctionParameters(f.Parameters)),
	}
}

func FunctionByName(name string) (Function, error) {
	if f, ok := functionsMap[name]; ok {
		return f, nil
	}

	return Function{}, fmt.Errorf("%s does not exist in functions", name)
}

// AllFunctions returns all functions
// WARNING: Call InitializeFunctions before calling this function
func AllFunctions() []Function {
	var fns []Function
	for _, f := range functionsMap {
		fns = append(fns, f)
	}
	return fns
}

func FunctionsMap() map[string]Function {
	return functionsMap
}

func marshalFuncArgs(args string, input any) error {
	return json.Unmarshal([]byte(args), &input)
}

const defaultSuccessReturning = "tool use succeeded."

func ExecFunction(l logger.Logger, store *corestore.Store, funcName FuncName, argsJson string) (string, error) {
	// TODO: make large switch statement smaller
	switch funcName {
	case FuncOpenFile:
		l.Info("functions: do %s\n", FuncOpenFile)
		input := OpenFileInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}
		file, err := OpenFile(input)
		if err != nil {
			return "", err
		}
		return file.Content, nil

	case FuncListFiles:
		l.Info("functions: do %s\n", FuncListFiles)
		input := ListFilesInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}
		files, err := ListFiles(input)
		if err != nil {
			return "", err
		}
		return strings.Join(files, "\n"), nil

	case FuncPutFile:
		l.Info("functions: do %s\n", FuncPutFile)
		input := PutFileInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}
		file, err := PutFile(input)
		if err != nil {
			return "", err
		}
		StoreFileAfterPutFile(store, file)
		return defaultSuccessReturning, nil

	case FuncModifyFile:
		l.Info("functions: do %s\n", FuncModifyFile)
		input := ModifyFileInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}
		file, err := ModifyFile(input)
		if err != nil {
			return "", err
		}
		StoreFileAfterModifyFile(store, file)
		return defaultSuccessReturning, nil

	case FuncSubmitFiles:
		l.Info("functions: do %s\n", FuncSubmitFiles)
		input := SubmitFilesInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}
		out, err := functionsMap[FuncSubmitFiles].Func.(SubmitFilesType)(input)
		if err != nil {
			return "", err
		}

		// NOTE: we would like to use any key, but for ease of implementation, we keep this as a simple implementation.
		SubmitFilesAfter(store, corestore.LastSubmissionKey, out)

		return fmt.Sprintf("%s\n%s\n",
			defaultSuccessReturning, out.Message), nil

	case FuncGetWebSearchResult:
		l.Info("functions: do %s\n", FuncGetWebSearchResult)
		input := GetWebSearchResultInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}

		r, err := GetWebSearchResult(input)
		if err != nil {
			return "", err
		}
		return r, nil

	case FuncGetWebPageFromURL:
		l.Info("functions: do %s\n", FuncGetWebPageFromURL)
		input := GetWebPageFromURLInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}

		r, err := GetWebPageFromURL(input)
		if err != nil {
			return "", err
		}
		return r, nil

	case FuncGetPullRequest:
		l.Info("functions: do %s\n", FuncGetPullRequest)
		input := GetPullRequestInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}
		fn, ok := functionsMap[FuncGetPullRequest].Func.(GetPullRequestType)
		if !ok {
			return "", fmt.Errorf("cat not call %s function", FuncGetPullRequest)
		}
		r, err := fn(input)
		if err != nil {
			return "", err
		}
		return r.ToLLMString(), nil

	case FuncSearchFiles:
		l.Info("functions: do %s\n", FuncSearchFiles)
		input := SearchFilesInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}
		r, err := SearchFiles(input)
		if err != nil {
			return "", err
		}
		return strings.Join(r, "\n"), nil

	case FuncRemoveFile:
		l.Info("functions: do %s\n", FuncRemoveFile)
		input := RemoveFileInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}
		err := RemoveFile(input)
		if err != nil {
			return "", err
		}
		return defaultSuccessReturning, nil

	case FuncSwitchBranch:
		l.Info("functions: do %s\n", FuncSwitchBranch)
		input := SwitchBranchInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}
		r, err := SwitchBranch(input)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s\n%s\n", defaultSuccessReturning, r), nil

	case FuncSubmitRevision:
		l.Info("functions: do %s\n", FuncSubmitRevision)
		input := SubmitRevisionInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}
		out, err := functionsMap[FuncSubmitRevision].Func.(SubmitRevisionType)(input)
		if err != nil {
			return "", err
		}
		return out.Message, nil

	case FuncGetIssue:
		l.Info("functions: do %s\n", FuncGetIssue)
		input := GetIssueInput{}
		if err := marshalFuncArgs(argsJson, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal args: %w", err)
		}
		out, err := functionsMap[FuncGetIssue].Func.(GetIssueType)(input)
		if err != nil {
			return "", err
		}
		return out.ToLLMString(), nil
	}

	return "", errors.New("function not found")
}
