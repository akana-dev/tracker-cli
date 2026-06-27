package aliases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"tracker/internal/config"
)

var aliasesFile = filepath.Join(config.ConfigDir, "aliases.json")

var (
	aliasesCache map[string]string
	aliasesMu    sync.RWMutex
)

var aliasNameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

var reservedNames = map[string]bool{
	"tracker": true, "login": true, "logout": true, "me": true,
	"register": true, "configure": true, "server": true, "task": true,
	"company": true, "help": true, "completion": true, "alias": true,
	"tag": true, "template": true, "config": true,
}

func init() {
	aliasesCache = nil
}

func loadAliases() (map[string]string, error) {
	aliasesMu.RLock()
	if aliasesCache != nil {
		cached := aliasesCache
		aliasesMu.RUnlock()
		return cached, nil
	}
	aliasesMu.RUnlock()

	data, err := os.ReadFile(aliasesFile)
	if err != nil {
		if os.IsNotExist(err) {
			setAliasesCache(make(map[string]string))
			return aliasesCache, nil
		}
		return nil, err
	}

	var aliases map[string]string
	if err := json.Unmarshal(data, &aliases); err != nil {
		return nil, fmt.Errorf("ошибка парсинга алиасов: %w", err)
	}

	if aliases == nil {
		aliases = make(map[string]string)
	}

	setAliasesCache(aliases)
	return aliases, nil
}

func setAliasesCache(aliases map[string]string) {
	aliasesMu.Lock()
	aliasesCache = aliases
	aliasesMu.Unlock()
}

func saveAliases(aliases map[string]string) error {
	if err := os.MkdirAll(config.ConfigDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(aliasesFile, data, 0600); err != nil {
		return err
	}

	setAliasesCache(aliases)
	return nil
}

func Get(name string) (string, bool) {
	aliases, err := loadAliases()
	if err != nil {
		return "", false
	}
	value, ok := aliases[name]
	return value, ok
}

func List() (map[string]string, error) {
	return loadAliases()
}

func Add(name, value string) error {
	name = strings.TrimSpace(name)
	value = strings.TrimSpace(value)

	if name == "" {
		return fmt.Errorf("имя алиаса не может быть пустым")
	}
	if value == "" {
		return fmt.Errorf("значение алиаса не может быть пустым")
	}
	if !aliasNameRegex.MatchString(name) {
		return fmt.Errorf("имя алиаса может содержать только буквы, цифры, дефисы, подчёркивания и точки")
	}
	if reservedNames[name] {
		return fmt.Errorf("имя %q зарезервировано и не может быть алиасом", name)
	}

	aliases, err := loadAliases()
	if err != nil {
		return err
	}

	aliases[name] = value
	return saveAliases(aliases)
}

func Remove(name string) error {
	aliases, err := loadAliases()
	if err != nil {
		return err
	}

	if _, ok := aliases[name]; !ok {
		return fmt.Errorf("алиас %q не найден", name)
	}

	delete(aliases, name)
	return saveAliases(aliases)
}

func Exists(name string) bool {
	_, ok := Get(name)
	return ok
}

func ExpandArgs(args []string) []string {
	if len(args) < 2 {
		return args
	}

	possibleAlias := args[1]

	if strings.HasPrefix(possibleAlias, "-") {
		return args
	}
	if reservedNames[possibleAlias] {
		return args
	}

	expansion, ok := Get(possibleAlias)
	if !ok {
		return args
	}

	expandedArgs := splitArgs(expansion)

	newArgs := make([]string, 0, 1+len(expandedArgs)+len(args)-2)
	newArgs = append(newArgs, args[0])
	newArgs = append(newArgs, expandedArgs...)
	if len(args) > 2 {
		newArgs = append(newArgs, args[2:]...)
	}

	return newArgs
}

func splitArgs(s string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)

	for _, r := range s {
		switch {
		case !inQuotes && (r == '"' || r == '\''):
			inQuotes = true
			quoteChar = r
		case inQuotes && r == quoteChar:
			inQuotes = false
			quoteChar = 0
		case !inQuotes && (r == ' ' || r == '\t'):
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
