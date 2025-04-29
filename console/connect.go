package console

import (
	"context"
	"fmt"
	"time"

	"github.com/abiosoft/readline"
	"github.com/fatih/color"
	"github.com/georgiy-belyanin/ttx/config"
	"github.com/tarantool/go-tarantool/v2"
)

func evalLua(text string) string {
	var evalLuaBase = `local yaml = require('yaml')
yaml.cfg{ encode_use_tostring = true }
local cmd = [[%s]]

local f, err = loadstring("return " .. cmd)
if not f then
	f, err = loadstring(cmd)
end
if not f then
	return yaml.encode({ error = err })
end

return yaml.encode({f()})`
	return fmt.Sprintf(evalLuaBase, text)
}

var connectTimeout = 1 * time.Second
var promptColor = color.New(color.FgCyan)

func connect(ctx context.Context, name string, dialer tarantool.Dialer) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := tarantool.Connect(ctx, dialer, tarantool.Opts{
		Timeout:    connectTimeout,
		SkipSchema: true,
	})
	if err != nil {
		return err
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       promptColor.Sprintf("%s> ", name),
		AutoComplete: CreateCompleter(conn),
		HistoryFile:  ".ttxhist",
	})
	if err != nil {
		return err
	}

	for {
		expr, err := rl.Readline()
		if err != nil {
			break
		}

		var data []string
		err = conn.Do(tarantool.NewEvalRequest(evalLua(expr))).GetTyped(&data)
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			fmt.Printf("%s\n", data[0])
		}
	}

	return nil
}

func ConnectByString(ctx context.Context, cfg *config.Config, s string) error {
	connectInfo := ParseConnectInfo(s)
	connectUris := make([]string, 0)

	if cfg != nil {
		instance, err := cfg.FindInstance(connectInfo.UriOrInstanceName)

		if err == nil && instance != nil {
			connectUris = instance.ConnectUris
		}
	}

	if len(connectUris) == 0 {
		connectUris = append(connectUris, connectInfo.UriOrInstanceName)
	}

	var err error = nil
	for _, connectUri := range connectUris {
		dialer := tarantool.NetDialer{
			Address:  connectUri,
			User:     connectInfo.User,
			Password: connectInfo.Password,
		}

		err = connect(ctx, connectInfo.UriOrInstanceName, dialer)
		if err == nil {
			break
		}
	}

	return err
}
