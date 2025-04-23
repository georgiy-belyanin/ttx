package runner

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/georgiy-belyanin/ttx/config"
)

var InstanceColors = []*color.Color{
	color.New(color.FgCyan),
	color.New(color.FgMagenta),
	color.New(color.FgRed),
	color.New(color.FgYellow),
	color.New(color.FgGreen),
}

var DefaultColor = color.New(color.Reset)
var InfoColor = color.New(color.Reset)
var WarnColor = color.New(color.FgYellow)
var ErrorColor = color.New(color.FgRed)
var FatalColor = color.New(color.FgRed)

func colorLogMessage(msg string) string {
	markersAndColors := [](struct {
		string
		*color.Color
	}){
		{"I>", InfoColor},
		{"W>", WarnColor},
		{"E>", ErrorColor},
		{"F>", FatalColor},
	}

	msgColor := DefaultColor
	for _, markerAndColor := range markersAndColors {
		if strings.Contains(msg, markerAndColor.string) {
			msgColor = markerAndColor.Color
		}
	}

	return msgColor.Sprint(msg)
}

func runInstanceColored(ctx context.Context, instanceName string, configPath string, color *color.Color) error {
	c1 := exec.Command("tarantool", "--name", instanceName, "--config", configPath)

	stderr, err := c1.StderrPipe()
	if err != nil {
		return err
	}

	scn := bufio.NewScanner(stderr)
	err = c1.Start()
	if err != nil {
		return err
	}

	coloredInstanceName := instanceName
	if color != nil {
		coloredInstanceName = color.Sprint(instanceName)
	}

	for scn.Scan() {
		coloredMsg := colorLogMessage(scn.Text())
		fmt.Println(coloredInstanceName, "|", coloredMsg)
	}

	<-ctx.Done()

	return c1.Wait()
}

func RunClusterFromConfig(ctx context.Context, configPath string) error {
	config, err := config.LoadYamlFile(configPath)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	instances, err := config.Instances()
	if err != nil {
		return err
	}

	for i, instance := range instances {
		color := InstanceColors[i%len(InstanceColors)]

		wg.Add(1)
		go func() {
			defer wg.Done()

			err := runInstanceColored(ctx, instance.Name, configPath, color)
			if err != nil {
				fmt.Println("Unable to start instance", instance.Name, err)
			}
		}()
	}
	wg.Wait()

	return nil
}

func RunClusterPartsFromConfig(ctx context.Context, configPath string, parts []string) error {
	conf, err := config.LoadYamlFile(configPath)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	instances, err := conf.Instances()
	if err != nil {
		return err
	}

	instancesToStart := make([]config.Instance, 0)
	for _, instance := range instances {
		for _, part := range parts {
			if instance.Name == part {
				instancesToStart = append(instancesToStart, instance)
				break
			}
			if instance.Replicaset == part {
				instancesToStart = append(instancesToStart, instance)
				break
			}
			if instance.Group == part {
				instancesToStart = append(instancesToStart, instance)
				break
			}
		}
	}

	for i, instance := range instancesToStart {
		color := InstanceColors[i%len(InstanceColors)]

		wg.Add(1)
		go func() {
			defer wg.Done()

			err := runInstanceColored(ctx, instance.Name, configPath, color)
			if err != nil {
				fmt.Println("Unable to start instance", instance.Name, err)
			}
		}()
	}
	wg.Wait()

	return nil
}
