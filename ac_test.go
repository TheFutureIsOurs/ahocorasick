/*
 * @Author: Daiming Liu (xingrufeng)
 * @
 */
/*
 * @Author: Daiming Liu (xingrufeng)
 */
package ahocorasick

import (
	"fmt"
	"testing"
)

func TestAb(t *testing.T) {
	kws := map[string]string{
		"hers": "hers",
		"his":  "his",
		"she":  "she",
		"he":   "he",
	}
	ac, err := Build(kws)
	if err != nil {
		fmt.Println(err)
	}
	ac.MultiPatternSearch([]rune("ushers hers"))
}
