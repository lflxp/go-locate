/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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

	"github.com/lflxp/go-locate/pkg"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	initPath  string
	gonum     int
	timesleep int
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化数据",
	Long: `重新初始化系统所有数据 For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.GetAllFile(initPath, gonum, timesleep)
		if err != nil {
			log.Errorln(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	home, err := homedir.Dir()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	initCmd.Flags().StringVarP(&initPath, "path", "p", home, "扫描路径，默认：~/")
	initCmd.Flags().IntVarP(&gonum, "go", "g", 900000, "最大并发goroutine数量")
	initCmd.Flags().IntVarP(&timesleep, "timesleep", "T", 500, "携程最大等待时间 单位：ms")
}
