package command

import (
	"path/filepath"
	"testing"

	"github.com/IceWhaleTech/CasaOS-Common/utils/constants"
)

func TestExecutePostStartScripts(t *testing.T) {
	scriptDirectory := filepath.Join(constants.DefaultConfigPath, "start.d")
	ExecuteScripts(scriptDirectory)
}
