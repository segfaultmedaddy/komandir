package komandir

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"go.segfaultmedaddy.com/komandir/internal/sliceutil"
)

////////////////////////////
// Public API             //
////////////////////////////

type FlagMarsheller interface {
	UnmarshalFlag(string) error
	MarshalFlag() (string, error)
}

type Command[TFlagSet any, TArgSet any] struct {
	flagSet *TFlagSet
	argsSet *TArgSet

	parent   *Command[any, any]
	commands []*Command[any, any] // child commands

	Name    string
	Short   string
	Aliases string

	Pre  Hook
	Post Hook

	// Action receives the context and the command itself.
	Action func(context.Context, *Command[TFlagSet, TArgSet]) error
}

func (c *Command[TFlagSet, TArgSet]) AddCommand(cmd *Command[any, any]) {
	c.commands = append(c.commands, cmd)
}

func (cmd *Command[TFlagSet, TArgSet]) Exec(ctx context.Context, args ...string) error {
	if args == nil {
		args = os.Args[1:]
	}

	// Filter out empty strings if any.
	args = sliceutil.Filter(args, func(arg string) bool {
		return strings.TrimSpace(arg) != ""
	})

	if err := cmd.prepare(); err != nil {
		return err
	}

	if err := cmd.parse(args); err != nil {
		return err
	}

	if err := cmd.Action(ctx, cmd); err != nil {
		return err
	}

	return nil
}

// prepare prepares the command for execution.
//
// It parses the flags and arguments definitions.
// It must be called before the command line arguments are parsed.
//
// It recursively prepares all subcommands.
func (cmd *Command[TFlagSet, TArgSet]) prepare() error {
	var flagSet TFlagSet
	var argSet TArgSet

	parsedFlags, err := parseFlagSetDefinition(&flagSet)
	if err != nil {
		return fmt.Errorf("komandir: failed to parse flags: %w", err)
	}

	for _, flag := range parsedFlags {
		println(flag.Name)
	}

	cmd.flagSet = &flagSet
	cmd.argsSet = &argSet

	for _, subcmd := range cmd.commands {
		if err := subcmd.prepare(); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Command[_, _]) parse(args []string) error {
	flags := make([]string, 0)
	isInFlag := false

	for _, arg := range args {
		switch {
		case strings.Contains(arg, "-"):
			{
				isInFlag = true
			}

		case strings.Contains(arg, "--"):

		case strings.Contains(arg, "="):

		}
	}

	return nil
}

// help prints the help message for the command.
// It is a default implementation of the help command that can be overridden by
// the user.
func (cmd *Command[_, _]) help() string {
	var sb strings.Builder

	return sb.String()
}

////////////////////////////
// Parser Definition      //
////////////////////////////

var (
	flagFieldName         string = "name"
	flagFieldDesc         string = "desc"
	flagFieldEnv          string = "env"
	flagFieldDefaultValue string = "default"
)

var argFieldName string = "name"

// Flag represents a parsed flag from the command flag set.
type Flag struct {
	// Parsed
	Name         string   // Name of the flag.
	Aliases      []string // Alias names of the flag. Usually short names.
	Desc         string   // Description of the flag.
	EnvName      string   // Environment variable name for the flag.
	DefaultValue any      // Default value of the flag.

	Value any // Current value of the flag.
}

type Arg[T any] struct {
	// Parsed
	Name         string // Name of the argument.
	Desc         string // Description of the argument.
	DefaultValue T      // Default value of the argument.

	Value T // Current value of the argument.
}

func parseFlagSetDefinition[TFlagSet any](flagSet *TFlagSet) ([]*Flag, error) {
	v := reflect.ValueOf(flagSet)
	if v.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("komandir: flags must be a pointer")
	}

	if v.IsNil() {
		return nil, fmt.Errorf("komandir: flags must not be nil")
	}

	structValue := v.Elem()
	if structValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("komandir: flags must be a struct")
	}

	structValueType := structValue.Type()

	parsedFlags := make([]*Flag, 0)
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structValueType.Field(i)

		// Skip unexported fields
		if fieldType.PkgPath != "" {
			continue
		}

		if field.Kind() != reflect.Struct {
			continue
		}

		flag := &Flag{
			Name:    field.Type().Name(),
			Desc:    field.Type().Field(i).Tag.Get(flagFieldDesc),
			EnvName: field.Type().Field(i).Tag.Get(flagFieldEnv),
		}

		if flag.Name == "" {
			continue
		}

		parsedFlags = append(parsedFlags, flag)
	}

	return parsedFlags, nil
}

type Hook interface {
	Hook() error
}
