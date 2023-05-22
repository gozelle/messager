package notifier_test

import (
	"errors"
	"testing"
	"time"

	"github.com/gozelle/notifier"
	"github.com/gozelle/notifier/feishurobot"
)

func TestNotify(t *testing.T) {

	//robot := dingrobot.NewRobot(
	//	"https://oapi.dingtalk.com/robot/send?access_token=c8487958cbb991620877cbbe645e3ddd5c82cf0ab99681a8f70fae6c5e63d217",
	//	"SECa9d4d94de78b8b9afe45548dc25559c1a2ae00628a0b3ae372b05d9506c8af9e",
	//)
	robot := feishurobot.NewRobot(
		"https://open.feishu.cn/open-apis/bot/v2/hook/29850c30-17a2-4e32-884f-90c2b0dce9e1",
		"TV86eKOTzkQAHVL9MqkgOf",
	)
	notify := notifier.NewNotify(robot)

	go func() {
		notify.Run()
	}()

	notify.Infof("日报生成成功")
	notify.Infof("日报生成成功")
	notify.Infof("日报生成成功")
	notify.Infof("日报生成成功")
	notify.Infof("日报生成成功")
	notify.Errorf("数据统计错误: %s", errors.New("invalid params"))
	notify.Errorf("数据统计错误: %s", errors.New("invalid params"))
	notify.Errorf("数据统计错误: %s", errors.New("invalid params"))
	time.Sleep(7 * time.Second)
	notify.Errorf("数据统计错误: %s", errors.New("invalid params"))
	notify.Errorf("数据统计错误: %s", errors.New("invalid params"))
	notify.Warnf("已经 3 天没有更新了")

	select {}
}
