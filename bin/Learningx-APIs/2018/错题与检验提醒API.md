## API

使用 JSON 交换数据。默认当输入有错误时，返回 Status Code 422 (StatusUnprocessableEntity) 和详细错误信息。

### *注：下面所有URL均有前缀/api/v3/students*

### *GET* /me/uploadTasks/

获取所有没有完成的上传标记结果的任务

**input**

no input

**output**

Status Code 200:

```javascript
[{
	time: number,  // 任务发生的时间对应的unix时间戳
	type: number,  // 任务类型（没标记的是错题为1，没标记的是检验题为2）
}]
```

Status Code 404:

```
"No upload tasks."
```

没有未完成的上传任务

### *GET* /me/uploadTasks/\<time\>/

获取任务内容，\<time\>是任务发生的时间对应的时间戳。返回上传标记任务当时获取的数据。

**input**

no input

**output**

Status Code 200:

```javascript
{
    totalNum: number, // 错题总数
    wrongProblems: [{
        type: string, // 题目类型
        problems: [{
            book: string,
            page: number,
            column: string,    // 栏目名称
            idx: number,    // 题目序号
            problemId: string,  // 题目识别码
            subIdx: number,  // 小问序号
            index: number,  // 这道题在这些错题当中的序号（最小为1，最大为totalNum）
            full: bool,  // 是否需要完成整一道题的所有小问
        }],
    }],
}
```

Status Code 404:

```
"Can not find the task."
```

找不到这个时间戳对应的任务

### *DELETE* /me/uploadTasks/\<time\>/

删除某一个任务。\<time\>是任务发生的时间对应的时间戳。

**input**

 no input

**output**

Status Code 200:

```
"Successfully deleted"
```

Status Code 404:

```
"Can not find the task.
```

找不到这个时间戳对应的任务。

### *POST* /me/uploadTasks/

新增一个上传任务，并保存任务细节。

例如：上传了如下数据
{
    time: 123456,
    type: 1,
    detail: "{'aaa': 'bbb'}"
}

则访问GET /me/uploadTasks/123456/将返回
{
    "aaa": "bbb"
}

建议将wrongProblem获取生成纠错本的错题得到的数据，用JSON.stringify()编码之后放进detail传进来

**input**

```javascript
{
    time: number,   // 任务对应的unix时间戳（即当前的时间戳）（时间精确到秒）
    type: number,   // 任务类型（没标记的是错题为1，没标记的是检验题为2）
    detail: string,    // 获取任务时需要获取的json数据对应的编码字符串
}
```

**output**

Status Code 200:

```
"Successfully create a task."
```