package command_test

import (
	"os"
	"testing"

	"github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"go.uber.org/goleak"
	"gotest.tools/assert"
)

func TestExecuteScripts(t *testing.T) {
	goleak.VerifyNone(t)

	// make a temp directory
	tmpDir, err := os.MkdirTemp("", "casaos-test-*")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpDir)

	command.ExecuteScripts(tmpDir)

	// create a sample script under tmpDir
	script := tmpDir + "/test.sh"
	f, err := os.Create(script)
	assert.NilError(t, err)
	defer f.Close()

	// write a sample script
	_, err = f.WriteString("#!/bin/bash\necho 123")
	assert.NilError(t, err)

	command.ExecuteScripts(tmpDir)
}
