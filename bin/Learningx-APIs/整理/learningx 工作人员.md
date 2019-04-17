# 目录

[TOC]



# 说明

- URL 均有前缀 /api/v3/staffs

- [ type ] 代表 type 构成的数组类型

- URL 中出现的 {} 代表参数，如 /templates/{templateID}/

- 部分数据结构在Models中进行定义

- 这里的Unix时间戳均精确到秒，即除以1000，如Math.floor(new Date().getTime() / 1000)

- 引用Models中定义的数据结构时，\`inline\`标志含义：在该处展开Model。如定义了Model A：

  ```javascript
  {
      a: int,
  }
  ```

  B：

  ```
  {
      b: int,
      A `inline`,
  }
  ```

  代表B为：

  ```
  {
      b: int,
      a: int,
  }
  ```

  

# Models

### Operation 

```javascript
{
    type: int,  // 1 变更格式 2 添加内容
    font: string, // 字体（仅当type为1变更格式时有效）
    fontSize: float, // 字号（仅当type为1变更格式时有效）
    fontEffect: [string], // 字体效果 bold 加粗， underlined 下划线 italic 倾斜（仅当type为1变更格式时有效）
    alignment: int, // 对齐，0 左对齐 1 居中对齐 2 右对齐（仅当type为1变更格式时有效）
    rowSpacing: float, // 行距（仅当type为1变更格式时有效）， 1  1.5  2
    content: string, // 内容（仅当type为2添加内容时有效）
}
```

### Template

```javascript
{
    name: string, // 模板名称
    info: string,  // 模板说明
    type: string, // 模板类型
    fileName: string,  // 文件名称
    pageType: string,    // 纸张大小，"A3"或者"A4"
    pageDirection: string,  // 纸张方向
    columnCount: int, // 分栏数
    marginTop: float, // 上下页边距
    marginLeft: float, // 左右页边距
    checkProblemList: {
        blankLinesForSelectProFillPro: int,  // 选择题或者填空题空多少行
        blankLinesForEachSubPro: int,  // 大题每小问空多少行
        minBlankLinesOfPro: int,  // 大题空的总行数最少值
        operations: [ Operation ],  // 操作
    },  // CHECK_PROBLEM_LIST变量每项定义
    contentList: {
        blankLinesForSelectProFillPro: int,  // 选择题或者填空题空多少行
        blankLinesForEachSubPro: int,  // 大题每小问空多少行
        minBlankLinesOfPro: int,  // 大题空的总行数最少值
        operations: [ Operation ],  // 操作
    },  // CONTENT_LIST变量每项定义
    operations: [ Operation ], // 文档正文的操作
}
```

### Product

```javascript
{
    problemCode: string,   // 问题代码，错题学习为"E"
    gradation: int,  // 层次， 1 2 3
    depth: int,  // 深度， 1 2 3
    name: string,  // 产品名称
    level: string,  // 产品级别
    object: string,  // 产品对象
    epu: int,  // EPU, 1 2 3
    problemMax: int,  // 题量控制
    wrongProblemStatus: int,  // 错题状态，1 现在仍错，2 曾经错过
    problemType: [ string ],  // 题目种类
    sameTypeMax: int,  // 同类最大题量
    sameTypeSource: [ string ], // 同类来源
    problemTemplateID: string,  // 题目模板
	answerTemplateID: string,  // 答案模板
    borderControl: string,  // 边界控制
    problemSource: [ string ],  // 错题源， 如： ["课本", "平时试卷"]
    serviceType: string,  // 服务类型
    serviceLauncher: string,  // 服务发起
    serviceStartTime: int64,  // 服务开始时间, unix时间戳
    serviceEndTime: int64,  // 服务结束时间, unix时间戳
    serviceTimes: int,  // 服务次数
    serviceDuration: string,  // 服务时长
    deliverType: string,  // 交付类型
    deliverPriority: int,  // 交付优先级
    deliverTime: [{
        day: int,  // 周日0,周一到周六分别是1到6,
        time: string,  // 时间，格式按照"08:00:00"
    }],   // 交付节点
    deliverExpected: int,  // 交付预期，预期多少小时内
	exceptionHandler: int,  // 异常处理，发现未标记：1 全部标记为对再生成， 2 全部标记为错再生成， 3 不生成
    price: int,  // 单价
    subject: string,  // 学科
    grade: string,  // 年级（全部直接用“全部”）
}
```

### Student

```javascript
{
    learnID: int,   // 学生学习号
    name: string,  // 学生名字
    createTime: int64, // 创建时间，Unix时间戳
    gender: string,  // 性别
    grade: string,   // 年级
    class: int,   // 班别
    level: int,  // 层级
}
```

### Book

```javascript
{
    bookID: string,    // 书本识别码
    type: int,   // 书本资料类型，1：课本,2：普通辅导书，3：培优资料
    time: string,  // 录入时间
    name: string,	// 书本名称
    term: string,    // 学期
    version: string,  // 教科书版本
    year: int,    // 教科书年份
    isbn: string,   // ISBN
    ediYear: string,  // 版次年
    ediMonth: string, // 版次月, 两位数字的字符串(所以用string...)
    ediVersion: string,  // 版次第几版
    impYear: string, // 印次年
    impMonth: string,  // 印次月
    impNum: string,  // 印次第几次印刷
    cipURL: string,  // CIP截图URL
    priceURL: string, // 价格截图URL
    coverURL: string,  // 封面URL
}
```

### Paper

```javascript
{
    paperID: string,  // 试卷识别码
    time: string,  // 录入时间
    type: string,   // 试卷类型
    name: string,   // 试卷名称
    fullScore: int,  // 满分	
    version: string,  // 版本
    year: int,    // 教科书版本年份
    choice: string,     // 选择题三个汉字
    blank: string, // 填空题三个汉字
    imageURL: string, // 题头截图URL	
}
```

### BasicProblemForCreatingFiles

```javascript
{
    type: string,  // 题目类型
    book: string,    // 书本名称
    page: int,   // 页码
    column: string,    // 栏目名称
    idx: int,    // 题目在原书中的题目序号
    problemID: string,  // 题目识别码
    subIdxs: [int], // 该题生成文档时需要用到的小问序号，没有小问的题目是 [-1]
    full: bool, // 是否需要完成整一道题的所有小问
    how: string,  // 出题方式
    reason: string,  // 选题依据
}
```

### ProblemForCreatingFiles

```javascript
{
    CheckProblemForCreatingFiles `inline`
    checkProblems: [ CheckProblemForCreatingFiles ],  // 检验题目
}
```



# API

使用 JSON 交换数据。默认当输入有错误时，返回 Status Code 422 (StatusUnprocessableEntity) 和详细错误信息。

**<u>注：下面所有URL均有前缀/api/v3/staffs</u>**



## 工作人员用户系统

### *POST* /archiveData/    封存数据进入下学期

**input**

```javascript
{
    needGradeTransform: bool,  // 是否需要进入下一年级
}
```

**output**

Status Code 200

"successfully archived"

### *POST* /login/    登录

**input**

```javascript
{
    staffID: string,  // 工作人员号码
    password: string,  // 密码
    remember: bool,  // 记住我
}
```

**output**

Status Code 200

"Successfully logged in."

### *POST* /me/logout/    登出

**no input**

**output**

Status Code 200

"Successfully logged out"

### *PUT* /me/password/    修改密码

**input**

```javascript
{
    password: string,  // 新密码
}
```

**output**

Status Code 200

"Successfully changed password."

### *GET* /me/profile/    获取个人信息

**no input**

**output**

```javascript
{
    staffID: string,  // 工作人员号码
    manageClasses: [{
        schoolName: string,	// 学校名称
        schoolID: string,    // 学校识别码
        grade: string,  // 年级 （一、二、三、四...）
        class: int,  // 班级
    }],  // 管理的班级
}
```

### *GET* /me/testLogin/    测试是否已经登陆

**no input**

**output**

Status Code 200:

```javascript
"ok"   // 已经登陆
```

Status Code 401:

没有登陆信息



## 内容库基本信息系统

### *GET* /info/chapsSects/    获取章节信息

**input**

```javascript
semester: string,   // 学期 全部、七上...
```

**output**

```javascript
[{
	chapter: int,
    chapterName: string,
    section: int,
    sectionName: string,
}]
```

### *GET* /info/knowledgePoint/    获取某章节的知识点

**input**

```javascript
chapter: int,   // 章
section: int,   // 节
```

**output**

```javascript
[{
	knowledgeNum: int,  // 知识点序号
    knowledgeName: string,  // 知识点名称
}]
```

### *GET* /info/bookProblems/    获取资料题目

**input**

```javascript
bookID: string,  // 资料识别码
page: int,  // 页码
```

**output**

```javascript
[{
    idx: int,  // 题目序号
    subIdx: int,  // 小问序号
    problemID: string,  // 题目识别码
}]
```

### *GET* /info/paperProblems/    获取试卷题目

**input**

```javascript
paperID: string,  // 试卷识别码
```

**output**

```javascript
[{
    idx: int,  // 题目序号
    subIdx: int,  // 小问序号
    problemID: string,  // 题目识别码
}]
```



## 题目查看系统

### *GET* /info/problemTypes/    获取特定课时的题型信息

**input**

```javascript
chapterStart: int,   // 起始章
sectionStart: int,    // 起始节
lessonStart:int,    // 起始课时
chapterEnd: int,   // 结束章
sectionEnd: int,    // 结束节
lessonEnd: int,    // 结束课时
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级
```

**output**

```javascript
[{
	chapter: int,
    section: int,
    lesson: int,  // 课时序号
    typeName: string,  // 题型名称
    priority: int,  // 学习顺序
    priorityTotal: int,  // 学习顺序总数
    // 展示中的学习顺序： priority/priorityTotal
    category: string,  // 题型大类
    unitExamProb: float, // 单元考试概率
    midtermProb: float, // 期中考试概率
    finalProb: float, // 期末考试概率
    seniorEntranceProb: float, // 中考概率
    examCount: int, // 已考次数
}]
```

### *GET* /info/problems/    获取特定范围的题目信息

**input**

```javascript
typeName: string,
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级
```

**output**

```javascript
[{
	how: string,  // 出题方式
    problemID: string,
    subIdx: int,
    used: int, // 使用情况， 1 已布置， 2 已考  3 已布置已考 4 未使用
    htmlURL: string,  // html 文件URL
    wordURL: string,  // word 文件URL
}]
```

### *POST* /info/getProblemsZip/    获取题目word文件压缩包

**input**

```javascript
{
    typeName: string, // 题型名称
    how: string,  // 出题方式（全部 则为 '全部'）
    problems: [{
        problemID: string,
    }]
}
```

**output**

```javascript
{
    URL: string,
}
```



## 模板系统

### *GET* /templates/   获取特定条件的模板信息

**input**

```javascript
type: string, // 模板类型, all 代表全部
```

**output**

```javascript
[{
    templateID: string, // 模板ID
    date: int64,   // 设计日期， Unix时间戳（精确到秒，即除以1000）
    Template `inline`
}]
```

### *POST* /templates/    新建模板

**input**

```javascript
{
    Template `inline`
}
```

**output**

```javascript
"successfully add a template"
```

### *GET* /templates/{templateID}/    获取某个模板信息

**no input**

**output**

```javascript
{
    templateID: string, // 模板ID
    date: int64,   // 设计日期， Unix时间戳（精确到秒，即除以1000）
    Template `inline`
}
```

### *PUT* /templates/{templateID}/    修改某个模板

**input**

```javascript
{
    Template `inline`
}
```

**output**

```javascript
"Successfully updated this template"
```

### *DELETE* /templates/{templateID}/    删除某个模板

**no input**

**output**

```javascript
"Successfully deleted this template"
```

### *GET* /templates/{templateID}/preview/    获取某个模板的预览PDF

**no input**

**output**

Status Code 200

```javascript
{
    docurl: string,
    pdfurl: string,
} // 目前这里只下载PDF文件就行了
```



## 产品系统

### *GET* /products/    获取特定条件的产品信息

**input**

```javascript
epu1: int,  // epu 类型, -1代表全部epu
object: string, // 产品对象, "all"代表全部
```

**output**

```javascript
[{
    productID: string,  // 产品编号
    date: int64,   // 设计日期， Unix时间戳
    status: bool,  // 服务状态
    Product `inline`
}]
```

### *POST* /products/    新增一个产品

**input**

```javascript
{
    Product `inline`
}
```

**output**

```javascript
"successfully uploaded a product"
```

### *GET* /products/{productID}/    获取某一个产品的信息

**no input**

**output**

```javascript
{
    productID: string,  // 产品编号
    date: int64,   // 设计日期， Unix时间戳
    status: bool,  // 服务状态
    Product `inline`
}
```

### *PUT*  /products/{productID}/    更新某个产品的信息

**input**

```javascript
{
    Product `inline`
}
```

**output**

```javascript
"Successfully updated this product"
```

### *PUT*  /products/{productID}/status/    更新某个产品的状态

**input**

```javascript
{
    status: bool,   // 服务状态,true运行 false停止
}
```

**output**

```javascript
"Successfully updated the status of this product"
```



## 学校系统

### *GET* /schools/    获取特定条件的学校

**input**

```javascript
province: string,  // 省，可以为""，""代表该参数不做限制，以下三个参数也是
city: string,  // 市
district: string,  // 区
county: string,   // 县
```

**output**

```javascript
[{
	name: string,	// 学校名称
	schoolID: string,    // 学校识别码
}]
```

### *POST* /schools/    新增一个学校

**input**

```javascript
{
    province: string,  // 省
    city: string,   // 市
    district: string,    // 区
    county: string,   // 县
    name: string,  // 学校
}
```

**output**

```javascript
{
    schoolID: string,   // 新增的学校的schoolID
}
```



## 班级系统

### *POST* /studentFile/    上传学生名单并预览

**input**

```javascript
file: file,  // 学生名单	采用multipart形式上传
```

**output**

```javascript
{
    "uid": "63362a46-0894-432a-8fc6-8b5594eccdfe", // UID 用于标识这一文件

    "columns": [{
        "title": "姓名",
        "dataIndex": "name"
    }, {
        "title": "性别",
        "dataIndex": "gender"
    }], // 预览表中的列

    "data": [{
        "name": string, // 姓名
        "gender": string, // 性别
    }], // 预览表中的数据
}
```

### *DELETE* /studentFile/{UID}/    删除UID对应的学生信息临时文件

**no inputs**

**output**

```
"Successfully deleted"
```

### *POST* /students/    上传班级学生信息

 如果对应班级不存在则新增. 返回学生账户列表下载URL

**input**

```javascript
{
    schoolID: string, // 学校识别码
    grade: string,   // 年级
    class: int,    // 班级
    studentFile: string,  // 学生名单的UID
}
```

**output**

```javascript
{
    URL: string,
}
```

### *POST* /students/addOne/    添加一个新学生信息

**input**

```javascript
{
    schoolID: string, // 学校识别码
    grade: string,   // 年级
    class: int,    // 班级
    name: string,  // 学生名字
    gender: string,  // 性别
}
```

**output**

```javascript
"successfully added"
```

### *GET* /classes/students/    获取某个班级特定条件的学生信息

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,  // 年级 （一、二、三、四...）
class: int,  // 班级，0代表全部
studentName: string,  // 学生姓名筛选（筛选姓名含有该字符串的学生），""代表不做限制
epu: string,  // epu，""代表不做限制
serviceType: string,  // 服务类型，""代表不做限制
productID: string,  // 产品ID，""代表不做限制
    
// productID 优先级高于 epu、serviceType，即设定了查询拥有某 productID 的学生，则 epu、serviceType 无效，以避免 productID 与 epuStr、serviceType 冲突问题
```

