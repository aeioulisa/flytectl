package upgrade

import (
	"io"
	"sort"
	"testing"

	"github.com/flyteorg/flytectl/pkg/util"
	"github.com/flyteorg/flytectl/pkg/util/githubutil"

	"github.com/flyteorg/flytectl/pkg/util/platformutil"

	"github.com/flyteorg/flyteidl/clients/go/admin/mocks"
	stdlibversion "github.com/flyteorg/flytestdlib/version"

	"context"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var (
	version = "v0.2.20"
	tempExt = "flyte.ext"
)

func TestUpgradeCommand(t *testing.T) {
	rootCmd := &cobra.Command{
		Long:              "flytectl is CLI tool written in go to interact with flyteadmin service",
		Short:             "flyetcl CLI tool",
		Use:               "flytectl",
		DisableAutoGenTag: true,
	}
	upgradeCmd := SelfUpgrade(rootCmd)
	cmdCore.AddCommands(rootCmd, upgradeCmd)
	assert.Equal(t, len(rootCmd.Commands()), 1)
	cmdNouns := rootCmd.Commands()
	// Sort by Use value.
	sort.Slice(cmdNouns, func(i, j int) bool {
		return cmdNouns[i].Use < cmdNouns[j].Use
	})

	assert.Equal(t, cmdNouns[0].Use, "upgrade")
	assert.Equal(t, cmdNouns[0].Short, upgradeCmdShort)
	assert.Equal(t, cmdNouns[0].Long, upgradeCmdLong)
}

//
func TestUpgrade(t *testing.T) {
	_ = util.WriteIntoFile([]byte("data"), tempExt)
	stdlibversion.Version = version
	githubutil.FlytectlReleaseConfig.OverrideExecutable = tempExt
	t.Run("Successful upgrade", func(t *testing.T) {
		message, err := upgrade(githubutil.FlytectlReleaseConfig)
		assert.Nil(t, err)
		assert.Contains(t, message, "Successfully updated to version")
	})
}

func TestCheckGoosForRollback(t *testing.T) {
	stdlibversion.Version = version
	linux := platformutil.Linux
	windows := platformutil.Windows
	darwin := platformutil.Darwin
	githubutil.FlytectlReleaseConfig.OverrideExecutable = tempExt
	t.Run("checkGOOSForRollback on linux", func(t *testing.T) {
		assert.Equal(t, true, isRollBackSupported(linux))
		assert.Equal(t, false, isRollBackSupported(windows))
		assert.Equal(t, true, isRollBackSupported(darwin))
	})
}

func TestIsUpgradeable(t *testing.T) {
	stdlibversion.Version = version
	githubutil.FlytectlReleaseConfig.OverrideExecutable = tempExt
	linux := platformutil.Linux
	windows := platformutil.Windows
	darwin := platformutil.Darwin
	t.Run("IsUpgradeable on linux", func(t *testing.T) {
		check, err := isUpgradeSupported(linux)
		assert.Nil(t, err)
		assert.Equal(t, true, check)
	})
	t.Run("IsUpgradeable on darwin", func(t *testing.T) {
		check, err := isUpgradeSupported(darwin)
		assert.Nil(t, err)
		assert.Equal(t, true, check)
	})
	t.Run("IsUpgradeable on darwin using brew", func(t *testing.T) {
		check, err := isUpgradeSupported(darwin)
		assert.Nil(t, err)
		assert.Equal(t, true, check)
	})
	t.Run("isUpgradeSupported failed", func(t *testing.T) {
		stdlibversion.Version = "v"
		check, err := isUpgradeSupported(linux)
		assert.NotNil(t, err)
		assert.Equal(t, false, check)
		stdlibversion.Version = version
	})
	t.Run("isUpgradeSupported windows", func(t *testing.T) {
		check, err := isUpgradeSupported(windows)
		assert.Nil(t, err)
		assert.Equal(t, false, check)
	})
}

func TestSelfUpgrade(t *testing.T) {
	stdlibversion.Version = version
	githubutil.FlytectlReleaseConfig.OverrideExecutable = tempExt
	goos = platformutil.Linux
	t.Run("Successful upgrade", func(t *testing.T) {
		ctx := context.Background()
		var args []string
		mockClient := new(mocks.AdminServiceClient)
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
		stdlibversion.Build = ""
		stdlibversion.BuildTime = ""
		stdlibversion.Version = version

		assert.Nil(t, selfUpgrade(ctx, args, cmdCtx))
	})
}

func TestSelfUpgradeError(t *testing.T) {
	stdlibversion.Version = version
	githubutil.FlytectlReleaseConfig.OverrideExecutable = tempExt
	goos = platformutil.Linux
	t.Run("Successful upgrade", func(t *testing.T) {
		ctx := context.Background()
		var args []string
		mockClient := new(mocks.AdminServiceClient)
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
		stdlibversion.Build = ""
		stdlibversion.BuildTime = ""
		stdlibversion.Version = "v"

		assert.NotNil(t, selfUpgrade(ctx, args, cmdCtx))
	})

}

func TestSelfUpgradeRollback(t *testing.T) {
	stdlibversion.Version = version
	githubutil.FlytectlReleaseConfig.OverrideExecutable = tempExt
	goos = platformutil.Linux
	t.Run("Successful rollback", func(t *testing.T) {
		ctx := context.Background()
		var args = []string{rollBackSubCommand}
		mockClient := new(mocks.AdminServiceClient)
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
		stdlibversion.Build = ""
		stdlibversion.BuildTime = ""
		stdlibversion.Version = version
		assert.Nil(t, selfUpgrade(ctx, args, cmdCtx))
	})

	t.Run("Successful rollback failed", func(t *testing.T) {
		ctx := context.Background()
		var args = []string{rollBackSubCommand}
		mockClient := new(mocks.AdminServiceClient)
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
		stdlibversion.Build = ""
		stdlibversion.BuildTime = ""
		stdlibversion.Version = "v100.0.0"
		assert.NotNil(t, selfUpgrade(ctx, args, cmdCtx))
	})

	t.Run("Successful rollback for windows", func(t *testing.T) {
		ctx := context.Background()
		var args = []string{rollBackSubCommand}
		mockClient := new(mocks.AdminServiceClient)
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
		stdlibversion.Build = ""
		stdlibversion.BuildTime = ""
		stdlibversion.Version = version
		goos = platformutil.Windows
		assert.Nil(t, selfUpgrade(ctx, args, cmdCtx))
	})

	t.Run("Successful rollback for windows", func(t *testing.T) {
		ctx := context.Background()
		var args = []string{rollBackSubCommand}
		mockClient := new(mocks.AdminServiceClient)
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
		stdlibversion.Build = ""
		stdlibversion.BuildTime = ""
		stdlibversion.Version = version
		githubutil.FlytectlReleaseConfig.OverrideExecutable = "/"
		assert.Nil(t, selfUpgrade(ctx, args, cmdCtx))
	})

}
