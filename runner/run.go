package runner

import (
	"bufio"
	"context"
	"errors"
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
	instanceNames := config.InstanceNames()

	for i, instanceName := range instanceNames {
		color := InstanceColors[i%len(InstanceColors)]

		wg.Add(1)
		go func() {
			defer wg.Done()

			err := runInstanceColored(ctx, instanceName, configPath, color)
			if err != nil {
				fmt.Println("Unable to start instance", instanceName, err)
			}
		}()
	}
	wg.Wait()

	return nil
}

func RunClusterFromNearestConfig(ctx context.Context) error {
	configPath, err := config.FindYamlFileAtPath(".")
	if err != nil {
		return errors.New("unable to seek for any of the configuration files in the current directory and it's parent directories")
	}

	fmt.Println("Found configuration at", configPath)

	err = RunClusterFromConfig(ctx, configPath)
	if err != nil {
		return fmt.Errorf("unable to start the cluster from the config file %s: %s", configPath, err)
	}

	return nil
}