**output**

```javascript
{
	total: int,    // 学生总人数
	learnIDs: [ Student ],  // 学生列表
}
```

### *GET* /classes/books/    获取某个班级的书本信息

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级, 0代表全部
```

**output**

```javascript
[ Book ]
```

### *GET* /classes/papers/    获取某个班级的试卷信息

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级
```

**output**

```javascript
[ Paper ]
```

### *GET* /books/search/    搜索一本资料

**input**

```javascript
isbn: string,   // ISBN
ediYear: string,  // 版次年
ediMonth: string, // 版次月, 两位数字的字符串(所以用string...)
ediVersion: string,  // 版次第几版
impYear: string, // 印次年
impMonth: string,  // 印次月
impNum: string,  // 印次第几次印刷
```

**output**

```javascript
[{
    bookID: string, // 资料识别码
    coverURL: string, // 封面图片URL
    cipURL: string,
    priceURL: string,  // 印版次数据即价格图片...
}]
```

### *POST* /classes/addBooks/    给班级添加书本

**input**

```javascript
{
    schoolID: string, // 学校识别码
    grade: string,   // 年级
    class: int,    // 班级, 0代表全部班级
    bookID: string, // 资料识别码
}
```

**output**

```javascript
"Successfully add a book"
```

### *GET* /papers/search/    搜索一张试卷

