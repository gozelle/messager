package feishurobot

import (
	"context"
	"testing"
)

func TestRobot(t *testing.T) {
	robot := NewRobot(
		"https://open.feishu.cn/open-apis/bot/v2/hook/29850c30-17a2-4e32-884f-90c2b0dce9e1",
		"TV86eKOTzkQAHVL9MqkgOf",
	)

	err := robot.Push(context.Background(), "✅ 100 🟡 30 ❌ 40", `2023-04-20 13:34:05 [5s]

✅ 快照生成成功
❌ 邮件发送失败
🟡 表情意指面部表情，圖標则是图形標誌的意思，可用來代表多种表情，如笑脸表示笑、蛋糕表示食物等。在香港除「表情圖標」外，也有稱作「繪文字」或「emoji」。
`)

	t.Log(err)
}
