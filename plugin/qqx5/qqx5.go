package qqx5

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"io"
	"strings"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	fcext "github.com/FloatTech/floatbox/ctxext"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
)

var lines []string
var max_length int

func init() {
 
	engine := control.Register("qqx5", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "炫舞答题助手",
		Help:             "- 炫舞答题[xxx]",
	})
	
	getdb := fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		resp, err := http.Get("https://raw.githubusercontent.com/ahckjhckxz/QQX5FireTableGenerator/master/x5.txt")
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		text := string(body)
		raw_lines := strings.Split(text, "\r\n")
		
		for _, str := range raw_lines {
			if str != "" {
				lines = append(lines, str)
			}
		}
		max_length = 0
		for i :=1; i < len(lines); i++ {
			if max_length < len(lines[i]) {
			   max_length = len(lines[i])
			}
		}
		log.Infof("读取%d最长问题", max_length)
		return true
	})

	

	engine.OnRegex("^炫舞答题([\u4E00-\u9FA5A-Za-z0-9]{1,25})$", getdb).SetBlock(true).Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			search_string := ctx.State["regex_matched"].([]string)[1]
			var answers []string
			if len(search_string) > 0 && len(search_string) < max_length {
				for i, s := range lines {
				   if strings.Contains(s, search_string) {
					  if i + 1 < len(lines) {
						answers = append(answers, s)
						answers = append(answers, lines[i + 1])
					  }
				   }                  
				}
				result := strings.Join(answers,"\n")
				ctx.SendChain(message.Text("答案：\n" + result))
			 }
		})
}