**input**

```javascript
choice: string,     // 选择题三个汉字
blank: string, // 填空题三个汉字
calculation: string, // 压轴题三个汉字
```

**output**

```javascript
[{
    paperID: string, // 试卷识别码
    name: string, // 试卷名称
    fullScore: int,  // 满分
    image: string,  // 题头图片URL
}]
```

### *POST* /classes/addPapers/    给班级添加试卷

**input**

```javascript
{
    schoolID: string, // 学校识别码
    grade: string,   // 年级
    class: int,    // 班级, 0代表全部班级
    paperID: string, // 试卷识别码
}
```

**output**

```javascript
"Successfully add a paper"
```

### *POST* /classes/deleteBooks/  删除某个班级与某本书的对应

**input**

```javascript
{
    schoolID: string, // 学校识别码
	grade: string,   // 年级
	class: int,    // 班级, 0代表全部班级
	bookID: string,    // 书本识别码
}
```

**output**

```javascript
"successfully deleted!"
```

### *POST* /classes/deletePapers/    删除某个班级与某个试卷的对应

**input**

```javascript
{
	schoolID: string, // 学校识别码
	grade: string,   // 年级
	class: int,    // 班级, 0代表全部班级
	paperID: string,   // 试卷识别码
}
```

**output**

```javascript
"successfully deleted!"
```



