## API

----------

使用 JSON 交换数据。默认当输入有错误时，返回 Status Code 422 (StatusUnprocessableEntity) 和详细错误信息。

### *注：下面所有URL均有前缀/api/v3/students*

### *POST* /me/newestWrongProblems/

获取现在仍错的错题信息

**input**

```javascript
{
    sort: number,   // 排序方式，1按出题方式，2按题目类型
    paper: number,   // 纸张大小，1 A3, 2 A4
    max: number,   // 题量
    bookPage: List<{
        book: string,  // 学习资料
        startPage: number,  // 开始页码
        endPage: number,  // 结束页码
    }>
}
```

**output**

Status Code 200:

```javascript
{
    totalNum: number, // 错题总数
    wrongProblems: List<{
        type: string, // 题目类型
        problems: List<{
            book: string,
            page: number,
            column: string,    // 栏目名称
            idx: number,    // 题目序号
            problemId: string,  // 题目识别码
            subIdx: number,  // 小问序号
            index: number,  // 这道题在这些错题当中的序号（最小为1，最大为totalNum）
            full: bool,  // 是否需要完成整一道题的所有小问
        }>,
    }>,
}
// List<T> 为T构成的列表
```

Status Code 404:

没有错题

### *POST* /me/onceWrongProblems/

获取曾经错过的所有错题信息

**input**

```javascript
{
    sort: number,   // 排序方式，1按出题方式，2按题目类型
    paper: number,   // 纸张大小，1 A3, 2 A4
    max: number,   // 题量
    bookPage: List<{
        book: string,  // 学习资料
        startPage: number,  // 开始页码
        endPage: number,  // 结束页码
    }>
}
```

**output**

Status Code 200:

```javascript
{
    totalNum: number, // 错题总数
    wrongProblems: List<{
        type: string, // 题目类型
        problems: List<{
            book: string,
            page: number,
            column: string,    // 栏目名称
            idx: number,    // 题目序号
            problemId: string,  // 题目识别码
            subIdx: number,  // 小问序号
            index: number,  // 这道题在这些错题当中的序号（最小为1，最大为totalNum）
            full: bool,  // 是否需要完成整一道题的所有小问
        }>,
    }>,
}
// List<T> 为T构成的列表
```

Status Code 404:

没有错题

### *POST* /me/getProblemsFile/

获取错题题目PDF文件URL。（与错题复习的API相同）
当full是true时，意味着只需要根据problemId获取题目，把相同problemId的元素过滤掉直到只剩一个，此时subIdx是无用信息。
当full是false时，意味着需要同时根据problemId和subIdx来获取题目，此时这个元素不需要额外处理

**input**

```javascript
{
    pageType: string,  // 纸张类型，"A3"或者"A4"
    problems: [{
        type: string,  // 类型
        // 当获取错题时排序方式是1按出题方式，得到的数据type是空字符串。
        // 则将problems对应的题目位置填入type中（因为此时problems是同一道题的不同小问，题目位置只有一个）
        // 题目位置样例： xx资料/P100/2        100是页码，2是在那本书的题目序号即idx
        problems: [{
            problemId: string,  // 题目识别码
            subIdx: number, // 小问序号
            index: number,  // 这道题在这些错题当中的序号（最小为1，最大为totalNum）
            full: bool, // 是否需要完成整一道题的所有小问
        }],
    }],
}
```

**output**

Status Code 200:

```javascript
{
    docurl: string,     // word文档URL
    pdfurl: string,     // PDF文档URL
}
```

Status Code 404:

找不到对应的题目PDF文件

### *POST* /me/getAnswersFile/

获取错题答案PDF文件URL（与错题复习的API相同）

**input**

```javascript
{
    pageType: string,    // 纸张类型，"A3"或者"A4"
    problems: [{
    	problemId: string,  // 题目识别码
        location: string,   // 题目位置
        // 题目位置样例： xx资料/P100/2        100是页码，2是在那本书的题目序号即idx
    	index: number,  // 这道题在这些错题当中的序号（最小为1，最大为totalNum）
	}]
}
```

**output**

Status Code 200:

```javascript
{
    docurl: string,     // word文档URL
    pdfurl: string,     // PDF文档URL
}
```

Status Code 404:

找不到对应的答案PDF文件

### *POST* /me/problemsRevised/

提交错题测试结果（与提交错题复习结果API相同）

**input**

```javascript
time: number,  // 完成错题复习的时间（实际即当前时间（不要仅仅精确到日）），unix时间戳
// 注意time时间戳精确到秒，即例如使用：Date.parse(new Date()) / 1000
problems: List<{
    isCorrect: bool,
    problemId: string,
    subIdx: number,
    smooth: number, // 是否顺利，顺利为1，不顺利为2
    // 若isCorrect是false，smooth为-1
    understood: number, // 是否学懂了，学懂了为1，否则2
    // 若isCorrect是true，understood为-1
}>
```

**output**

Status Code 200:

"Succeeded."
