## API

----------

使用 JSON 交换数据。默认当输入有错误时，返回 Status Code 422 (StatusUnprocessableEntity) 和详细错误信息。

### *POST* /api/v3/login/

登录

**input**

```javascript
learnId: number,  // 学习号
password: string,  // 密码
remember: boolean,  // 记住我
```

**output**

Status Code 200 No Content

### *注：下面所有URL均有前缀/api/v3/students*

### *GET* /me/problems/

获取有哪些题目（这里只获取该用户没有布置过的题目）

**filters**

```javascript
book: string,  // 资料名称
page: number,  // 页码
```

**output**

```javascript
List<{
    lessonName, string,  // 课时名称
    column: string,  // 栏目名称
    idx: number,  // 题目序号
    subIdx, number,  // 小问序号
    problemId, string,  // 题目识别码
}>

// List<T>代表T类型构成的列表，以下不再说明
```

Status Code 404:

提供的书本与页码下没有对应信息

### *POST* /me/logout/

登出

**no filters**

**output**

Status Code 200 No Content

### *PUT* /me/password/

修改密码

**input**

```javascript
{
    password: string,  // 新密码
}
```

**output**

Status Code 200 No Content

### *GET* /me/profile/

获取用户信息

**no filters**

**output**

```javascript
{
    learnId: number,  // 学习号
    school: string,  // 学校
    grade: string,  // 年级
    classId: number,  // 在该学校中的班级号
    realName: string,  // 真实姓名
    gender: string,  // 性别
    telephone: string,  // 电话
}
```

### *PATCH* /me/profile/

更新用户信息

**input**

```javascript
{
    school: string,  // 学校
    grade: string,  // 年级
    classId: number,  // 在该学校中的班级号
    realName: string,  // 真实姓名
    gender: string,  // 性别
    telephone: string,  // 电话
}
```

**output**

Status Code 200 No Content

### *GET* /me/books/

获取有哪些书本

**no filters**

**output**

```javascript
List<string>    // 直接是书本名字字符串构成的数组
```

### *POST* /me/problems/

提交录入信息

**input**

```javascript
time: number,  // 作业布置时间，unix时间戳
// 注意time时间戳精确到秒，即例如使用：Date.parse(new Date()) / 1000
problems: List<{
    isCorrect: bool,
    problemId: string,
    subIdx: number,
}>
```

**output**

Status Code 200 No Content

### *GET* /me/info/

获取具体信息（大类、章、节、知识点）

**filters**

```javascript
block: number,  // 大类
chapter: number,  // 章
section: number,  // 节
point: number,  // 知识点

// 当某个filter值为1，代表需要该filter对应的信息，为0则该filter返回的信息为无效信息。
```

**output**

Status Code 200:

```javascript
List<{
	block: string,
	chapter: number,
	chapterName: string,  // 章名称
	section: number,
	sectionName: string,  // 节名称
	point: string,
}>
```

### *GET* /me/problemsSortByTime/

获取按时间分组过的题目信息

**filters**

```javascript
chapter: number,
section: number,
```

**output**

```javascript
List<{
    time: number,  // 布置时间：unix时间戳
    problems: List<{
        book: string,  // 资料名称
        page: number,  // 页码
        lessonName, string,  // 课时名称
        column: string,  // 栏目名称
        idx: number,  // 题目序号
        subIdx: number,  // 小问序号
        isCorrect: number,  // 正确与否: 1.最新做对了，以前没错过，2.最新做对了，以前错过，3.最新做错了
        category: string,  // 题型大类
    }>,
}>

// 注意顺序已经确定，按顺序展示，不需要变动
```

### *GET* /me/problemsSortByType/

获取按题型分组过的题目信息

**filters**

```javascript
chapter: number,
section: number,
```

**output**

```javascript
List<{
    type: string,  // 题型
    category: string,  // 大类
    problems: List<{
        book: string,  // 资料名称
        page: number,  // 页码
        lessonName, string,  // 课时名称
        column: string,  // 栏目名称
        idx: number,  // 题目序号
        subIdx: number,  // 小问序号
        isCorrect: number,  // 正确与否: 1.最新做对了，以前没错过，2.最新做对了，以前错过，3.最新做错了
        assignDate: number,  // 作业布置时间
    }>,
}>

// 注意顺序已经确定，按顺序展示，不需要变动
```

### *GET* /me/problemsSortByAccuracy/

获取按掌握程度分组过的题目信息

**filters**

```javascript
chapter: number,
section: number,
```

**output**

```javascript
List<{
    type: string,  // 题型
    category: string,  // 大类
    problems: List<{
        book: string,  // 资料名称
        page: number,  // 页码
        lessonName, string,  // 课时名称
        column: string,  // 栏目名称
        idx: number,  // 题目序号
        subIdx: number,  // 小问序号
        isCorrect: bool,  // 正确与否
        time: number,  // 作业布置时间
    }>,
}>

// 注意顺序已经确定，按顺序展示，不需要变动
```
