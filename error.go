// 错误模块

package filter

import (
	"strings"
)

type Error struct {
	Type   ErrorType
	Fields []string
}

type ErrorType uint

func (e *Error) Error() string {
	e.Fields = append([]string{appErrorCodes[e.Type]}, e.Fields...)
	return strings.Join(e.Fields, ":")
}

// application error
// value should keep synchronous with api.Error*
const (
	ErrorObjectNotExist = iota
	ErrorObjectDuplicated
	ErrorNoObjectUpdated
	ErrorNoObjectDeleted
	ErrorMissingParam
	ErrorInvalidParam
	ErrorQuotaExceed
	ErrorPermissionDenied
	ErrorOperationFailed
	ErrorInternalError
)

// application error codes
var appErrorCodes = map[ErrorType]string{
	ErrorObjectNotExist:   "ObjectNotExist",
	ErrorObjectDuplicated: "ObjectDuplicated",
	ErrorNoObjectUpdated:  "NoObjectUpdated",
	ErrorNoObjectDeleted:  "NoObjectDeleted",
	ErrorMissingParam:     "MissingParam",
	ErrorInvalidParam:     "InvalidParam",
	ErrorQuotaExceed:      "QuotaExceed",
	ErrorPermissionDenied: "PermissionDenied",
	ErrorOperationFailed:  "OperationFailed",
	ErrorInternalError:    "InternalError",
}

// NewError return a filter error.
func NewError(errorType ErrorType, fields ...string) *Error {
	return &Error{errorType, fields}
}

