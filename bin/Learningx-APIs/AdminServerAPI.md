## API

注：下面所有URL均有前缀/api/admin-v1 

----------

使用 JSON 交换数据。默认当输入有错误时，返回 Status Code 422 (StatusUnprocessableEntity) 和详细错误信息。

### *GET* /info/

获取具体信息（大类、章、节、知识点）

**filters**

```
block: number,  // 大类
chapter: number,  // 章
section: number,  // 节
point: number,  // 知识点
```

当某个filter值为1，代表需要该filter对应的信息，不为1则该filter返回的信息为无效信息。

**output**

Status Code 200:

```
List<{
	block: string,
	chapter: number,
	chapterName: string,  // 章名称
	section: number,
	sectionName: string,  // 节名称
	point: string,
}>

//此处List<T>表示T组成的list，下同，不再说明
```

### *GET* /examplesMeta/

获取总体分析结果（例题与题型信息）

**filters**

```
chapter: number,  // -2代表全部
section: number,  // -2代表全部

//当chapter为-2（全部章）时，section必须为-2
```

**output**

Status Code 200:

```
List<{
	exampleNum: number,  // 例题数量
	typeNum: number,  // 题型数量
}>
```

status Code 404:

该条件下没有任何题型。

### *GET* /detailedExamplesMeta/

获取给定例题个数下的详细分析结果，该接口默认chapter、section均为全部。

**filters**

```
exampleNum: number,
```

**output**

Status Code 200:

```
List<{
	chapter: number,
	typeNum: number,  // 题型数量
}>

// 例题个数为传给服务器的exampleNum
```

### *GET* /typeWithoutExample/

获取给定章节下的没有例题的题型信息。

**filters**

```
chapter: number,
section: number,
```

**output**

Status Code 200:

```javascript
{
	typeName: string,  // 题型名称
	problems: List<{
		problemId: string,  // 题目识别码
		subIdx: number,  // 小问序号
		image: URL,  // 题目对应的图片的URL
	}>,  // problems: 这个题型下的所有习题
}
```

Status Code 404:

这个章节对应的所有题型都已经存在例题了。

### *GET* /answerImage/

获取某个题目的答案URL

**filters**

```
problemId: string,
subIdx: number,
```

**output**

Status Code 200:

```
{
	url: string,  // 答案URL
}
```

Status Code 404:

找不到答案图片。

### *POST* /examples/

新增一道例题。

**input**

```javascript
{
	typeName: string,
	problemId: string,
	subIdx: number,
}
```

**output**

Status Code 200:

添加成功。

### *GET* /problemNotKnowChapter/

获取一道还没有确定章的习题。

**no filters**

**output**

Status Code 200:

```javascript
{
	problemId: string,
	image: string,  // 图案URL
}
```

Status Code 404:

没有还没确定章的题目了。

### *PUT* /problemChapter/

更新题目章信息。

**input**

```
{
	problemId: string,
	chapter: number,
}

// 当chapter为-1时，代表撤销该章信息。
```

**output**

Status Code 200:

更新成功。

### *PUT* /problemSection/

更新题目节信息。

**input**

```
{
	problemId: string,
	section: number,
}

// 当section为-1时，代表撤销该节信息。
```

**output**

Status Code 200:

更新成功。

### *GET* /problemNotKnowSection/

获取某一章下一道还没有确定节的习题。

**filters**

```
chapter: number
```

**output**

Status Code 200:

```javascript
{
	problemId: string,
	image: string,  // 图案URL
}
```

Status Code 404:

没有还没确定节的题目了。

### *GET* /problemNotKnowType/

在某一章节下获取一道还没有确定题型的习题。

**filters**

```
chapter: number,
section: number,
```

**output**

Status Code 200:

```
{
	problemId: string,
	subIdx: number,
	image: string,  // 图案URL
	examples: List<{
		problemId: string,
		subIdx: number,
		image: string,
	}>,
}
```

Status Code 404:

当前条件下没有还没确定题型的题目了。

### *PUT* /subProblemType/

更新小问题型信息。

**input**

```
{
	problemId: string,
	subIdx: number,
	example: {
		problemId: string,
		subIdx: number,
	},  // 例题信息
}
```

**ouput**

Status Code 200:

更新成功。

### *PUT* /subProblemSimilarity/

更新小问相似度信息。

**input**

```
{
	problemId: string,
	subIdx: number,
	similarity: number,  //相似度
}
```

**ouput**

Status Code 200:

更新成功。

### *GET* /problemNotKnowSimilarity/

在某一章节下获取一道还没有确定相似度的习题。

**filters**

```
chapter: number,
section: number,
```

**output**

Status Code 200:

```
{
	problemId: string,
	subIdx: number,
	image: string,  // 图案URL
	examples: List<{
		problemId: string,
		subIdx: number,
		image: string,
	}>,
}
```

### *GET* /similarityInfo/

获取相似度统计信息。

**filters**

```
chapter: number,
// chapter为-2代表全部
```

**output**

Status Code 200:

```
{
	totalNum: number,  // 总题目数量
	notFormalNum: number,  // 未正式确定题型的题目数量
	bothFourOrAbove: number,  // 评分均为4或以上的题目数量
	threeAndFour: number,
	threeAndFive: number,
	bothThreeOrBelow: number,
	notCompleted: number,
}
```

### *POST* /confirmSimilarity/

**no filters**

**output**

Status Code 200:

确认完成。
