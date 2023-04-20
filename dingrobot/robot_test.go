package dingrobot

import "testing"

func TestRobot(t *testing.T) {
	robot := NewRobot(
		"https://oapi.dingtalk.com/robot/send?access_token=c8487958cbb991620877cbbe645e3ddd5c82cf0ab99681a8f70fae6c5e63d217",
		"SECa9d4d94de78b8b9afe45548dc25559c1a2ae00628a0b3ae372b05d9506c8af9e",
	)
	
	err := robot.Request("✅100 ⚠️30 ❌40", `
	2023-04-20 13:34:05 [5s]
	✅ 快照生成成功
	❌ 邮件发送失败
	⚠️ 表情意指面部表情，圖標则是图形標誌的意思，可用來代表多种表情，如笑脸表示笑、蛋糕表示食物等。在香港除「表情圖標」外，也有稱作「繪文字」或「emoji」。
`)
	
	t.Log(err)
}
