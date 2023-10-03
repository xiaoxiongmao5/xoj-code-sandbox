package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

/** 生成登录验证token（jwt）
 */
func GenerateToken(userAccount, userRole string) (string, error) {
	claims := jwt.MapClaims{
		"user_account": userAccount,
		"user_role":    userRole,
		"exp":          time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("SecretKey"))
}

/** 加密密码（bcrypt 哈希算法）
 */
func HashPasswordByBcrypt(password string) (string, error) {
	// bcrypt.DefaultCost: 这是 bcrypt 哈希算法的工作因子（cost factor），表示计算哈希时使用的迭代次数。工作因子越高，计算哈希所需的时间和资源就越多，因此更难受到暴力破解。bcrypt.DefaultCost 是库中预定义的默认工作因子值。
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

/** 校验加密密码
 */
func CheckHashPasswordByBcrypt(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}

/** 生成包含N个随机数字的字符串
 */
func GenetateRandomString(length int) string {
	// 设置随机数种子，以确保每次运行生成的随机数都不同
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 定义一个包含数字字符的字符集
	charset := "0123456789"
	charsetLength := len(charset)

	// 生成随机数字并拼接字符串
	randomString := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex := r.Intn(charsetLength)
		randomChar := charset[randomIndex]
		randomString[i] = randomChar
	}
	return string(randomString)
}

// /** 生成随机字符串
//  */
// func GenerateRandomKey(length int) (string, error) {
// 	randomBytes := make([]byte, length)
// 	// 生成随机的字节序列
// 	_, err := rand.Read(randomBytes)
// 	if err != nil {
// 		return "", err
// 	}
// 	// 转为base64编码的字符串
// 	return base64.StdEncoding.EncodeToString(randomBytes), nil
// }

/** 生成带盐的哈希值（SHA-256 哈希算法）
 */
func HashBySHA256WithSalt(data, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data + salt))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func MD5(str string) string {
	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}
