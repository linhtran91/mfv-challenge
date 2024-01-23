package services

import (
	"mfv-challenge/internal/constants"
	"net/url"
	"strconv"
)

func getLimitOffset(values url.Values) (int, int) {
	page := getValueFromUrl(values, "page", constants.DefaultPage)
	size := getValueFromUrl(values, "size", constants.DefaultSize)
	if size >= constants.MaximumSize {
		size = constants.DefaultSize
	}
	if page == 1 {
		return size, 0
	}

	return size, (page-1)*size + 1
}

func getValueFromUrl(values url.Values, key string, defaultVal int) int {
	val := values.Get(key)
	if val == "" {
		return defaultVal
	}
	t, err := strconv.Atoi(val)
	if err != nil || t == 0 {
		return constants.DefaultPage
	}
	return t

}
