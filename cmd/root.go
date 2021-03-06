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
	"fmt"
	"os"
	"time"

	"github.com/lflxp/go-locate/pkg"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	logLevel bool
	all      bool
	info     string
	isKey    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-locate",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if info == "" {
			log.Errorln("查询字段不能为空, -i")
			return
		}
		now := time.Now()
		if !all {
			err := pkg.SearchAll(info, isKey)
			if err != nil {
				log.Errorln(err)
			}
			elapsed := time.Since(now)
			log.WithField("耗时", fmt.Sprint(elapsed)).Println("search all finished")
		} else {
			err := pkg.SearchPrefix(info)
			if err != nil {
				log.Errorln(err)
			}
			elapsed := time.Since(now)
			log.WithField("耗时", fmt.Sprint(elapsed)).Println("search all prefix finished")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-locate.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 日志级别
	rootCmd.PersistentFlags().BoolVarP(&logLevel, "debug", "v", false, "是否debug日志输出")
	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "扫描方式： 默认全表扫描，true为prefix扫描")
	rootCmd.Flags().BoolVarP(&isKey, "isKey", "I", false, "扫描方式： 默认按Key键进行扫描，true为value扫描")
	rootCmd.Flags().StringVarP(&info, "info", "i", "", "support regexp, check message,default: ''")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-locate" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-locate")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// 设置日志格式为json格式
	log.SetFormatter(&log.JSONFormatter{})

	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)

	// 设置日志级别为warn以上
	if logLevel {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

}