## 学生管理系统

### *GET* /students/{learnID}/    获取某个学生的个人信息

**no input**

**output**

```javascript
{
    name: string,  // 姓名
    gender: string,  // 性别
    school: string,  // 学校名称
    grade: string,   // 年级
    class: int,   // 班别
    productID: string,  // 在用产品信息, ""空字符串代表没有在用产品
}
```

### *PUT* /students/{learnID}/productID/     修改某个学生的个人产品信息

**input**

```javascript
{
    productID: string,  // 产品信息
}
```

**output**

```javascript
"Successfully updated productID of this student"
```

### *PUT* /students/{learnID}/     修改某个学生的个人信息

**input**

```javascript
{
    name: string,  // 学生名字
    gender: string,  // 性别
    class: int,   // 班别
}
```

**output**

```javascript
"successfully updated this student"
```

### *DELETE* /students/{learnID}/    删除某个学生

**no input**

**output**

```javascript
"successfully deleted!"
```



## 文档生成系统

### *POST* /batchDownloads/    生成一个批量生成文档的任务ID

后续根据该ID追踪该批量生成任务的完成情况

**input**

```javascript
{
    school: string, // 学校名称
	grade: string,   // 年级
	class: int,    // 班级
    students: [{
        name: string,
        learnID: int,
    }],  // 需要生成纠错本的学生的信息
}
```

**output**

```javascript
{
    batchID: string,
}
```

### *DELETE* /batchDownloads/{batchID}/     删除批量生成文档的任务ID

**no input**

**output**

```javascript
"successfully deleted batchID"
```

### *GET* /batchDownloads/    获取批量生成文档的情况

**no input**

**output**

```javascript
[{
    batchID: string,  // 批量生成文档下载任务ID
    createTime： int64, // 该任务创建时间，Unix时间戳
	school: string,      // 学校名称
	grade: string,       // 年级
	class: int,          // 班级
    finishTime: int64,   // 完成时间或预计完成时间，unix时间戳
    problemFilesFinished: int,  // 已处理的题目文件数目（无论是否成功）
    answerFilesFinished: int,  // 已处理的答案文件数目（无论是否成功）
    problemFilesSuccessful: int,  // 已完成并成功生成的题目文件数目
    answerFilesSuccessful: int,   // 已完成并成功生成的答案文件数目
	students: [{
		name: string,
		learnID: int,                   
		problemFileStatus: bool,      // 题目文件生成状态(是否成功)
		answerFileStatus: bool,       // 答案文件生成状态(是否成功)
		problemStatusCode: int,    // 题目文件生成请求得到的状态码
		answerStatusCode: int,     // 答案文件生成请求得到的状态码
		problems: [ ProblemForCreatingFiles ],  // 错题信息, 与epu生成文档流程中获取错题信息得到的数据一样
    }],  // 需要生成纠错本的学生的信息
}]
```

说明：

总人数：students数组长度

文档状态有：初始：还没处理完，即开始处理了或者没开始处理（statusCode: 0, status: false）。处理完并成功生成了（statusCode: 200, status: true），处理完并出错了（statusCode: 4xx or 5xx, status false）

