package token

import (
	"beaconCloud/rpcServer/user-srv/ssoEncoding"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"math"
	"strings"
	"sync"
	"time"
)

var (
	encodeX64  = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_")
	decodeX64  = make([]byte, 256)
	base64Code = base64.URLEncoding
)

type esToken struct {
	loginSn      uint64
	loginSnMutex sync.Mutex
	tokenType    byte
	tokenVer     byte
	equiId       string
}

func newToken(tokenType, tokenVer byte, equiId string) (t *esToken) {
	t = new(esToken)
	t.loginSn = 0
	t.tokenType = tokenType
	t.tokenVer = tokenVer
	t.equiId = equiId
	return
}

func (t *esToken) genLoginSn(i uint64) uint64 {
	t.loginSnMutex.Lock()
	defer t.loginSnMutex.Unlock()

	t.loginSn += i
	return t.loginSn
}

func (t *esToken) genRandSn(randBytes []byte) (sn uint64) {
	//计算剔除位
	b := 0
	for _, j := range randBytes {
		b += int(j)
	}
	b = b % len(randBytes)
	//计算随机数
	sn = 0
	weight := uint64(1)
	m := 0.0
	for k, j := range randBytes {
		if k != b {
			sn += uint64(j) * weight
			weight *= 256
			m++
		}
	}
	switch {
	case sn < 0:
		sn = sn + math.MaxInt32
	case sn == 0:
		sn = sn + 1
	}
	return

}

func (t *esToken) genVerifyCode(equiId, accX64, nodeGrpX64, userX64, timeX64,
	snX64, addX64 string) string {
	weights := []byte{2, 3, 5, 7, 11, 13, 17, 19, 23, 29,
		31, 37, 41, 43, 47, 53, 59, 61, 67, 71,
		73, 79, 83, 89, 97}
	weightLen := len(weights)
	verifyFunc := func(code string) int {
		codeV := 0
		codes := []byte(code)
		codeL := len(code)
		for i, c := range codes {
			codeV += int(weights[(codeL-1-i)%weightLen]) * int(c)
		}
		codeV = codeV % (64 * 64)
		return codeV
	}
	//计算节点编号校验结果
	equi := verifyFunc(equiId)

	//计算账号Hash校验结果
	acc := verifyFunc(accX64)

	//节点组校验结果
	nodeGrp := verifyFunc(nodeGrpX64)

	//用户校验结果
	user := verifyFunc(userX64)

	//计算系统时间校验结果
	time := verifyFunc(timeX64)

	//计算登录序号校验结果
	sn := verifyFunc(snX64)

	//计算登录附加码校验结果
	add := verifyFunc(addX64)

	//"="的AscII=61
	verify := equi*int(weights[0]) + acc*int(weights[1]) +
		nodeGrp*int(weights[2]) + user*int(weights[3]) + time*int(weights[4]) +
		sn*int(weights[5]) + add*int(weights[6])
	verify = verify % (64 * 64)
	//计算分隔符号
	verify = (verify + 61*int(weights[0])) % (64 * 64)
	verify = (verify + 61*int(weights[1])) % (64 * 64)
	verify = (verify + 61*int(weights[2])) % (64 * 64)
	verify = (verify + 61*int(weights[3])) % (64 * 64)
	verify = (verify + 61*int(weights[4])) % (64 * 64)
	verify = (verify + 61*int(weights[5])) % (64 * 64)
	verify = (verify + 61*int(weights[6])) % (64 * 64)
	return string([]byte{encodeX64[verify%64], encodeX64[verify/64]})
}

func (t *esToken) tokenVerifyFull(token, equiId string, account string) (ok bool) {
	p := strings.Split(token, ".")
	if len(p) != 8 {
		return false
	}

	accCode, errAcc := ssoEncoding.DecodeBase64ToUint64(p[1])
	if errAcc != nil {
		return false
	}

	//计算帐户的hash值
	accHash := fnv.New64()
	accHash.Write([]byte(account))

	if p[0] == equiId && accCode == accHash.Sum64() {
		code := ssoToken.genVerifyCode(p[0], p[1], p[2], p[3], p[4], p[5], p[6])
		return code == p[7]
	} else {
		return false
	}
}

