package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"community-elderly-care-platform/pkg/crypto"
)

// 测试用固定密钥（32 字节）
// 生产环境请使用 openssl rand -base64 32 生成
var testEncryptKey = []byte("scare-test-id-card-key-32-bytes!")
var testTokenSecret = "your_jwt_secret_key_here"
var hmacKey = []byte(testTokenSecret + ":id_card_token")

type mockUser struct {
	id     int64
	name   string
	idCard string
}

func main() {
	users := []mockUser{
		// B端用户
		{1, "系统管理员", "110101198001010011"},
		{2, "李站长", "110101197505150022"},
		{3, "王站长", "110101197808200033"},
		{4, "王小红", "110101199003100044"},
		{5, "刘师傅", "110101198507220055"},
		{6, "陈护士", "110101199211050066"},
		{7, "赵大哥", "110101198809180077"},
		// C端用户
		{8, "张大爷", "110101195005150088"},
		{9, "李奶奶", "110101195503200099"},
		{10, "王爷爷", "110101194811100100"},
		{11, "孙女士", "110101199006250111"},
		{12, "赵先生", "110101196502140122"},
		{13, "小明", "110101201803150133"},
		// 新增 C端用户
		{14, "周阿姨", "110101194903080144"},
		{15, "吴大爷", "110101195107120155"},
		{16, "郑奶奶", "110101195209250166"},
		{17, "冯爷爷", "110101194605180177"},
		{18, "陈阿姨", "110101195312020188"},
		{19, "杨大爷", "110101195008150199"},
		{20, "黄女士", "110101199205100200"},
		{21, "林先生", "110101197003200211"},
		{22, "何大爷", "110101194712080222"},
		{23, "马奶奶", "110101195601150233"},
		// 新增 B端用户
		{24, "孙小明", "110101199508120244"},
		{25, "周护工", "110101199112050255"},
	}

	fmt.Println("-- 身份证号加密数据（使用测试密钥生成）")
	fmt.Println("-- 测试加密密钥 base64:", base64.StdEncoding.EncodeToString(testEncryptKey))
	fmt.Println()

	for _, u := range users {
		encrypted, err := crypto.Encrypt(u.idCard, testEncryptKey)
		if err != nil {
			fmt.Printf("-- ERROR: 用户 %d (%s) 加密失败: %v\n", u.id, u.name, err)
			continue
		}

		mac := hmac.New(sha256.New, hmacKey)
		mac.Write([]byte("id_card:"))
		mac.Write([]byte(u.idCard))
		hmacVal := hex.EncodeToString(mac.Sum(nil))

		masked := maskIDCard(u.idCard)

		fmt.Printf("-- 用户 %d: %s (身份证: %s)\n", u.id, u.name, u.idCard)
		fmt.Printf("UPDATE `users` SET `id_card` = '%s', `id_card_hmac` = '%s', `id_card_masked` = '%s' WHERE `id` = %d;\n",
			encrypted, hmacVal, masked, u.id)
		fmt.Println()
	}
}

func maskIDCard(idCard string) string {
	if idCard == "" {
		return ""
	}
	runes := []rune(idCard)
	n := len(runes)
	if n <= 8 {
		if n <= 2 {
			return strings.Repeat("*", n)
		}
		return string(runes[:1]) + strings.Repeat("*", n-2) + string(runes[n-1:])
	}
	return string(runes[:4]) + strings.Repeat("*", n-8) + string(runes[n-4:])
}
