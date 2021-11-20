package rcmd

import (
	"errors"
	"fmt"
	"github.com/panjf2000/gnet/logging"
	"go-rwhps/core"
	"go-rwhps/server"
	"reflect"
	"strconv"
	"strings"
)

var QcCmd *Cmd

type ParseResult = int

type Variadic []string

const (
	Success ParseResult = iota
	NotEqualPrefix
	CommandNotFound
	ManyArguments
	FewArguments
	UnsuitableVariable
	ErrorArguments
	NoArguments
	VeryGoodDesign
	Failed
)

type Cmd struct {
	PrefixStr string

	Commands []*Command
}

func init() {
	QcCmd = NewCommandHandler("-qc ")
	type SetTeam struct {
		Target int
	}

	QcCmd.Register(&Command{
		Name:     "-self_team",
		Params:   reflect.TypeOf(SetTeam{}),
		Children: nil,
		Callback: func(sender *core.Player, args interface{}) error {
			sender.Team = args.(SetTeam).Target - 1
			fmt.Println("set successfully")
			return nil
		},
	})

	type SetMove struct {
		Target, TargetTeam int
	}
	QcCmd.Register(&Command{
		Name:     "-self_move",
		Params:   reflect.TypeOf(SetMove{}),
		Children: nil,
		Callback: func(sender *core.Player, args interface{}) error {
			d := args.(SetMove)
			server.GameServer.Players.Move(sender, d.Target-1, sender.Admin)
			if d.TargetTeam > -1 {
				sender.Team = d.TargetTeam
			}

			fmt.Println("move successfully")
			return nil
		},
	})
}

func NewCommandHandler(prefixStr string) *Cmd {
	return &Cmd{
		PrefixStr: prefixStr,
		Commands:  []*Command{},
	}
}

func (cmd *Cmd) Register(command *Command) {
	cmd.Commands = append(cmd.Commands, command)
}

func (cmd *Cmd) ParseChat(sender *core.Player, args string) (bool, ParseResult) {
	if strings.HasPrefix(strings.TrimSpace(args), cmd.PrefixStr) {
		return true, cmd.Parse(sender, args)
	} else {
		return false, NotEqualPrefix
	}
}

func (cmd *Cmd) Parse(sender *core.Player, args string) ParseResult {
	args = strings.TrimSpace(args)
	if args2 := strings.TrimPrefix(args, cmd.PrefixStr); args2 == args {
		return NotEqualPrefix
	} else {
		args = args2
	}
	strList := strings.Split(args, " ")
	fmt.Println(strList)
	//if cmd.Prefix != "" && !checkPrefix(cmd.Prefix, strList) {
	//	return NotContainPrefix
	//}
	var command *Command
	for _, c := range cmd.Commands {
		if strList[0] == c.Name {
			command = c
			break
		}
	}
	if command == nil {
		return CommandNotFound
	}
	if command.Params == nil && command.Callback != nil {
		call(command, sender, nil)
	} else if command.Children != nil {
		command2, args := search(command, strList)
		if command2 != nil {
			return call(command2, sender, args)
		} else if command.Callback != nil {
			return call(command, sender, strList[1:])
		} else {
			return ErrorArguments
		}
	} else if command.Params != nil && command.Callback != nil {
		return call(command, sender, strList[1:])
	}
	return VeryGoodDesign
}

func search(command *Command, args []string) (*Command, []string) {
	name := args[1]
	args = args[1:]
	for _, child := range command.Children {
		if child.Name == name {
			if child.Children != nil {
				c, args2 := search(command, args)
				if c != nil {
					return c, args2
				}
			}
			return child, args
		}
	}
	return nil, nil
}

//
//func checkPrefix(prefix string, strList []string) bool {
//	for i, s := range strList {
//		if !strings.HasPrefix(s, prefix) {
//			return false
//		} else {
//			strList[i] = strings.TrimPrefix(s, prefix)
//		}
//	}
//	return true
//}

func setInt(s string, v reflect.Value, bitSize int) error {
	parseInt, err := strconv.ParseInt(s, 10, bitSize)
	if err != nil {
		return err
	}
	v.SetInt(parseInt)
	return nil
}

func setUint(s string, v reflect.Value, bitSize int) error {
	parseUint, err := strconv.ParseUint(s, 10, bitSize)
	if err != nil {
		return err
	}
	v.SetUint(parseUint)
	return nil
}

func setFloat(s string, v reflect.Value, bitSize int) error {
	parseFloat, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return err
	}
	v.SetFloat(parseFloat)
	return nil
}

func setParams(s string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int64:
		err := setInt(s, v, 64)
		if err != nil {
			return err
		}
	case reflect.Int8, reflect.Int16, reflect.Int32:
		err := setInt(s, v, 32)
		if err != nil {
			return err
		}
	case reflect.Uint, reflect.Uint64:
		err := setUint(s, v, 64)
		if err != nil {
			return err
		}
	case reflect.Uint8, reflect.Uint32, reflect.Uint16:
		err := setUint(s, v, 32)
		if err != nil {
			return err
		}
	case reflect.Float32:
		err := setFloat(s, v, 64)
		if err != nil {
			return err
		}
	case reflect.Float64:
		err := setFloat(s, v, 32)
		if err != nil {
			return err
		}
	case reflect.String:
		v.SetString(s)
	default:
		return errors.New("invalid type: " + v.Kind().String())
	}
	return nil
}

func call(command *Command, sender *core.Player, args []string) ParseResult {
	var params interface{}
	if args != nil {
		var result ParseResult
		result, params = checkParams(command.Params, args)
		if result != Success {
			return result
		}
	}
	err := command.Callback(sender, params)
	if err != nil {
		logging.LogErr(err)
		return Failed
	}
	return Success
}

func checkParams(params reflect.Type, strList []string) (ParseResult, interface{}) {
	if params == nil && len(strList) != 0 {
		return NoArguments, nil
	}
	v := reflect.New(params).Elem()

	if v.Kind() == reflect.Struct {
		l := v.NumField()
		if l != len(strList) {
			if v.Field(l-1).Type().Name() == "Variadic" {
				if l < len(strList) {
					return FewArguments, nil
				}
			} else {
				if l > len(strList) {
					return ManyArguments, nil
				} else {
					return FewArguments, nil
				}
			}
		}
		for i := 0; i < l; i++ {
			s, v := strList[i], v.Field(i)
			if i != l-1 && v.Type().Name() == "Variadic" {
				return UnsuitableVariable, nil
			} else if i == l-1 && v.Type().Name() == "Variadic" {
				v.Set(reflect.ValueOf(strList[i:]))
				continue
			}
			err := setParams(s, v)
			if err != nil {
				logging.LogErr(err)
				return Failed, nil
			}
		}
	} else {
		if v.Type().Name() == "Variadic" {
			v.Set(reflect.ValueOf(strList))
		} else {
			if len(strList) > 1 {
				return ManyArguments, nil
			} else {
				err := setParams(strList[0], v)
				if err != nil {
					logging.LogErr(err)
					return Failed, nil
				}
			}
		}
	}

	return Success, v.Interface()
}

type Command struct {
	Name string

	Params reflect.Type

	Children []*Command

	Callback func(sender *core.Player, args interface{}) error
}
