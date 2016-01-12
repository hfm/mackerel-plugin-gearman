package main

import (
	"bytes"
	"fmt"
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
