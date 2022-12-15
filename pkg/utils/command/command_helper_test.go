package command

import (
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestExecuteScripts(t *testing.T) {
	// make a temp directory
	tmpDir, err := os.MkdirTemp("", "casaos-test-*")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpDir)

	ExecuteScripts(tmpDir)

	// create a sample script under tmpDir
	script := tmpDir + "/test.sh"
	f, err := os.Create(script)
	assert.NilError(t, err)
	defer f.Close()

	// write a sample script
	_, err = f.WriteString("#!/bin/bash\necho 123")
	assert.NilError(t, err)

	ExecuteScripts(tmpDir)
}
