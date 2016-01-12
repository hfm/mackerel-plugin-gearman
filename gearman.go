package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	flag "github.com/docker/docker/pkg/mflag"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

/* graphdef is Graph definition for mackerelplugin.

gearman.#.status:
	total ... The number of jobs in the queue
	running ... The number of running jobs
	available_workers ... The number of capable workers

see also: http://gearman.org/protocol/
*/
var graphdef = map[string](mp.Graphs){
	"gearman.status.#": mp.Graphs{
		Label: "Gearmand function status",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "total", Label: "total", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "running", Label: "running", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "available_workers", Label: "available_workers", Diff: false, Type: "uint64"},
		},
	},
}

// GearmanPlugin mackerel plugin for MoglieFS.
type GearmanPlugin struct {
	Target string
}

// FetchMetrics interface for mackerelplugin.
func (m GearmanPlugin) FetchMetrics() (map[string]interface{}, error) {
	raddr, err := net.ResolveTCPAddr("tcp", m.Target)
	if err != nil {
		_ = fmt.Errorf("Relosve error: %v\n", err)
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		_ = fmt.Errorf("DialTCP error: %v\n", err)
		return nil, err
	}

	fmt.Fprintln(conn, "status")

	return m.parseStatus(conn)
}

func (m GearmanPlugin) parseStatus(conn io.Reader) (map[string]interface{}, error) {
	scanner := bufio.NewScanner(conn)
	status := make(map[string]interface{})

	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)
		if s == "." {
			return status, nil
		}

		res := strings.Split(s, "\t")
		status["gearman.status."+res[0]+".total"] = res[1]
		status["gearman.status."+res[0]+".running"] = res[2]
		status["gearman.status."+res[0]+".available_workers"] = res[3]
	}

	if err := scanner.Err(); err != nil {
		return status, err
	}

	return nil, nil
}

// GraphDefinition interface for mackerelplugin.
func (m GearmanPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

// Parse flags and Run helper (MackerelPlugin) with the given arguments.
func main() {
	// Flags
	var (
		host     string
		port     string
		tempfile string
		version  bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)

	flags.StringVar(&host, []string{"H", "host"}, "127.0.0.1", "Host of gearmand")
	flags.StringVar(&port, []string{"p", "port"}, "4730", "Port of gearmand")
	flags.StringVar(&tempfile, []string{"t", "tempfile"}, "/tmp/mackerel-plugin-gearman", "Temp file name")
	flags.BoolVar(&version, []string{"v", "version"}, false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(os.Args[1:]); err != nil {
		os.Exit(ExitCodeError)
	}

	// Show version
	if version {
		fmt.Fprintf(os.Stderr, "%s version %s\n", Name, Version)
		os.Exit(ExitCodeOK)
	}

	// Create MackerelPlugin for Gearman
	var gearman GearmanPlugin
	gearman.Target = net.JoinHostPort(host, port)
	helper := mp.NewMackerelPlugin(gearman)
	helper.Tempfile = tempfile

	helper.Run()
}
