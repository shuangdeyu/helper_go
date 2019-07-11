package comhelper

import (
	"github.com/satori/go.uuid"
	"regexp"
)

const (
	CODE_TIME_BASE      = 1 // 时间序列生成，注意：严重不推荐
	CODE_NAME_HASH_MD5  = 3 //基于md5生成
	CODE_RANDOM         = 4 //基于随机数生成，建议使用这种方式
	CODE_NAME_HASH_SHA1 = 5 //基于sha1生成，
)

/**
 * 获取uuid
 */
func Get_uuid(code int) string {
	uu_id := ""
	switch code {
	case CODE_TIME_BASE: // 不推荐使用，因为要求使用系统命令，生产环境会有问题
		UUID, _ := uuid.NewV1()
		//UUID, _ := uuid.NewV1()
		uu_id = UUID.String()
	case CODE_NAME_HASH_MD5:
		UUID := uuid.NewV3(uuid.NamespaceDNS, "php.net")
		uu_id = UUID.String()
	case CODE_RANDOM:
		UUID, _ := uuid.NewV4()
		//UUID, _ := uuid.NewV4()
		uu_id = UUID.String()
	case CODE_NAME_HASH_SHA1:
		UUID := uuid.NewV5(uuid.NamespaceDNS, "php.net")
		uu_id = UUID.String()
	}
	return uu_id
}

/**
 * 校验uuid
 */
func Check_uuid(uu_id string) bool {
	len := len(uu_id)
	if len == 36 {
		if m, _ := regexp.MatchString(`^[0-9a-f]{8}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{12}$`, uu_id); !m {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}
