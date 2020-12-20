/*
 * @Author: Daiming Liu (xingrufeng)
 */
package ahocorasick

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func readFile(input string) []rune {
	file, err := os.Open(input)
	if err != nil {
		return nil
	}
	defer file.Close()
	content := make([]rune, 0)
	bufReader := bufio.NewReader(file)
	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			break
		}
		keyword := strings.TrimSpace(string(line))
		if keyword == "" {
			continue
		}
		content = append(content, []rune(keyword)...)
	}
	return content

}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func writeFile(filename string, contents []Hit) {
	var f *os.File
	var err1 error
	if checkFileIsExist(filename) { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err1 = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
	}
	defer f.Close()
	if err1 != nil {
		panic(err1)
	}
	for _, v := range contents {
		str := strconv.Itoa(v.Begin) + "\t" + strconv.Itoa(v.End) + "\t" + string(v.Value) + "\n"
		io.WriteString(f, str)
	}

}

func TestBuild(t *testing.T) {
	start := time.Now().UnixNano()
	ac, err := BuildFromFile("./dictionary.txt")
	//ac, err := BuildFromFile("./black.txt")
	fmt.Println(err)
	runTime := (time.Now().UnixNano() - start) / 1000 / 1000
	fmt.Println("字典加载时间(ms)", runTime)
	start = time.Now().UnixNano()
	content := readFile("./text.txt")
	runTime = (time.Now().UnixNano() - start) / 1000 / 1000
	fmt.Println("测试文件加载时间(ms)", runTime)
	start = time.Now().UnixNano()
	ac.MultiPatternIndexes(content)
	/*
		search := ac.MultiPatternSearch([]rune("一群"))
		for _, v := range search {
			fmt.Printf("%d\t%d\t%s\n", v.Begin, v.End, string(v.Value))
		}
	*/
	runTime = (time.Now().UnixNano() - start) / 1000 / 1000
	fmt.Println("检索时间(ms)", runTime)

	//writeFile("result2", result)

}

func TestAb(t *testing.T) {
	/*
		kws := []string{
			"hers", "his", "she", "he",
		}
	*/
	/*
		kws := []string{
			"中华人民共和国", "中华人民", "人民共和国", "中华人民",
		}
	*/
	kws := []string{
		"一", "群", "一群羊",
	}
	ac, err := Build(kws)
	//ac, err := Build(kws)
	//ac, err := BuildFromFile("./dictionary.txt")
	if err != nil {
		fmt.Println(err)
	}
	//search := ac.MultiPatternSearch([]rune("中华人民共和国"))
	//search := ac.MultiPatternSearch([]rune("ushers"))
	search := ac.MultiPatternSearch([]rune("一群"))
	for _, v := range search {
		fmt.Printf("%d\t%d\t%s\n", v.Begin, v.End, string(v.Value))
	}
}

func TestWord(t *testing.T) {
	fmt.Println([]rune("一"), []rune("群"))
}
