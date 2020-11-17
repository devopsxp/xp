/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/devopsxp/xp/pipeline"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	systemdService string // 服务名
	systemdState   string // 操作行为
	systemdReload  bool   // daemonReload
	systemdEnabled bool
)

// systemdCmd represents the systemd command
var systemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "用于管理服务运行状态",
	Long: `官方文档：https://docs.ansible.com/ansible/latest/modules/service_module.html#service-module
eg: ./xp cli systemd -n docker -s restart -r -e`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Cli args: %v", args)
		if systemdService == "" {
			log.Error("未指定需要操作的服务名")
			os.Exit(1)
		} else if len(args) == 0 {
			log.Error("未检测到目标主机，请确认！ [eg: ./xp cli systemd 127.0.0.1-20 -n docker -s status]")
			os.Exit(1)
		}

		data := map[string]interface{}{
			"host":        args,
			"remote_user": cliUser,
			"remote_pwd":  cliPwd,
			"remote_port": cliPort,
			"roles":       []interface{}{"systemd"},
			"stage":       []interface{}{"systemd"},
			"vars":        map[string]interface{}{},
			"hooks":       []interface{}{map[interface{}]interface{}{"type": "none"}},
			"config": []interface{}{map[interface{}]interface{}{
				"stage": "systemd",
				"name":  "服务管理",
				"systemd": map[interface{}]interface{}{
					"name":         systemdService,
					"state":        systemdState,
					"daemonReload": systemdReload,
					"enabled":      systemdEnabled,
				},
				"terminal": true,
			}},
		}

		config := pipeline.DefaultPipeConfig("systemd").
			WithInputName("cli").SetArgs(data).
			WithFilterName("shell").
			WithOutputName("console")

		p := pipeline.Of(*config)
		p.Init()
		p.Start()
		p.Exec()
		p.Stop()
	},
}

func init() {
	cliCmd.AddCommand(systemdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// systemdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// systemdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	systemdCmd.Flags().StringVarP(&systemdService, "name", "n", "", "服务名")
	systemdCmd.Flags().StringVarP(&systemdState, "state", "s", "", "服务状态，eg: start|stop|status|restart|reload")
	systemdCmd.Flags().BoolVarP(&systemdReload, "daemonReload", "r", false, "是否systemd daemon-reload")
	systemdCmd.Flags().BoolVarP(&systemdEnabled, "enable", "e", false, "是否开机启动")
}