错误状态：当生成状态是true时，statusCode字段为200：正常生成文档。若生成状态是false，statusCode字段判断错误状态，400：没有找到错题（未标记或者不存在错题）；504：处理超时； 500：内部未知错误； 404：题目或者答案文档缺失；403：存在未标记的纠错本，不允许生成新文档；0：初始状态，此时还没处理完成。

### *POST* /students/{learnID}/getWrongProblems/    根据产品获取某个学生的错题信息

**input**

```javascript
{
    productID: string, // 产品ID
    wrongProblemType: int,   // 错题类型，1 现在仍错的错题，2 曾经错过的错题
    sort: int,   // 排序方式，1按出题方式，2按题目类型
    bookPage: [{
        bookID: string,  // 学习资料识别码
        startPage: int,  // 开始页码
        endPage: int,  // 结束页码
    }],  // 书本，仅当产品EPU为1时有效
    paperIDs: [string],    // 字符串列表，试卷识别码，仅当产品EPU为1时有效
    batchID: string,   // 批量生成文档的任务ID
}
```

**output**

Status Code 200:

```javascript
[ ProblemForCreatingFiles ]
```

Status Code 404:

"No wrong problems"

### *POST* /students/{learnID}/documents/    生成一份文档

**input**

```javascript
{
    productID: string, // 产品ID
    docType: int,  // 文档类型 1 生成题目文件 2 生成答案文件
    batchID: string,   // 批量生成文档的任务ID
    contents: [ ProblemForCreatingFiles ],  // 文档内容数据
}
```

**output**

Status Code 200:

```javascript
"file is generating"
```

### *POST* /students/getDocumentZip/    文件打包，获取压缩包下载URL

**input**

```javascript
{
    grade: string,	// 年级
    class: int, // 班别
    batchID: string,  // 批量下载ID
}
```

**output**

```javascript
{
	URL: string,
}
```



## 学生标记任务系统

### *POST* /students/markTasks/    批量新增标记任务

**input**

```javascript
[{
    time: int64,   // 任务对应的unix时间戳（即当前的时间戳）
    type: int,   // 任务类型（没标记的是错题为1，没标记的是检验题为2）
    learnID: int,    // 学习号
    problems: [ ProblemForCreatingFiles ],  // 待标记的错题
}]
```

**output**

Status Code 200:

```
"Successfully created tasks."
```

Status Code 404:

```
[int, ]  // learnID构成的列表
```

当有部分用户找不到或者添加task失败的时候，状态码404并同时返回失败的用户的learnId构成的列表。没在列表中的用户则是成功创建了的。

### *GET* /classes/markTasks/    获取班级所有学生没有完成的标记任务

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级， 0 代表全部
```

**output**

Status Code 200:

```javascript
[{
    learnID: int, // 学生学习号
    name: string, // 学生名字
	time: int64,  // 任务发生的时间对应的unix时间戳
	type: int,  // 任务类型（没标记的是错题为1，没标记的是检验题为2）
}]
```

Status Code 404:

```
"No upload tasks."
```

没有未完成的标记任务

### *DELETE* /classes/markTasks/    批量删除标记任务

**input**

```javascript
[{
    learnID: int, // 学生学习号
	time: int64,  // 任务发生的时间对应的unix时间戳
}]
```

**output**

Status Code 200:

```
"Successfully deleted"
```

### *GET* /students/{learnID}/markTasks/    获取某个学生没有完成的标记任务

**input**

no input

**output**

Status Code 200:

```javascript
[{
	time: int64,  // 任务发生的时间对应的unix时间戳
	type: int,  // 任务类型（没标记的是错题为1，没标记的是检验题为2）
}]
```

Status Code 404:

```
"No upload tasks."
```

没有未完成的上传任务

### *GET* /students/{learnID}/markTasks/{time}/    获取标记任务内容

**input**

no input

**output**

Status Code 200:

```javascript
[ ProblemForCreatingFiles ]
```

Status Code 404:

```
"Can not find the task."
```

找不到这个时间戳对应的任务

### *DELETE* /students/{learnID}/markTasks/{time}/    删除学生某一个标记任务

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



## 学生标记情况系统

### *GET* /students/problemRecords/    查看标记情况记录信息

附：

页面中如果纠错本未标记（wrongProblemStatus == 1）或者试卷未标记（paperStatus == 1）显示未标记（不管bookStatus）

每个学生的标记评估计算：10 × （ status 是 0 的个数 / status 的总数），然后四舍五入取整。（这里的status，包括wrongProblemStatus，paperStatus 以及 bookStatus 里面的 status）

**input**

```javascript  
[ int ]    // learnID 构成的列表
```

**output**

```javascript
[{
    learnID: int,  // 学习号
    wrongProblemStatus: int,		// 纠错本状态，1未标记，0已标记
    paperStatus: int,		// 试卷状态，1未标记，0已标记
    bookStatus: [{
        book: string,    // 资料名称
        status: int,    // 状态，0最近一周有标记，1没有
    }],  // 书本状态
}]

