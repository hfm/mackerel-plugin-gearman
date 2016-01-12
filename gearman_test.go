package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestGraphDefinition(t *testing.T) {
	var gearman GearmanPlugin

	graphdef := gearman.GraphDefinition()
	if len(graphdef) != 1 {
		t.Errorf("GetTempfilename: %d should be 1", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var gearman GearmanPlugin
	stub := `awesome_function	0	0	16
super_function	20	20	60
beyond_function	3	2	18
.
`

	gearmanStatus := bytes.NewBufferString(stub)

	status, err := gearman.parseStatus(gearmanStatus)
	fmt.Println(status)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestParseWithDot(t *testing.T) {
	var gearman GearmanPlugin
	stub := `function.with.dot	0	0	10
.
`

	gearmanStatus := bytes.NewBufferString(stub)

	status, err := gearman.parseStatus(gearmanStatus)
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Println(status)

	expect := map[string]interface{}{
		"gearman.status.function_with_dot.total": "0",
		"gearman.status.function_with_dot.running": "0",
		"gearman.status.function_with_dot.available_workers": "10",
	}
	fmt.Println(expect)

	if ! reflect.DeepEqual(status, expect) {
		t.Errorf("expects %s to equal %s", status, expect)
	}
}
