package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// Validator 接口允许结构体定义自定义验证逻辑
type Validator interface {
	Validate() error
}

// 缓存相关结构
type fieldValidator struct {
	fieldIndex int
	rules      []validationRule
}

type validationRule struct {
	name  string
	param string
}

// 全局缓存和锁
var (
	validationCache = make(map[reflect.Type][]fieldValidator)
	cacheMutex      sync.RWMutex

	// 预定义错误
	errRequired   = errors.New("字段不能为空")
	errMinValue   = errors.New("字段值不能小于")
	errMaxValue   = errors.New("字段值不能大于")
	errInvalidLen = errors.New("字段长度不符合要求")
)

// Validate 验证结构体
func Validate(obj interface{}) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)

	// 处理指针
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}

	// 只处理结构体
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("验证函数需要传入结构体或结构体指针")
	}

	typ := value.Type()

	// 获取缓存的验证规则
	validators := getValidators(typ)

	// 遍历字段验证
	for _, validator := range validators {
		field := value.Field(validator.fieldIndex)
		structField := typ.Field(validator.fieldIndex)

		// 处理嵌套结构体
		if field.Kind() == reflect.Struct {
			if err := Validate(field.Interface()); err != nil {
				return fmt.Errorf("%s.%s: %w", typ.Name(), structField.Name, err)
			}
			continue
		}

		// 处理嵌套结构体指针
		if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Struct {
			if !field.IsNil() {
				if err := Validate(field.Elem().Interface()); err != nil {
					return fmt.Errorf("%s.%s: %w", typ.Name(), structField.Name, err)
				}
			}
			continue
		}

		// 验证字段
		for _, rule := range validator.rules {
			if err := validateFieldByRule(field, rule); err != nil {
				return fmt.Errorf("%s.%s: %w", typ.Name(), structField.Name, err)
			}
		}
	}

	// 检查自定义验证器
	if validator, ok := obj.(Validator); ok {
		return validator.Validate()
	}

	return nil
}

// 从缓存获取验证规则
func getValidators(t reflect.Type) []fieldValidator {
	cacheMutex.RLock()
	validators, exists := validationCache[t]
	cacheMutex.RUnlock()

	if !exists {
		cacheMutex.Lock()
		defer cacheMutex.Unlock()

		// 双重检查
		if validators, exists = validationCache[t]; !exists {
			initValidationCache(t)
			validators = validationCache[t]
		}
	}

	return validators
}

// 初始化验证规则缓存
func initValidationCache(t reflect.Type) {
	var validators []fieldValidator

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules := parseValidationRules(tag)
		if len(rules) > 0 {
			validators = append(validators, fieldValidator{
				fieldIndex: i,
				rules:      rules,
			})
		}
	}

	validationCache[t] = validators
}

// 解析验证规则 - 修复版本
func parseValidationRules(tag string) []validationRule {
	var rules []validationRule
	// 使用状态机解析带引号的字符串
	var buf strings.Builder
	inQuotes := false
	var quoteChar byte
	escape := false

	for i := 0; i < len(tag); i++ {
		c := tag[i]

		switch {
		case escape:
			buf.WriteByte(c)
			escape = false
		case c == '\\':
			escape = true
		case inQuotes:
			if c == quoteChar {
				inQuotes = false
			}
			buf.WriteByte(c)
		case c == '"' || c == '\'':
			inQuotes = true
			quoteChar = c
			buf.WriteByte(c)
		case c == ',' && !inQuotes:
			// 遇到逗号且不在引号内，分割规则
			if buf.Len() > 0 {
				rules = appendRule(rules, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteByte(c)
		}
	}

	// 处理最后一个规则
	if buf.Len() > 0 {
		rules = appendRule(rules, buf.String())
	}

	return rules
}

// 添加解析后的规则
func appendRule(rules []validationRule, s string) []validationRule {
	s = strings.TrimSpace(s)
	if s == "" {
		return rules
	}

	// 分割规则名和参数
	parts := strings.SplitN(s, "=", 2)
	rule := validationRule{name: strings.TrimSpace(parts[0])}

	if len(parts) > 1 {
		rule.param = strings.TrimSpace(parts[1])
		// 去除参数两端的引号
		if len(rule.param) >= 2 {
			if (rule.param[0] == '"' && rule.param[len(rule.param)-1] == '"') ||
				(rule.param[0] == '\'' && rule.param[len(rule.param)-1] == '\'') {
				rule.param = rule.param[1 : len(rule.param)-1]
			}
		}
	}

	return append(rules, rule)
}

// 根据规则验证字段
func validateFieldByRule(field reflect.Value, rule validationRule) error {
	if !field.IsValid() || !field.CanInterface() {
		return nil
	}

	switch rule.name {
	case "required":
		if isEmpty(field) {
			return errRequired
		}

	case "min":
		param, err := strconv.ParseFloat(rule.param, 64)
		if err != nil {
			return fmt.Errorf("无效的最小值参数: %s", rule.param)
		}

		if num, ok := getNumericValue(field); ok {
			if num < param {
				return fmt.Errorf("%w %v", errMinValue, param)
			}
		}

	case "max":
		param, err := strconv.ParseFloat(rule.param, 64)
		if err != nil {
			return fmt.Errorf("无效的最大值参数: %s", rule.param)
		}

		if num, ok := getNumericValue(field); ok {
			if num > param {
				return fmt.Errorf("%w %v", errMaxValue, param)
			}
		}

	case "len":
		param, err := strconv.Atoi(rule.param)
		if err != nil {
			return fmt.Errorf("无效的长度参数: %s", rule.param)
		}

		if field.Kind() == reflect.String || field.Kind() == reflect.Slice || field.Kind() == reflect.Map {
			if field.Len() != param {
				return fmt.Errorf("%w %v", errInvalidLen, param)
			}
		}
	}

	return nil
}

// 判断字段是否为空
func isEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		return value.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Ptr, reflect.Interface:
		return value.IsNil()
	default:
		return false
	}
}

// 获取数值类型的值
func getNumericValue(value reflect.Value) (float64, bool) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(value.Uint()), true
	case reflect.Float32, reflect.Float64:
		return value.Float(), true
	default:
		return 0, false
	}
}
