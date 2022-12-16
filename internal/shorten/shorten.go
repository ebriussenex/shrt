package shorten

import (
	"net/url"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const alphLen = uint32(len(alphabet))

func Shorten(id uint32) string {
	var builder strings.Builder
	nums := make([]uint32, 0)

	for id > 0 {
		n := id % alphLen
		nums = append(nums, n)
		id = id / alphLen
	}

	reverse(nums)

	for _, num := range nums {
		builder.WriteString(string(alphabet[num]))
	}

	return builder.String()
}

func reverse(nums []uint32) {
	j := len(nums) - 1
	for i := 0; i < j; i++ {
		nums[i], nums[j] = nums[j], nums[i]
		j--
	}
}

func CreateResUrl(baseUrl, id string) (string, error) {
	parsed, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}
	parsed.Path = id
	return parsed.String(), nil

}
