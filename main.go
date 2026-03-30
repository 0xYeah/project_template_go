// Package project_template_go /main.go
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/0xYeah/project_template_go/api"
	"github.com/0xYeah/project_template_go/common"
	"github.com/0xYeah/project_template_go/config"
	"github.com/0xYeah/project_template_go/custom_cmd"
	"github.com/george012/gtbox"
	"github.com/george012/gtbox/gtbox_cmd"
	"github.com/george012/gtbox/gtbox_log"
)

var (
	mRunMode       = ""
	mGitCommitHash = ""
	mGitCommitTime = ""
	mPackageOS     = ""
	mPackageTime   = ""
	mGoVersion     = ""
)

func setupApp() {
	runMode := gtbox.RunModeDebug
	switch mRunMode {
	case "debug":
		runMode = gtbox.RunModeDebug
	case "test":
		runMode = gtbox.RunModeTest
	case "release":
		runMode = gtbox.RunModeRelease
	default:
		runMode = gtbox.RunModeDebug
	}

	config.CurrentApp = config.NewApp(
		config.ProjectName,
		config.ProjectBundleID,
		fmt.Sprintf("%s service",
			config.ProjectName,
		),
		runMode,
	)

	if config.CurrentApp.CurrentRunMode == gtbox.RunModeDebug {
		cmdMap := map[string]string{
			"git_commit_hash": "git show -s --format=%H",
			"git_commit_time": "git show -s --format=\"%ci\" | cut -d ' ' -f 1,2 | sed 's/ /_/'",
			"build_os":        "go env GOOS",
			"go_version":      "go version | awk '{print $3}'",
		}
		cmdRes := gtbox_cmd.RunWith(cmdMap)

		if cmdRes != nil {
			mGitCommitHash = cmdRes["git_commit_hash"]
			mGitCommitTime = cmdRes["git_commit_time"]
			mPackageOS = cmdRes["build_os"]
			mGoVersion = cmdRes["go_version"]
			mPackageTime = time.Now().UTC().Format("2006-01-02_15:04:05")
		}
	}

	config.CurrentApp.GitCommitHash = mGitCommitHash
	config.CurrentApp.GitCommitTime = mGitCommitTime
	config.CurrentApp.GoVersion = mGoVersion
	config.CurrentApp.PackageOS = mPackageOS
	config.CurrentApp.PackageTime = mPackageTime

	custom_cmd.HandleCustomCmds(os.Args, config.CurrentApp)

	gtbox.SetupGTBox(config.CurrentApp.AppName,
		config.CurrentApp.CurrentRunMode,
		config.CurrentApp.AppLogPath,
		30,
		gtbox_log.GTLogSaveHours,
		config.CurrentApp.HTTPRequestTimeOut,
	)
}

func main() {
	runtime.LockOSThread()

	setupApp()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			gtbox_log.LogErrorf("系统崩溃:%v", r)
			gtbox_log.LogErrorf("系统崩溃:\n%s", debug.Stack())
		}
	}()

	if config.CurrentApp.CurrentRunMode == gtbox.RunModeDebug {
		config.CurrentApp.AppConfigFilePath = "./example_files/config_local.yaml"
	}

	config.SyncConfigFile(func(err error) {
		gtbox_log.LogErrorf("[sync config eeror]:%v", err)
	})

	api.StartAPIServices(config.GlobalConfig.ApiCfg)

	common.LoadSigHandle(nil, nil)

}