var ErrorWordMap = map[string]map[string]string{
	"en_us": map[string]string{
		// Internal error
		"InvalidValidator": "invalid validator",

		// Common
		"NotInSet":     "not in set",
		"ItemNotInSet": "item not in set",

		// Interger
		"NotInt":    "not int",
		"NotInt64":  "not int64",
		"NotUint":   "not uint",
		"NotUint64": "not uint64",
		"TooSmall":  "too small",
		"TooLarge":  "too large",

		// String
		"NotString":       "not string",
		"TooShort":        "too short",
		"TooLong":         "too long",
		"WrongFormat":     "wrong format",
		"NotNumeric":      "not numeric",
		"NotDigit":        "not digit",
		"NotAlpha":        "not alpha",
		"NotAlphaNumeric": "not alphanumeric",

		// Email
		"NotEmail": "not email",

		// IP、CIDR
		"NotCIDR": "not cidr",
		"NotIP":   "not ip",
		"NotIPv4": "not IPv4",
		"NotIPv6": "not IPv6",

		// Json
		"NotJson": "not json",

		// Timestamp、Time
		"NotTimestamp": "not timestamp",
		"NotTime":      "not date",
		"TooEarly":     "too early",
		"TooLate":      "too late",

		// Range Distance
		"TooNear":    "too near",
		"TooFar":     "too far",
		"WrongRange": "wrong range",

		// Integer Range
		"NotIntRange":    "not int range",
		"NotInt64Range":  "not int64 range",
		"NotUintRange":   "not uint range",
		"NotUint64Range": "not uint64 range",
		"LeftTooSmall":   "left of range is too small",
		"LeftTooLarge":   "left of range is too large",
		"RightTooSmall":  "right of range is too small",
		"RightTooLarge":  "right of range is too large",

		// Timestamp/Time Range
		"NotTimestampRange": "not timestamp range",
		"NotTimeRange":      "not date range",
		"LeftTooEarly":      "left of range is too early",
		"LeftTooLate":       "left of range is too late",
		"RightTooEarly":     "right of range is too early",
		"RightTooLate":      "right of range is too late",

		// Set Count
		"TooFew":  "too few",
		"TooMany": "too many",

		// Integer Set
		"NotIntSet":    "not int set",
		"NotInt64Set":  "not int64 set",
		"NotUintSet":   "not uint set",
		"NotUint64Set": "not uint64 set",
		"ItemTooSmall": "item too small",
		"ItemTooLarge": "item too large",

		// String Set
		"NotStringSet":        "not string set",
		"ItemTooShort":        "item too short",
		"ItemTooLong":         "item too long",
		"ItemNotNumeric":      "item not numeric",
		"ItemNotDigit":        "item not digit",
		"ItemNotAlpha":        "item not alpha",
		"ItemNotAlphaNumeric": "item not alphanumeric",
		"ItemWrongFormat":     "item wrong format",

		// Email Set
		"NotEmailSet": "not email set",

		// IP Set、CIDR Set
		"NotCIDRSet":  "not cidr set",
		"NotIPSet":    "not ip set",
		"ItemNotIPv4": "item not IPv4",
		"ItemNotIPv6": "item not IPv6",
	},
	"zh_cn": map[string]string{
		// Internal error
		"InvalidValidator": "验证器非法",

		// Common
		"NotInSet":     "不在集合中",
		"ItemNotInSet": "元素不在集合中",

		// Interger
		"NotInt":    "非int型",
		"NotInt64":  "非int64型",
		"NotUint":   "非uint型",
		"NotUint64": "非uint64型",
		"TooSmall":  "太小",
		"TooLarge":  "太大",

		// String
		"NotString":       "非字符串",
		"TooShort":        "太短",
		"TooLong":         "太长",
		"WrongFormat":     "格式错误",
		"NotNumeric":      "不是数值",
		"NotDigit":        "不是纯数字组成",
		"NotAlpha":        "不是纯字母组成",
		"NotAlphaNumeric": "不是由纯字母和数字组成",

		// Email
		"NotEmail": "非邮箱地址",

		// IP、CIDR
		"NotCIDR": "非CIDR地址",
		"NotIP":   "非IP地址",
		"NotIPv4": "不是IPv4",
		"NotIPv6": "不是IPv6",

		// Json
		"NotJson": "非JSON字符串",

		// Timestamp、Time
		"NotTimestamp": "非时间戳",
		"NotTime":      "非日期",
		"TooEarly":     "太早",
		"TooLate":      "太晚",

		// Range Distance
		"TooNear":    "太近",
		"TooFar":     "太远",
		"WrongRange": "区间错误",

		// Integer Range
		"NotIntRange":    "非int区间",
		"NotInt64Range":  "非int64区间",
		"NotUintRange":   "非uint区间",
		"NotUint64Range": "非uint64区间",
		"LeftTooSmall":   "区间左值太小",
		"LeftTooLarge":   "区间左值太大",
		"RightTooSmall":  "区间右值太小",
		"RightTooLarge":  "区间右值太大",

		// Timestamp/Time Range
		"NotTimestampRange": "非时间戳区间",
		"NotTimeRange":      "非时间区间",
		"LeftTooEarly":      "区间左值太早",
		"LeftTooLate":       "区间左值太晚",
		"RightTooEarly":     "区间右值太早",
		"RightTooLate":      "区间右值太晚",

		// Set Count
		"TooFew":  "太少",
		"TooMany": "太多",

		// Integer Set
		"NotIntSet":    "非int集合",
		"NotInt64Set":  "非int64集合",
		"NotUintSet":   "非uint集合",
		"NotUint64Set": "非uint64集合",
		"ItemTooSmall": "集合中元素值太小",
		"ItemTooLarge": "集合中元素值太大",

		// String Set
		"NotStringSet":        "非字符串集合",
		"ItemTooShort":        "集合中字符串太短",
		"ItemTooLong":         "集合中字符串太长",
		"ItemNotNumeric":      "集合中字符串不是数值",
		"ItemNotDigit":        "集合中字符串不是纯数字组成",
		"ItemNotAlpha":        "集合中字符串不是纯字母组成",
		"ItemNotAlphaNumeric": "集合中字符串不是由纯字母和数字组成",
		"ItemWrongFormat":     "集合中字符串格式错误",

		// Email Set
		"NotEmailSet": "非邮箱地址集合",

		// IP Set、CIDR Set
		"NotCIDRSet":  "非CIDR集合",
		"NotIPSet":    "非IP地址集合",
		"ItemNotIPv4": "集合中的元素不是IPv4",
		"ItemNotIPv6": "集合中的元素不是IPv6",
	},
}