// 即统一为0无异常，1有异常
```



## 学生错题信息录入

### *GET* /students/{learnID}/notMarkedPapers/    获取某个学生的未标记试卷信息

**no input**

**output**

```javascript
[{
    name: string,   // 试卷名称
    paperID: string,  // 试卷识别码
}]
```

### *GET* /students/{learnID}/paperProblems/    获取试卷中还没录入结结果的题目

**input**

```javascript
paperID: string,  // 试卷识别码
```

**output**

```javascript
[{
    lessonName: string,  // 课时名称（恒为""）
    column: string,  // 栏目名称（恒为""）
    idx: int,  // 题目序号
    subIdx: int,  // 小问序号
    problemID: string,  // 题目识别码
}]
```

### *GET* /students/{learnID}/bookProblems/    获取资料中还没录入结结果的题目

**input**

```javascript
book: string,  // 资料识别码
page: int,  // 页码
```

**output**

```javascript
[{
    lessonName, string,  // 课时名称
    column: string,  // 栏目名称
    idx: int,  // 题目序号
    subIdx: int,  // 小问序号
    problemID: string,  // 题目识别码
}]
```

### *POST* /students/{learnID}/problems/    提交错题录入信息

**input**

```javascript
{
    time: int64,  // 作业布置时间，unix时间戳  使用当前时间即可
	problems: [{
	    isCorrect: bool,
	    problemID: string,
	    subIdx: int,
        sourceID: string,  // 题目来源（对应的bookID或者paperID，纠错本的话为""）
       	sourceType: int, // 题目来源类型（书本 1 试卷 2 纠错本 3）
	}]
}
```

**output**

```javascript
"successfully uploaded"
```



## 班级学期起始标记

### *GET* /classes/semester/    获取班级学期起始标记结果

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级, 0 代表全部
```

**output**

```javascript
{
    semester: string,  // 学期，值为 "上" "下" "未定"
    startTime: int64,  // 学期开始时间，unix时间戳
    endTime: int64,  // 学期结束时间，unix时间戳
}
```

### *POST* /classes/semester/    提交班级学期起始标记

**input**

```javascript
{
    schoolID: string, // 学校识别码
    grade: string,   // 年级
    class: int,    // 班级, 0 代表全部
    semester: string,  // 学期，值为 "上" "下" "未定"
    startTime: int64,  // 学期开始时间，unix时间戳
    endTime: int64,  // 学期结束时间，unix时间戳
}
```

**output**

```javascript
"successfully uploaded"
```



## 班级作业布置情况系统

### *GET* /classes/bookProblems/    获取资料中还没布置的题目

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级
bookID: string,  // 资料识别码
page: int,  // 页码
```

**output**

```javascript
[{
    idx: int,  // 题目序号
    subIdx: int,  // 小问序号
    problemID: string,  // 题目识别码
}]
```

### *POST* /classes/assignments/    提交已经布置的作业

只提交勾选了布置了的题目

**input**

```javascript
{
    schoolID: string, // 学校识别码
    grade: string,   // 年级
    class: int,    // 班级, 0 代表全部
    time: int64,  // 作业布置时间，unix时间戳
	problems: [{
	    problemID: string,
	    subIdx: int,
	}]
}
```

**output**

```javascript
"successfully uploaded"
```



## 班级考试成绩录入

### *GET* /classes/papersForMarkScore/    获取某个班级的用于标记考试成绩的试卷信息

与直接获取试卷信息不同，这里标记完成的试卷排序在后面

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级
```

**output**

```javascript
[{
    paperID: string,  // 试卷识别码
    name: string,   // 试卷名称
    fullScore: int,  // 试卷满分
    marked: bool, // 是否已经标记了成绩
}]
```

### *POST* /classes/scoreFile/    上传班级成绩excel表获取数据

表格格式：第一行："姓名" "成绩"，后面每行分别是每个学生的姓名和成绩

**input**

```javascript
file: file,  // 学生名单	采用multipart形式上传
```

**output**

Status Code 200:

```javascript
[{
    name: string,  // 姓名（暂不考虑班级内出现重名的情况）
    score: float,   // 成绩
}]
```

Status Code 403:

```javascript
"wrong format"  // 表格格式不正确
```

### *POST* /classes/examScores/    录入班级所有学生考试成绩

**input**

```javascript
{
    schoolID: string, // 学校识别码
    grade: string,   // 年级
    class: int,    // 班级, 0 代表全部
    time: int64,  // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
    paperID: string,  // 考试试卷ID
	scores: [{
	    learnID: int,  // 学生学习号
        name: string, // 学生姓名
	    score: double,  // 成绩
	}],  // 成绩
}
```

**output**

```javascript
"successfully uploaded"
```

### *GET* /classes/examScores/    获取班级某次考试成绩

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级, 0 代表全部
paperID: string,  // 考试试卷ID
```

**output**

```javascript
{
    time: int64,  // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
	scores: [{
	    learnID: int,  // 学生学习号
        name: string, // 学生姓名
	    score: double,  // 成绩
	}],  // 成绩
}
```



## 班级知识讲解标记

### *GET* /classes/knowledgeLearned/nextChapSect/    获取班级知识讲解标记预测章节

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级, 0 代表全部
```

**output**

```javascript
{
    chapter: int,   // 章
    section: int,    // 节
}
```

### *POST* /classes/knowledgeLearned/    提交班级知识讲解标记结果

**input**

```javascript
{
    schoolID: string, // 学校识别码
    grade: string,   // 年级
    class: int,    // 班级, 0 代表全部
   	time: int64,  // 讲课时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
    knowledges: [{
        chapter: int,   // 章
        section: int,    // 节
        knowledgeNum: [ int ],  // 知识点序号构成的数组
    }],
}
```

**output**

```javascript
"successfully uploaded"
```



## 班级题目讲解标记

