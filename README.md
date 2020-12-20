# ahocorasick
Aho Corasick算法基于双数组Trie树的一种极速实现。可用于全文匹配，模糊匹配等。


# 如何使用

### 下载

>go get -u github.com/TheFutureIsOurs/ahocorasick

### 使用

import "github.com/TheFutureIsOurs/ahocorasick"

第一步：构造算法。

可以通过字符串列表构造：

```go

dictionary := []string{"hers", "his", "she", "he"}

ac, err := ahocorasick.Build(dictionary)

```
或通过文件构造:

dictionary.txt文件内容为：

	hers
	his
	she
	he

```go

ac, err := BuildFromFile("./dictionary.txt")

```

第二步：进行匹配。

api列表：

MultiPatternSearch(content []rune) []Hit 返回所有命中的字符串在原字符串中的起终点及字符串。

MultiPatternIndexes(content []rune) []int 返回所有命中的字符串在原字符中的起点。

MultiPatternHit(content []rune) bool 返回content是否命中了字典。如果命中了，会立即返回true，不会再继续查询下去。

一个完整的例子如下：

```go

// 从字符串列表构造

dictionary := []string{"hers", "his", "she", "he"}

ac, err := ahocorasick.Build(dictionary)

search := ac.MultiPatternSearch([]rune("ushers")) // 会返回所有命中的字符串

/*
// 返回命中的所有字符串列表如下。
// 字符串在原字符串开始的位置&结束位置&命中的字符串
// 原字符串起始值为0
1	3	she
2	3	he
2	5	hers
*/

for _, v := range search {

    fmt.Printf("%d\t%d\t%s\n", v.Begin, v.End, string(v.Value))

}

```

性能：

对比一个star数较多的[cloudflare/ahocorasick](https://github.com/cloudflare/ahocorasick)

字典：dictionary.txt 字符串个数：153151。字符数：401552

待检索文件：text.txt 字符数：815006

（字典和待检索文件已在仓库内）

运行机器：联想小新pro13 Ryzen 5 3550H

对待检索文件匹配100遍

go版本：1.15

待匹配前先进行一次gc(runtime.Gc())

| 仓库                       | 一次gc时间(ms)|  100遍全文匹配(ms)  |inuse_space|
| --------                   | -----:  | :----:  | :----: |
| cloudflare/ahocorasick     | 686  |   14910     |4.67G|
| TheFutureIsOurs/ahocorasick| 15   |   10334   |31.78M|

可以看出，性能比cloudflare/ahocorasick 快30%（检索一遍81.5w的字符仅需103ms）,得益于双数组实现，执行一次gc的时间大幅下降。另外占用内存仅为31.78M.








