package gutils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
)

// 密钥
type SecretKey struct {
	Curve      string
	PrivateKey string
	PublicKey  string
}

// 新建一个密钥
func NewSecreKey() *SecretKey {
	s := new(SecretKey)
	s.GenearteKey()
	return s
}

// 通过文件获取密钥
func NewSecreKeyByFile(path string) (*SecretKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	secretKey := new(SecretKey)
	err = json.Unmarshal(keyBytes, secretKey)
	return secretKey, err
}

// 通过椭圆曲线算法和随机数生成密钥
func (s *SecretKey) GenearteKey() {
	curve := s.SecrettKeyCurve()
	//生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("密钥创建失败")
	}
	encode := base64.StdEncoding
	s.PrivateKey = encode.EncodeToString(privateKey.D.Bytes())
	buffer := bytes.NewBuffer(nil)
	buffer.Write(privateKey.X.Bytes())
	buffer.Write([]byte("+++"))
	buffer.Write(privateKey.Y.Bytes())
	s.PublicKey = encode.EncodeToString(buffer.Bytes())
}

// 使用私钥进行数据签名判断数据是否有被更改及其他使用
func (s *SecretKey) Signature(msg string) string {
	PrivateKey, err := s.ReducePrivate()
	if err != nil {
		fmt.Println("私钥获取失败")
		return ""
	}
	sign, err := ecdsa.SignASN1(rand.Reader, PrivateKey, HashAlgorithm(msg))
	if err != nil {
		fmt.Println("签名失败")
		return ""
	}
	return base64.StdEncoding.EncodeToString(sign)
}

// 使用公钥对数据进行验证
func (s *SecretKey) Verify(msg, sign string) bool {
	publicKey, err := s.ReducePublic()
	if err != nil {
		fmt.Println("私钥获取失败")
		return false
	}
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		fmt.Println("数据编码失败:", err)
		return false
	}
	return ecdsa.VerifyASN1(publicKey, HashAlgorithm(msg), signBytes)
}

// 默认选取P256曲线
func (s *SecretKey) SecrettKeyCurve() elliptic.Curve {
	switch s.Curve {
	case "P224":
		return elliptic.P224()
	case "P256":
		return elliptic.P256()
	case "P384":
		return elliptic.P384()
	case "P521":
		return elliptic.P521()
	default:
		return elliptic.P256()
	}
}

// 还原公钥
func (s *SecretKey) ReducePublic() (*ecdsa.PublicKey, error) {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(s.PublicKey)
	if err != nil {
		return nil, err
	}
	XY := bytes.Split(pubKeyBytes, []byte("+++"))
	if len(XY) != 2 {
		return nil, errors.New("公钥解析失败")
	}
	publickey := &ecdsa.PublicKey{
		Curve: s.SecrettKeyCurve(),
		X:     big.NewInt(0).SetBytes(XY[0]),
		Y:     big.NewInt(0).SetBytes(XY[1]),
	}
	return publickey, nil
}

// 还原私钥
func (s *SecretKey) ReducePrivate() (*ecdsa.PrivateKey, error) {
	priKeyBytes, err := base64.StdEncoding.DecodeString(s.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: s.SecrettKeyCurve(),
		},
		D: big.NewInt(0).SetBytes(priKeyBytes),
	}, nil
}

// 通过信息摘要算法加密数据
func HashAlgorithm(msg string) []byte {
	hash := sha256.New()
	hash.Write([]byte(msg))
	return hash.Sum(nil)
}

// 保存密钥到本地
func SavaSecretKey(path string, s *SecretKey) error {
	key, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, key, 0755)
	return err
}