### *POST* /classes/problemsLearned/methodOne/    提交班级题目讲解标记结果(方式1)

**input**

```javascript
{
    schoolID: string, // 学校识别码
    grade: string,   // 年级
    class: int,    // 班级, 0 代表全部
   	time: int64,  // 讲课时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
    problems: [{
        problemHow: string,  // 出题方式（选择题、填空题、解答题）
        source: string,  // 题目来源
    }],
}
```

**output**

```javascript
"successfully uploaded"
```

### *POST* /classes/problemsLearned/methodTwo/    提交班级题目讲解标记结果(方式2)

**input**

```javascript
{
    schoolID: string, // 学校识别码
    grade: string,   // 年级
    class: int,    // 班级, 0 代表全部
   	time: int64,  // 讲课时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
    problems: [{
        problemID: string,  // 题目识别码
        subIdx: int,  // 小问序号
    }],   // 这里只传勾选了讲解的题目
}
```

**output**

```javascript
"successfully uploaded"
```



## 班级分层系统

### *GET* /classes/totalLevel/    获取班级的总层数

**input**

```javascript
schoolID: string,  // 学校识别码
grade: string, // 年级 （一、二、三、四...）
class: int,  // 班级
```

**output**

```javascript
{
    totalLevel: int,
}
```

status Code 404:

还没进行设定

### *POST* /classes/totalLevel/    修改当前班级的总层数

**input**

```javascript
{
    schoolID: string, // 学校识别码
    grade: string, // 年级 （一、二、三、四...）
    class: int, // 班级
    totalLevel: int,
}
```

**output**

```javascript
"successfully update total level of this class"
```

### *PUT* /classes/students/level/    更新学生层级信息

**input**

```javascript
[{
    learnID: int,  // 学生学习号
    level: int,    // 层级
}]
```

**output**

```javascript
"Successfully updated students' levels."
```



## 班级产品配置系统

### *GET* /classes/productID/    获取某个班级的产品信息

**input**

```javascript
schoolID: string,  // 学校识别码
grade: string,   // 年级
class: int,   // 班别
level: int,  // 层级，-1获取这个班级的整体产品信息
```

**output**

```javascript
{
    productID: [string],  // 在用产品信息, []空列表代表没有在用产品
}
```

### *PUT* /classes/productID/    修改某个班级的产品信息

**input**

```javascript
{
    schoolID: string,  // 学校识别码
    grade: string,   // 年级
    class: int,   // 班别
    level: int， // 层级，-1代表修改整个班级整体产品信息
    productID: [string],  // 产品信息
}
```

**output**

```javascript
"Successfully updated productID of this class"
```



## 目标规划系统

### *GET* /classes/targets/    获取某个班级的目标规划信息

**input**

```javascript
schoolID: string,  // 学校识别码
grade: string,   // 年级
class: int,   // 班别
level: int,  // 层级
exam: string,  // 考试
chapter: int,  // 章 0代表全部
section: int,  // 节 0代表全部
typename: string,  // 搜索题型名称
semester: string,   // 学期 全部、七上...
```

**output**

```javascript
[{
    status: bool,  // 是否已加入
    chapter: int, // 章
    section: int,  // 节
    typename: string,  // 题型名称
    originalKP： string,  // 最新原始知识点
}]
```

### *POST* /classes/targets/    将某些目标题型加入到班级该层中

**input**

```javascript
{
    schoolID: string,  // 学校识别码
    grade: string,   // 年级
    class: int,   // 班别
    level: int,  // 层级
    exam: string,  // 考试
    targets: [{
        chapter: int, // 章
        section: int,  // 节
        typename: string,  // 题型名称
    }],
}
```

**output**

```javascript
"successfully added"
```

### *DELETE* /classes/targets/    将某些目标题型从班级该层中移除

**input**

```javascript
{
    schoolID: string,  // 学校识别码
    grade: string,   // 年级
    class: int,   // 班别
    level: int,  // 层级
    exam: string,  // 考试
    targets: [{
        chapter: int, // 章
        section: int,  // 节
        typename: string,  // 题型名称
    }],
}
```

**output**

```javascript
"successfully deleted"
```



## 班级错误率分析系统

### *POST* /classes/getErrorRateAnalysis/    获取班级错误率分析结果

**input**

```javascript
{
    wrongProblemStatus: int,   // 错题状态，1现在仍错的题目，2曾经错过的
    bookPage: [{
        bookID: string,  // 学习资料识别码
        startPage: int,  // 开始页码
        endPage: int,  // 结束页码
    }],  // 书本
    paperIDs: [string],    // 字符串列表，试卷识别码
    schoolID: string, // 学校识别码
    grade: string,   // 年级
    class: int,    // 班级
    level： int,  // 层级，-1代表全部
    exam: string, // 考试
    dateBefore: int64, // Unix时间戳，分析什么日期之前的错题，除以1000精确到秒
}
```

**output**

Status Code 200:

```javascript
[{
    source: string,  // 来源， 书本名称或者试卷名称
    page: int,
    column: string,    // 栏目名称
    idx: int,    // 题目序号
    problemID: string,  // 题目识别码
    subIdx: int,  // 小问序号
    probability: float, // 考试概率
    errorRate: float,  // 错误率
    wrongStudents: [string],  // 错误的学生名单
    totalStudents: int, // 分析的学生总数（因为选择了分析某个层级的学生，所以这里可能不等于班级学生总数）
}]
```

### *GET* /classes/practiceProblems/    根据预选择的题目获取真正的训练题目

