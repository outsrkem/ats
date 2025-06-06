package validate

import (
	"fmt"
	"testing"
)

// 示例结构体
type User struct {
	Name  string `validate:"required"`
	Age   int    `validate:"required,min=18,max=100"`
	Email string `validate:"required,len=10"`
}

func TestValidation(t *testing.T) {
	// 测试用例1：有效的数据
	user := User{
		Name:  "张三",
		Age:   25,
		Email: "1234567890", // 10个字符
	}
	fmt.Println(Validate(user))

	// 测试用例2：年龄小于最小值
	invalidUser := User{
		Name:  "李四",
		Age:   16,
		Email: "1234567890",
	}
	fmt.Println(Validate(invalidUser))
}
