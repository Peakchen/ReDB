## API

----------

使用 JSON 交换数据。默认当输入有错误时，返回 Status Code 422 (StatusUnprocessableEntity) 和详细错误信息。

### *注：下面所有URL均有前缀/api/contactWithUserServer-v1*

### *POST* /getProblemsFile/

获取错题题目PDF文件URL：
如果full是true，根据problemId提取题目，给题目加上index对应的题目序号。
如果full是false，根据problemId提取题目，给题目加上index对应的题目序号后，插入：“仅需要做第几小问”。
然后拼接并转换成PDF文件，返回文件URL

**input**

```javascript
// 原先：List<{
//     type: string,  // 题目类型
//     problems: List<{
//         problemId: string,  // 题目识别码
//         subIdx: number, // 小问序号
//         index: number,  // 这道题在这些错题当中的序号（最小为1，最大为totalNum）
//         full: bool, // 是否需要完成整一道题的所有小问
//     }>,
// }>,

// 改为：
{
    school: string,  // 学校
    grade: int,  // 年级
    classID: int,  // 班级号
    name: string,  // 名字
    problems: List<{
        type: string,  // 题目类型
        problems: List<{
            problemId: string,  // 题目识别码
            subIdx: number, // 小问序号
            index: number,  // 这道题在这些错题当中的序号（最小为1，最大为totalNum）
            full: bool, // 是否需要完成整一道题的所有小问
        }>,
    }>,
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

### *POST* /getAnswersFile/

获取错题答案PDF文件URL：根据problemId提取题目对应的答案，给题目加上index对应的题目序号，然后拼接并转换成PDF文件，返回文件URL

**input**

```javascript
// 原来：List<{
//     problemId: string,  // 题目识别码
//     index: number,  // 这道题在这些错题当中的序号（最小为1，最大为totalNum）
// }>,

// 改为：
{
    school: string,  // 学校
    grade: int,  // 年级
    classID: int,  // 班级号
    name: string,  // 名字
    problems: List<{
        problemId: string,  // 题目识别码
        index: number,  // 这道题在这些错题当中的序号（最小为1，最大为totalNum）
    }>,
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

### *POST* /getPointsFile/

获取错题知识点PDF文件URL：
得到该章节所有知识点（注：传过来的题目应该是处于同一章节中的），
根据problemId提取题目并给题目加上index对应的题目序号，然后在文档第一页加上题目序号、小问序号与知识点的对应表，转换成PDF文件，返回文件URL

**input**

```javascript
// 原来：List<{
//     problemId: string,  // 题目识别码
//     subIdx: number, // 小问序号
//     index: number,  // 这道题在这些错题当中的序号（最小为1，最大为totalNum）
// }>,

// 改为：
{
    school: string,  // 学校
    grade: int,  // 年级
    classID: int,  // 班级号
    name: string,  // 名字
    problems: List<{
        problemId: string,  // 题目识别码
        subIdx: number, // 小问序号
        index: number,  // 这道题在这些错题当中的序号（最小为1，最大为totalNum）
    }>,
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

找不到对应的知识点PDF文件

