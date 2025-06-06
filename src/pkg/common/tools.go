package common

// CheckUuId 验证是否为32位纯十六进制UUID(无连字符)
// e642aba60d874f3d9e63f77a07b0a955
func CheckUuId(uuid string) bool {
	// 快速检查长度
	if len(uuid) != 32 {
		return false
	}

	for _, r := range uuid {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}

// 支持标准UUID格式(带连字符)
// e642aba6-0d87-4f3d-9e63-f77a07b0a955
func CheckStandardUUID(uuid string) bool {
	if len(uuid) != 36 {
		return false
	}

	// 检查连字符位置
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return false
	}

	// 验证每个部分
	parts := []struct{ start, end int }{
		{0, 8}, {9, 13}, {14, 18}, {19, 23}, {24, 36},
	}

	for _, p := range parts {
		for i := p.start; i < p.end; i++ {
			c := rune(uuid[i])
			if !((c >= '0' && c <= '9') ||
				(c >= 'a' && c <= 'f') ||
				(c >= 'A' && c <= 'F')) {
				return false
			}
		}
	}

	return true
}
