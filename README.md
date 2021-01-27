# ahocorasick
An extremely fast implementation of Aho Corasick algorithm based on Double Array Trie.Can be used for fuzzy matching, prohibited word marking.

[中文](https://github.com/TheFutureIsOurs/ahocorasick/blob/master/README-ZH.md)

Common usage scenarios are:

> 1.When chatting, block when the user-entered
qury hits a prohibited word.  

> 2.Determine if the user input word hits the prohibited word blacklist,and then follow up with logic.  

> 3.other scenarios that require high-performance fuzzy matching.  

### Characteristic

In addition to the construction based on double arrays, the algorithm also has additional optimization to ensure higher efficiency

> 1.When constructing the algorithm, the double array and the fail pointer and output items are built at the same time, greatly increasing the speed of construction and reducing the time it takes to build a dictionary.  

> 2.Compress the output items to reduce build time and memory useage and gc.  

> 3.Compressed leaf nodes, reduced memory consumption, improve retrieval efficiency.  

Explanation of compressed output items: 

e.g. the prohibited word dictionary is:

	he 
	she

The input is:  ushe,then the returned hit item is: she 

Explanation: 

If we don't do output compression, the hit item is:he,she;  

Considering that in practice, whether it is blocking prohibited words (e.g., ushe above, blockingprohibited words afteru,) 
or just determining whether a prohibited word has been hit,

there is no need for redundant output items, 

compression of output items, not only can reduce the build time, but also reduce memory consumption, 

for go, can also reduce gc scanning.  


### Performance Compare 

 A popular project [cloudflare/ahocorasick](https://github.com/cloudflare/ahocorasick) and 
 return the same results (using api:MultiPatternIndexes(content) (see test file forcomparison code).  

Dictionary:
> dictionary.txt
> 
> Number of strings: 153151. 
> 
> Number of characters: 401552 (average 2.6 characters per article)

File to be retrieved: 

> text.txt
> 
> Characters: 815006
> 

Dictionary and file to be retrieved are already in the project

Running machine: 

Lenovo Small New Pro13 Ryzen 5 3550H 

Match the file to be retrieved 100 times

go Version: 1.15 

 gc before matching (runtime. Gc())

| Warehouse   |Dictionary BuildTime(ms)| One gc Time(ms)|  100 full-text matches(ms)  |inuse_space|inuse_objects|
| --------                   |-----:| -----:  | :----:  | :----: |:----:|
| cloudflare/ahocorasick     |59| 686  |   14910     |4.67G|  360455|
| TheFutureIsOurs/ahocorasick|2431| 0   |   5341       |14.2M|  4  |

As you can see, performance is 64% faster than cloudflare/ahocorasick (retrieving 81.5w characters takes only 54ms), thanks to double array implementations and optimization of output items, the time to perform a gc can be ignored (only four slice headers are held). It also uses only 14.2M of memory.  

Thanks to the various optimizations above, the build time is optimized to 2.4s, although it is still longer than cloudflare/ahocorasick, 

but considering that the build time is more complex, and most applications only need to be built once at startup, the second stage for the sensory time is short, the build time is controllable (the author with a business 61w line black word test, build time 3.5s), but the higher quantifiable indicators of the return, this is worth it.  




# How to use 

### Download

> go get -u github.com/TheFutureIsOurs/ahocorasick


### Useage

import "github.com/TheFutureIsOurs/ahocorasick"

First step: construct an algorithm.  

It can be constructed from a list of strings:

```go

dictionary := []string{"hers", "his", "she", "he"}

ac, err := ahocorasick.Build(dictionary)

```

or by file construction:

dictionary.txt file content is: 

	hers
	his
	she
	he
 
```go

ac, err := BuildFromFile("./dictionary.txt")

```

Step 2: Match.  

api list:

MultiPatternSearch(content []rune) []Hit

> returns all hit strings at the beginning and end of the original string andstrings.  

MultiPatternIndexes(content []rune) []int

> returns the starting point of all hit strings in the original character.  

MultiPatternHit(content []rune) bool

> returns whether the content hit the dictionary. If hit, true is returned immediately and the query is not continued.  

A complete example is as follows:

```go

// Constructs dictionary from the list of 

dictionary := []string{"hers", "his", "she", "he"}

ac, err := ahocorasick.Build(dictionary)

search := ac.MultiPatternSearch([]rune("ushers")) // returns the longest string of all hits

/*

// Returns a list of all strings hit as follows:
// The starting position of the string at the beginning and end of the original string and the hit string
// The startingvalue of the original string is 0
1	3	she
2	5	hers

*/

for _, v := range search {

    fmt.Printf("%d\t%d\t%s\n", v.Begin, v.End, string(v.Value))

}

```


### Thanks 

The double array trie was inspired by the open source project of the [darts-java](https://github.com/komiya-atsushi/darts-java)



