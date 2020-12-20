# ahocorasick
Aho Corasick算法基于双数组Trie树的一种极速实现。可用于模糊匹配，违禁词标记。

注意：该算法为了更高的性能，针对输出的构造做了优化，仅会返回最长命中项。

如：违禁词词典为：
	
	he
	she

输入为: ushe，则仅返回的最长命中为：she.(因为在she结束时，本应该输出she和he，但是由于she比he更长，故只返回she。)


常用使用场景为：
>1.聊天时，对用户输入的query进行命中违禁词时进行屏蔽。
>2.判断用户输入词是否命中了违禁词黑名单，然后就行后续逻辑处理。

等需要高性能模糊匹配的场景。

### 优化

本算法基于双数组进行构建，为了更高的效率，又做了额外的优化，几个大的优化点：

1.在构造算法时，双数组和fail指针及输出项会同时构建，确保更快的构建速度，减少构建时间。

2.压缩输出项，减少构建时长及内存占用及gc。

3.压缩了叶子节点，减少内存占用，提升检索效率。

### 性能

对比一个star数较多的[cloudflare/ahocorasick](https://github.com/cloudflare/ahocorasick)
并返回相同的结果（使用api:MultiPatternIndexes(content []rune) []int）进行对比（对比代码见test文件）。

字典：dictionary.txt 字符串个数：153151。字符数：401552 （平均每条2.6个字符）

待检索文件：text.txt 字符数：815006

（字典和待检索文件已在仓库内）

运行机器：联想小新pro13 Ryzen 5 3550H

对待检索文件匹配100遍

go版本：1.15

待匹配前先进行一次gc(runtime.Gc())

| 仓库                       |字典构建时间(ms)| 一次gc时间(ms)|  100遍全文匹配(ms)  |inuse_space|inuse_objects|
| --------                   |-----:| -----:  | :----:  | :----: |:----:|
| cloudflare/ahocorasick     |59| 686  |   14910     |4.67G|  360455|
| TheFutureIsOurs/ahocorasick|2431| 0   |   5341       |14.2M|  4  |

可以看出，性能比cloudflare/ahocorasick快64%（检索一遍81.5w的字符仅需54ms）,得益于双数组实现和对输出项的优化，执行一次gc的时间可以忽略不记(持有的指针仅为四个数组header)。另外占用内存仅为14.2M。

但得益于上面的各种优化，使得构建时间优化到2.4s，虽然仍较cloudflare/ahocorasick长，但是考虑到构建时较复杂，而且大部分应用只需在启动时构建一次，秒级对于感官时间较短，构建时长可控（我的一个真实业务黑名单词61w行，构建时间3.5s），但带来的更高的可量化的各项指标收益，这个是值得的。




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

search := ac.MultiPatternSearch([]rune("ushers")) // 会返回所有命中的最长字符串

/*
// 返回命中的所有字符串列表如下。
// 字符串在原字符串开始的位置&结束位置&命中的字符串
// 原字符串起始值为0
1	3	she
2	5	hers
*/

for _, v := range search {

    fmt.Printf("%d\t%d\t%s\n", v.Begin, v.End, string(v.Value))

}

```







