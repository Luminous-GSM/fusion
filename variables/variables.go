package variables

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/utils"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const globalScope = "----"

var (
	_conce          sync.Once
	variableService *VariableService
)

type VariableService struct {
	scopedVariables   map[string]map[string]string
	variableReplacers map[string]*strings.Replacer
	muVariables       sync.Mutex
	muReplacers       sync.Mutex
}

func Instance() *VariableService {
	_conce.Do(func() {
		variableService = &VariableService{
			scopedVariables:   make(map[string]map[string]string, 1),
			variableReplacers: make(map[string]*strings.Replacer, 1),
		}
		variableService.scopedVariables[globalScope] = make(map[string]string)
	})

	return variableService
}

func (vs *VariableService) RefreshAllVariables() {
	vs.RefreshGeneratedVariables()
	vs.RefreshConfigVariables()
}

func (vs *VariableService) refreshGlobalReplacer() {
	vs.muReplacers.Lock()
	defer vs.muReplacers.Unlock()

	oldNewVariableList := make([]string, 0)
	for key, val := range vs.scopedVariables[globalScope] {
		oldNewVariableList = append(oldNewVariableList, key, val)
	}
	vs.variableReplacers[globalScope] = strings.NewReplacer(oldNewVariableList...)
}

func (vs *VariableService) RefreshConfigVariables() {
	vs.muVariables.Lock()
	defer vs.muVariables.Unlock()

	defer vs.refreshGlobalReplacer()
	scopedMap := vs.scopedVariables[globalScope]

	conf := config.Get()
	scopedMap["{{fusion.directory.root}}"] = conf.RootDirectory
	scopedMap["{{fusion.directory.data}}"] = conf.DataDirectory
	scopedMap["{{fusion.directory.logs}}"] = conf.LogDirectory
	scopedMap["{{fusion.directory.certs}}"] = conf.CertsDirectory

	for _, role := range conf.SystemRoles {
		scopedMap[fmt.Sprintf("{{fusion.role.%v.password}}", role.Username)] = role.Password
	}
	vs.scopedVariables[globalScope] = scopedMap
}

func (vs *VariableService) RefreshGeneratedVariables() {
	vs.muVariables.Lock()
	defer vs.muVariables.Unlock()

	defer vs.refreshGlobalReplacer()
	scopedMap := vs.scopedVariables[globalScope]

	id, err := utils.GenerateSecureRandomString(12, false)
	if err != nil {
		zap.S().Errorw("variables: could not generated id. Using default id", "error", err, "id", "wdf356trqwd5")
		id = "wdf356trqwd5"
	}
	scopedMap["{{fusion.generated.id}}"] = id

	name := petname.Generate(2, " ")
	name = cases.Title(language.English).String(name)
	scopedMap["{{fusion.generated.name}}"] = name

	password, err := utils.GenerateSecureRandomString(16, false)
	if err != nil {
		zap.S().Errorw("variables: could not generated password. Using default id", "error", err, "password", "sdJfw4*7sF3s9@ds")
		password = "sdjfw4*7sf3s9(ds"
	}
	scopedMap["{{fusion.generated.password}}"] = password

	hashedPassword, _ := utils.HashPasswordBasedOnArgon2(password, true)
	scopedMap["{{fusion.generated.password_hashed}}"] = hashedPassword

	vs.scopedVariables[globalScope] = scopedMap
}

func (vs *VariableService) ReplaceGlobalVariablesInString(line string) string {
	return vs.variableReplacers[globalScope].Replace(line)
}

func (vs *VariableService) ReplaceScopedVariablesInString(scope, line string) string {
	return vs.variableReplacers[scope].Replace(line)
}

func (vs *VariableService) ReplaceGlobalAndScopedVariablesInString(scope, line string, failOnLeftOverVariables bool) (string, error) {
	val := vs.ReplaceGlobalVariablesInString(line)
	val = vs.ReplaceScopedVariablesInString(scope, val)
	if !vs.IsAllVariablesReplaced(val) {
		zap.S().Errorw("found unrecognised fusion variable", "line", val)
		return "", errors.New("found unrecognised fusion variable. Line : " + val)
	}
	return val, nil
}

func (vs *VariableService) AddVariableToScopedReplacer(scope, name, value string) {
	vs.muVariables.Lock()
	defer vs.muVariables.Unlock()

	scopedMap := vs.scopedVariables[scope]
	if scopedMap == nil {
		scopedMap = make(map[string]string)
	}
	scopedMap[fmt.Sprintf("{{%v}}", name)] = value
	vs.scopedVariables[scope] = scopedMap

	oldNewVariableList := make([]string, 0)
	for key, val := range vs.scopedVariables[scope] {
		oldNewVariableList = append(oldNewVariableList, key, val)
	}
	vs.variableReplacers[scope] = strings.NewReplacer(oldNewVariableList...)
}

func (vs *VariableService) RemoveScopedReplacer(scope string) {
	vs.muReplacers.Lock()
	defer vs.muReplacers.Unlock()

	delete(vs.variableReplacers, scope)
}

func (vs *VariableService) RemoveScopedVariables(scope string) {
	vs.muVariables.Lock()
	defer vs.muVariables.Unlock()

	delete(vs.scopedVariables, scope)
}

func (vs *VariableService) RemoveScopeOperations(scope string) {
	vs.RemoveScopedReplacer(scope)
	vs.RemoveScopedVariables(scope)
}

func (vs *VariableService) IsAllVariablesReplaced(line string) bool {
	match, err := regexp.MatchString("\\{{(.*?)\\}}", line)
	if err != nil {
		match = false
	}
	return !match
}