**input**

```javascript
{
    bookPage: [{
        bookID: string,  // 学习资料识别码
        startPage: int,  // 开始页码
        endPage: int,  // 结束页码
    }],  // 书本
    paperIDs: [string],    // 字符串列表，试卷识别码
    // 前两个字段是错误率分析界面一开始选择的数据
    problems: [{
        problemID: string,  // 题目识别码
        subIdx: int,  // 小问序号
    }],  // 用户选择训练的题目
}
```

错误率分析获取题目与文件暂时不改



## 班级考试分析系统

### *GET* /classes/examAnalysis/average/    获取某个班级的均分分析

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级, 0 代表全部
startTime: int64,  // 分析开始时间unix时间戳
endTime: int64,  // 分析结束时间unix时间戳
standardFullScore: int,   // 标定的满分标准，所有考试的分数都会转化成满分分数为standardFullScore的形式
```

**output**

```javascript
{
    latestTop10: [{
        learnID: int,  // 学习号
    	name: string,  // 学生名字
        score：float,  // 分数
    }],  // 最新考试前10名（右侧展示）
    latestLast10: [{
        learnID: int,  // 学习号
    	name: string,  // 学生名字
        score：float,  // 分数
    }],  // 最新考试后10名（右侧展示）
    exams: [{
    	time: int64,  // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
    	paperID: string,  // 考试试卷ID
    	averageScore: float,  // 均分
	}],  // 已经按时间顺序排序
}
```

### *GET* /classes/examAnalysis/rankingLevelAverage/    获取某个班级的排名段均分分析

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级, 0 代表全部
startTime: int64,  // 分析开始时间unix时间戳
endTime: int64,  // 分析结束时间unix时间戳
standardFullScore: int,   // 标定的满分标准，所有考试的分数都会转化成满分分数为standardFullScore的形式
```

**output**

```javascript
[{
    latestRankinglevel: string,  // 最新考试的排名段
    data: [{
        time: int64,  // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
    	paperID: string,  // 考试试卷ID
    	averageScore: float,  //  该分数层这次考试均分
	}],   // 考试已经按时间顺序排序
}]
```

### *GET* /classes/examAnalysis/scoreProportion/    获取某个班级的分数段占比分析

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级, 0 代表全部
startTime: int64,  // 分析开始时间unix时间戳
endTime: int64,  // 分析结束时间unix时间戳
standardFullScore: int,   // 标定的满分标准，所有考试的分数都会转化成满分分数为standardFullScore的形式
```

**output**

```javascript
[{
    scoreSegment: int,  // 分数段
    data: [{
        time: int64,  // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
    	paperID: string,  // 考试试卷ID
    	rate: float,  // 该分数层这次分数段占比
	}],   // 考试已经按时间顺序排序
}]
```

### *GET* /classes/examAnalysis/studentScore/    获取某个班级的个人分数分析

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级, 0 代表全部
startTime: int64,  // 分析开始时间unix时间戳
endTime: int64,  // 分析结束时间unix时间戳
standardFullScore: int,   // 标定的满分标准，所有考试的分数都会转化成满分分数为standardFullScore的形式
```

**output**

Status Code 200:

```javascript
[{
    level: int,  // 层级
    students: [{
        learnID: int,  // 学习号
        name: string,  // 学生名字
        latestRanking: int, // 最新排名（右侧表格展示）
        data: [{
            time: int64,  // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
            paperID: string,  // 考试试卷ID
            score: float,  // 这次考试分数
        }],   // 考试已经按时间顺序排序
    }],
}]
```

Status Code 403:

存在尚未分层的学生，请先完成分层

```javascript
"some students do not have levels"
```

### *GET* /classes/examAnalysis/studentRanking/    获取某个班级的个人排名分析

**input**

```javascript
schoolID: string, // 学校识别码
grade: string,   // 年级
class: int,    // 班级, 0 代表全部
startTime: int64,  // 分析开始时间unix时间戳
endTime: int64,  // 分析结束时间unix时间戳
standardFullScore: int,   // 标定的满分标准，所有考试的分数都会转化成满分分数为standardFullScore的形式
```

**output**

Status Code 200:

```javascript
[{
    level: int,  // 层级
    students: [{
        learnID: int,  // 学习号
        name: string,  // 学生名字
        latestRanking: int, // 最新排名（右侧表格展示）
        data: [{
            time: int64,  // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
            paperID: string,  // 考试试卷ID
            ranking: int,  // 这次考试排名
        }],   // 考试已经按时间顺序排序
    }],
}]
```

Status Code 403:

存在尚未分层的学生，请先完成分层

```javascript
"some students do not have levels"
```

### *POST* /classes/examAnalysis/thoughts/    提交教学思考    

**input**

```javascript
{
    schoolID: string,  // 学校识别码
    grade: string,   // 年级
    class: int,   // 班别
    paperID: string,  // 试卷识别码，注：提交的是最新的那次考试的识别码，即，包含了paperID那个数组中的最后一个元素中的paperID
    imageType: int,  // 1 班级均分分析图，2 班级排名段分析图 3 班级分数段分析图  4 个人分数分析图  5 个人排名分析图
    scoreSegment: int,  // 分数段(当为班级分数段分析图才有效)
    level: int,  // 层级(个人分数分析或者个人排名分析才有效)
    thoughts: string,  // 教学思考
}
```

**output**

```javascript
"successfully added thoughts"
```

### 