func (t *esToken) tokenVerifyAccount(token, account string) (ok bool) {
	p := strings.Split(token, ".")
	if len(p) != 8 {
		return false
	}

	accCode, errAcc := ssoEncoding.DecodeBase64ToUint64(p[1])
	if errAcc != nil {
		return false
	}

	//计算帐户的hash值
	accHash := fnv.New64()
	accHash.Write([]byte(account))

	return accCode == accHash.Sum64()
}

func (t *esToken) genToken(account string, nodeGrpId int32, userId int64) (token string, err error) {
	return t.genTokenFull(t.equiId, nodeGrpId, userId, account)
}

//生成登录令牌,生成规则如下
//令牌     00.000000000000.000000000000.000000000000.000000000000.000000000000.000000000000.00
//子域    节点| 账号Hash64 |   节点组ID   |    用户Id  |   Unix时间 | 登录顺序号   |    附加码    |验证码
//实际长度  2 |     8     |      4      |      8     |     8     |     8      |      2      |2
//附加长度  0 |     1     |      5      |      1     |     1     |     1+4    |      7      |0
//该规则被探测命中的概率为：20*8=160位(2^160=1.4615016373309029182036848327163e+48)
func (t *esToken) genTokenFull(equiId string, nodeGrpId int32, userId int64, account string) (
	token string, err error) {
	randByte := make([]byte, (1+1)+(1+5)+(1+1)+(1+1)+(1+5)+(1+9)+8)
	n, err := rand.Read(randByte)
	if err != nil || n < len(randByte) {
		return
	}

	//计算帐户的hash值
	accHash := fnv.New64()
	accHash.Write([]byte(account))
	accX64 := ssoEncoding.EncodeUint64ToBase64(accHash.Sum64(), randByte[1:2]) //1byte随机数

	//当前节点组值
	nodeGrpX64 := ssoEncoding.EncodeUint32ToBase64(uint32(nodeGrpId), randByte[3:8]) //5byte随机数

	//用户Id
	userX64 := ssoEncoding.EncodeUint64ToBase64(uint64(userId), randByte[9:10]) //1byte随机数

	//获取系统当前时间
	now := time.Now()
	timeX64 := ssoEncoding.EncodeUint64ToBase64(uint64(now.Unix()), randByte[11:12]) //1byte随机数

	//获取随机数
	randSn := t.genRandSn(randByte[13:18]) //5byte随机数
	seedSn := uint64(randSn % math.MaxUint32)
	if seedSn <= 0 {
		seedSn = 1
	}
	//获取登录顺序号
	randSn = t.genLoginSn(seedSn)
	snX64 := ssoEncoding.EncodeUint64ToBase64(randSn, randByte[19:20]) //1byte随机数

	//获取附加码
	addRand := randByte[21:30] //9byte随机数
	addRand[0] = t.tokenVer    //令牌版本
	addRand[1] = t.tokenType   //令牌类型
	addX64 := base64Code.EncodeToString(addRand)

	err, token = nil, fmt.Sprintf("%s.%s.%s.%s.%s.%s.%s.%s", equiId,
		accX64, nodeGrpX64, userX64, timeX64, snX64, addX64,
		t.genVerifyCode(equiId, accX64, nodeGrpX64, userX64, timeX64, snX64, addX64))
	return
}

var ssoToken = new(esToken)

func GetToken(equiId string, nodeGrpId int32, userId int64, account string) (
	token string, err error) {
	return ssoToken.genTokenFull(equiId, nodeGrpId, userId, account)
}

func TokenVerifyFull(token, equiId string, account string) (ok bool) {
	return ssoToken.tokenVerifyFull(token, equiId, account)
}

func TokenVerifyAccount(token, account string) (ok bool) {
	return ssoToken.tokenVerifyAccount(token, account)
}

func TokenType(token string) (typ, ver byte, ok bool) {
	p := strings.Split(token, ".")
	if len(p) != 8 {
		return 0, 0, false
	}
	addCode, err := base64Code.DecodeString(p[6])
	ok = err == nil
	if ok {
		ver = addCode[0]
		typ = addCode[1]
	}
	return
}

func TokenRoute(token string) (nodeGrpId int32, userId int64, ok bool) {
	p := strings.Split(token, ".")
	if len(p) != 8 {
		return 0, 0, false
	}
	node, errNode := ssoEncoding.DecodeBase64ToUint32(p[2])
	user, errUser := ssoEncoding.DecodeBase64ToUint64(p[3])
	ok = errNode == nil && errUser == nil
	if ok {
		nodeGrpId = int32(node)
		userId = int64(user)
	}
	return
}
