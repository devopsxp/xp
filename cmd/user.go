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
	userName     string
	userPassword string
	userDelete   bool
	userLock     bool
	userUnLock   bool
	userForce    bool
	userStatus   bool
	userMax      string
	userMin      string
	userWarn     string
	userInactive string
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "远程批量用户管理",
	Long:  `官方文档：https://docs.ansible.com/ansible/latest/modules/user_module.html#user-module`,
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
			"roles":       []interface{}{"user"},
			"stage":       []interface{}{"user"},
			"vars":        map[string]interface{}{},
			"hooks":       []interface{}{map[interface{}]interface{}{"type": "none"}},
			"config": []interface{}{map[interface{}]interface{}{
				"stage": "user",
				"name":  "用户管理",
				"user": map[interface{}]interface{}{
					"name":     userName,
					"password": userPassword,
					"delete":   userDelete,
					"lock":     userLock,
					"unlock":   userUnLock,
					"force":    userForce,
					"status":   userStatus,
					"maximum":  userMax,
					"minimum":  userMin,
					"warning":  userWarn,
					"inactive": userInactive,
				},
				"terminal": true,
			}},
		}

		config := pipeline.DefaultPipeConfig("user").
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
	cliCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	userCmd.Flags().StringVarP(&userName, "user", "U", "", "用户名")
	userCmd.Flags().StringVarP(&userPassword, "password", "p", "", "修改密码")
	userCmd.Flags().StringVarP(&userMax, "maximum", "x", "", "两次密码修正的最大天数，后面接数字；仅能root权限操作")
	userCmd.Flags().StringVarP(&userMin, "minimum", "n", "", "两次密码修改的最小天数，后面接数字，仅能root权限操作")
	userCmd.Flags().StringVarP(&userWarn, "warning", "w", "", "在距多少天提醒用户修改密码；仅能root权限操作")
	userCmd.Flags().StringVarP(&userInactive, "inactive", "i", "", "在密码过期后多少天，用户被禁掉，仅能以root操作")
	userCmd.Flags().BoolVarP(&userDelete, "delete", "D", false, "删除用户密码，仅能以root权限操作")
	userCmd.Flags().BoolVarP(&userLock, "lock", "L", false, "锁住用户无权更改其密码，仅能通过root权限操作")
	userCmd.Flags().BoolVarP(&userUnLock, "unlock", "u", false, "解除锁定")
	userCmd.Flags().BoolVarP(&userForce, "force", "f", false, "强制操作；仅root权限才能操作")
	userCmd.Flags().BoolVarP(&userStatus, "status", "S", false, "查询用户的密码状态，仅能root用户操作")
}
