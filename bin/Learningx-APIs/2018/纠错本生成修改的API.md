## API

使用 JSON 交换数据。默认当输入有错误时，返回 Status Code 422 (StatusUnprocessableEntity) 和详细错误信息。

## 修改：

#### 注：下面所有URL均有前缀/api/v3/students

### *GET* /me/books/

获取有哪些书本

**no input**

**output**

```javascript
[{
    name: string,	// 书本名称
    bookID: string,    // 识别码
	type: number,    // 书本资料类型，1：课本,2：普通辅导书，3：培优资料
}]
// [] 代表列表
```

### *POST* /me/problems/

提交题目录入信息（以前的平时作业、试卷、纠错本录入都改为这个接口）

**input**

```javascript
{
    time: number,  // unix时间戳
    // 如果是平时作业、课本、试卷录入（即界面有选择日期的情况下），这里是选择的日期
    // 如果是纠错本录入（界面没有选择日期），这里是当前时间（不要仅仅精确到日）），
	// 注意time时间戳精确到秒，即例如使用：Date.parse(new Date()) / 1000
    type: number, // 1：课本,2：普通辅导书，3：培优资料,（前三个在获取书本的接口中有） 4: 试卷， 5： 纠错本
	problems: [{
    	isCorrect: bool,
    	problemId: string,
    	subIdx: number,
		// 原先的smooth、understood删去了
    }]
}
```

**output**

Status Code 200:

"Succeeded."

#### 注：下面所有URL均有前缀/api/v3/staffs

### *GET* /classes/books/

获取某个班级的书本信息

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: number,    // 班级
```

**output**

```javascript
[{
    name: string,	// 书本名称
    bookID: string,    // 识别码
    type: number,    // 书本资料类型，1：课本,2：普通辅导书，3：培优资料
}]
// [] 代表列表
```

### *POST* /students/\<learnID\>/problems/

提交错题录入信息

**input**

```javascript
{
    time: number,  // 作业布置时间，unix时间戳  使用当前时间即可
	// 注意time时间戳精确到秒，即例如使用：Date.parse(new Date()) / 1000
    type: number, // 1：课本,2：普通辅导书，3：培优资料,（前三个在获取书本的接口中有） 4: 试卷， 5： 纠错本
	problems: [{
	    isCorrect: bool,
	    problemId: string,
	    subIdx: number,
	}]
}
```

**output**

```javascript
"successfully uploaded"
```

### 