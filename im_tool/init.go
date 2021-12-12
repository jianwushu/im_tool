package jingdong

import (
	"strconv"
	"fmt"
	"os"
	"bufio"
	"strings"
	"path/filepath"
	"github.com/cdle/sillyGirl/core"
)

var qq = core.NewBucket("qq")
var groups = qq.Get("onGroups")
var menu = core.Bucket("menu")

var ExecPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

func init() {
	// 加载配置
	readMenu(ExecPath+"/conf/")
	//向core包中添加命令
	core.AddCommand("", []core.Function{
		{
			Admin: true,
			Rules: []string{"^去监听 ?"},
			Handle: func(s core.Sender) interface{} {
				id := s.Get(0)
				arr :=strings.Split(groups,"&")
				for i,p := range arr {
					if(p == id){
						return "已经监听群号:" + id
					}
					_ = i
				}
				groups = groups + "&" + id
				qq.Set("onGroups",groups)
				return "成功监听群号:" + id
			},
		},
		{
			Admin: true,
			Rules: []string{"^取消监听本群"},
			Handle: func(s core.Sender) interface{} {
				id := s.GetChatID()
				arr :=strings.Split(groups,"&")
				
				new_groups := []string{}
				str := strconv.Itoa(id)
				for i,p := range arr {
					if(p == str){
						continue
					}
					_ = i
					new_groups = append(new_groups,p)
				}				
				qq.Set("onGroups",strings.Replace(strings.Trim(fmt.Sprint(new_groups), "[]"), " ", "&", -1))
				return "成功取消监听群号:" + str
			},
		},
		{
			Rules: []string{"^菜单"},
			Handle: func(s core.Sender) interface{} {
				menu_info := menu.Get("info")
				return menu_info
			},
		},
	})
}

//菜单
func readMenu(confDir string) {
	path := confDir + "menu.txt"
	s := []string{}
	if _, err := os.Stat(confDir); err != nil {
		os.MkdirAll(confDir, os.ModePerm)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			s = append(s,line)
		}
	}
	ss := fmt.Sprintf(strings.Join(s, "\n"))
	menu.Set("info",ss)
	f.Close()
}