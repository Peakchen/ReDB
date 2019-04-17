Learningx Admin Server
----------------------

# 目录

* [创建配置文件](#%E5%88%9B%E5%BB%BA%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)
* [创建 MySQL Server](#mysql)
* [构建应用](#%E6%9E%84%E5%BB%BA%E5%BA%94%E7%94%A8)
* [API](#api)
    * [/api/admin-v1/login](#apiadmin-v1login)
    * [/api/admin-v1/logout](#apiadmin-v1logout)
    * [/api/admin-v1/users](#apiadmin-v1users)
    * [/api/admin-v1/blocks](#apiadmin-v1blocks)
    * [/api/admin-v1/chapters](#apiadmin-v1chapters)
    * [/api/admin-v1/sections](#apiadmin-v1sections)
    * [/api/admin-v1/typenames](#apiadmin-v1typenames)
    * [/api/admin-v1/probtypes](#apiadmin-v1probtypes)
    * [/api/admin-v1/probmetas](#apiadmin-v1probmetas)
    * [/api/admin-v1/hows](#apiadmin-v1hows)
    * [/api/admin-v1/images](#apiadmin-v1images)
    * [/api/admin-v1/contents](#apiadmin-v1contents)
    * [/api/admin-v1/extremums](#apiadmin-v1extremums)

# 创建配置文件

```bash
cp config.example.yml ./conf/config.yml

## 修改配置
vim ./conf/config.yml
```

# MySQL

```bash
docker pull mariadb:10.1.26

# 参考文档 https://hub.docker.com/r/library/mariadb/
# 根据实际情况修改

docker run -d \
    -e MYSQL_ROOT_PASSWORD=secret-root-passwd \
    -e MYSQL_DATABASE=learningx \
    -e MYSQL_USER=learningx \
    -e MYSQL_PASSWORD=learningx \
    -p 127.0.0.1:3306:3306 \
    -v $DATA:/var/lib/mysql \
    --restart unless-stopped \
    --name db mariadb:10.1.26 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

# 构建应用

```bash
docker pull golang:1.9

docker run --rm -i \
    -v $PWD:/go/src/gitee.com/learningx2017/learningx_admin_server \
    -w /go/src/gitee.com/learningx2017/learningx_admin_server \
    golang:1.9 ./build.sh
```

# API

## /api/admin-v1/login

### POST

* 描述: 用户登录

请求体:

```
{
    "userName": "admin",
    "password": "123456"
}
```

参数:

* remember
    * 描述: 是否记住登录状态
    * 可选值: "0", "1"

返回值:

```
{
    "uid": "0", // 用户 ID
    "userName": "admin", // 用户名
    "roles": ["p1", "p2"], // 允许访问的页面别称
}
```

响应会把 JWT 保存在名为 `TOKEN` 的 Cookie 中.

允许访问的页面:
{
    "p1": "题目数据 Excel 表",
    "p2": "基本 Excel 表",
    "p3": "信息序号最值表",
    "p4": "例题文档",
    "p5": "习题和答案图片",
}

## /api/admin-v1/logout

### POST

* 描述: 用户退出登录
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

响应会把名为 Cookie 中的 `TOKEN` 的值清空.

## /api/admin-v1/me

### GET

* 描述: 返回用户信息
* 要求: 发送的请求需带有 JWT

返回值:

```
{
    "uid": "0", // 用户 ID
    "userName": "admin", // 用户名
    "roles": ["p1", "p2"], // 允许访问的页面别称
}
```

## /api/admin-v1/users

### GET

* 描述: 返回所有用户信息
* 要求: 发送的请求需带有 JWT 并且用户是 `admin`

返回值:

```
[
{
    "uid": "0", // 用户 ID
    "userName": "admin", // 用户名
    "roles": ["p1", "p2"], // 允许访问的页面别称
},
{
    "uid": "1", // 用户 ID
    "userName": "staff1", // 用户名
    "roles": ["p1", "p2", "p3"], // 允许访问的页面别称
},
...
]
```

### PUT

* 描述: 更新某用户允许访问的页面
* 要求: 发送的请求需带有 JWT 并且用户是 `admin`

请求体:

```
{
    "uid": "1", // 目标用户 ID
    "roles": ["p1", "p2", "p3"], // 允许访问的页面别称
}
```

返回值:

```
None
```

## /api/admin-v1/blocks

### POST

* 描述: 上传知识块信息表 (并没有真正保存)
* 要求: 发送的请求需带有 JWT

请求体:

```
File: 知识块信息表.xlsx
```

返回值:

```
{
    "uid": "63362a46-0894-432a-8fc6-8b5594eccdfe", // UUID 用于标识这一操作

    "columns": [{
        "title": "教材识别码",
        "key": "教材识别码",
        "dataIndex": "bookID"
    }, {
        "title": "章序号",
        "key": "章序号",
        "dataIndex": "chapNum"
    },
    ....
    ], // 预览表中的列

    "data": [{
        "chapNum": 1, // 章序号
        "sectNum": 1, // 节序号
        "num": 1, // 知识块序号
        "name": "样例", // 知识块名称

        "key": "uniquekey", // 唯一标识
    },
    ...
    ], // 预览表中的数据
}
```

## /api/admin-v1/blocks/:uid

### POST

* 描述: 确认保存 `:uid` 对应的知识块信息表
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

### DELETE

* 描述: 不保存 `:uid` 对应的知识块信息表并删除临时存储的数据
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

## /api/admin-v1/chapters

### POST

* 描述: 上传章名称信息表 (并没有真正保存)
* 要求: 发送的请求需带有 JWT

请求体:

```
File: 知识块信息表.xlsx
```

返回值:

```
{
    "uid": "63362a46-0894-432a-8fc6-8b5594eccdfe", // UUID 用于标识这一操作

    "columns": [{
        "title": "教材识别码",
        "key": "教材识别码",
        "dataIndex": "bookID"
    }, {
        "title": "章序号",
        "key": "章序号",
        "dataIndex": "num"
    },
    ....
    ], // 预览表中的列

    "data": [{
        "num": 1, // 知识块序号
        "name": "样例", // 知识块名称

        "key": "uniquekey", // 唯一标识
    },
    ...
    ], // 预览表中的数据
}
```

## /api/admin-v1/chapters/:uid

### POST

* 描述: 确认保存 `:uid` 对应的表
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

### DELETE

* 描述: 不保存 `:uid` 对应的表并删除临时存储的数据
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

## /api/admin-v1/sections

### POST

* 描述: 上传节名称信息表 (并没有真正保存)
* 要求: 发送的请求需带有 JWT

请求体:

```
File: 节名称信息表.xlsx
```

返回值:

```
{
    "uid": "63362a46-0894-432a-8fc6-8b5594eccdfe", // UUID 用于标识这一操作

    "columns": [{
        "title": "教材识别码",
        "key": "教材识别码",
        "dataIndex": "bookID"
    }, {
        "title": "章序号",
        "key": "章序号",
        "dataIndex": "chapNum"
    },
    ....
    ], // 预览表中的列

    "data": [{
        "chapNum": 1, // 章序号
        "num": 1, // 节序号
        "name": "样例", // 节名称

        "key": "uniquekey", // 唯一标识
    },
    ...
    ], // 预览表中的数据
}
```

## /api/admin-v1/sections/:uid

### POST

* 描述: 确认保存 `:uid` 对应的表
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

### DELETE

* 描述: 不保存 `:uid` 对应的表并删除临时存储的数据
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

## /api/admin-v1/typenames

### POST

* 描述: 上传题型汇总表 (并没有真正保存)
* 要求: 发送的请求需带有 JWT

请求体:

```
File: 题型汇总表.xlsx
```

返回值:

```
{
    "uid": "63362a46-0894-432a-8fc6-8b5594eccdfe", // UUID 用于标识这一操作

    "columns": [{
        "title": "节",
        "key": "节",
        "dataIndex": "sectNum"
    }, {
        "title": "章序号",
        "key": "章序号",
        "dataIndex": "chapNum"
    },
    ....
    ], // 预览表中的列

    "data": [{
        "chapNum": 1, // 章序号
        "sectNum": 1, // 节序号
        "name": "样例", // 题型名称
        "priority": 1, // 题型学习顺序
        "category": "基础类", // 题型大类
        "source": "辅导书", // 题型来源
        "researcher": "研究人员", // 研究人员

        "key": "uniquekey", // 唯一标识
    },
    ...
    ], // 预览表中的数据
}
```

## /api/admin-v1/typenames/:uid

### POST

* 描述: 确认保存 `:uid` 对应的表
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

### DELETE

* 描述: 不保存 `:uid` 对应的表并删除临时存储的数据
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

## /api/admin-v1/probtypes

### POST

* 描述: 上传题型确定表 (并没有真正保存)
* 要求: 发送的请求需带有 JWT

请求体:

```
File: 题型确定表.xlsx
```

返回值:

```
{
    "uid": "63362a46-0894-432a-8fc6-8b5594eccdfe", // UUID 用于标识这一操作

    "columns": [{
        "title": "辅助识别码",
        "key": "辅助识别码",
        "dataIndex": "problemID"
    }, {
        "title": "小问序号",
        "key": "小问序号",
        "dataIndex": "subIdx"
    },
    ....
    ], // 预览表中的列

    "data": [{
        "problemID": "12345", // 辅助识别码
        "subIdx": -1, // 小问序号, 如果值为 -1 表示表中的值为空
        "typeName": "样例", // 题型名称

        "key": "uniquekey", // 唯一标识
    },
    ...
    ], // 预览表中的数据
}
```

## /api/admin-v1/probtypes/:uid

### POST

* 描述: 确认保存 `:uid` 对应的表
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

### DELETE

* 描述: 不保存 `:uid` 对应的表并删除临时存储的数据
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

## /api/admin-v1/probmetas

### POST

* 描述: 上传零级属性表 (并没有真正保存)
* 要求: 发送的请求需带有 JWT

请求体:

```
File: 零级属性表.xlsx
```

返回值:

```
{
    "uid": "63362a46-0894-432a-8fc6-8b5594eccdfe", // UUID 用于标识这一操作

    "columns": [{
        "title": "辅助识别码",
        "key": "辅助识别码",
        "dataIndex": "problemID"
    }, {
        "title": "题目来源",
        "key": "题目来源",
        "dataIndex": "source"
    },
    ....
    ], // 预览表中的列

    "data": [{
        "problemID": "12345", // 辅助识别码
        "source": "来源", // 题目来源
        "chapNum": 1, // 章序号
        "sectNum": 1, // 节序号
        "column": "样例", // 栏目名称
        "page": 1, // 页码
        "num": 1, // 题目序号

        "key": "uniquekey", // 唯一标识
    },
    ...
    ], // 预览表中的数据
}
```

## /api/admin-v1/probmetas/:uid

### POST

* 描述: 确认保存 `:uid` 对应的表
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

### DELETE

* 描述: 不保存 `:uid` 对应的表并删除临时存储的数据
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

## /api/admin-v1/hows

### POST

* 描述: 上传出题方式表 (并没有真正保存)
* 要求: 发送的请求需带有 JWT

请求体:

```
File: 出题方式表.xlsx
```

返回值:

```
{
    "uid": "63362a46-0894-432a-8fc6-8b5594eccdfe", // UUID 用于标识这一操作

    "columns": [{
        "title": "辅助识别码",
        "key": "辅助识别码",
        "dataIndex": "problemID"
    }, {
        "title": "出题方式",
        "key": "出题方式",
        "dataIndex": "how"
    },
    ....
    ], // 预览表中的列

    "data": [{
        "problemID": "12345", // 辅助识别码
        "how": "解答体", // 出题方式

        "key": "uniquekey", // 唯一标识
    },
    ...
    ], // 预览表中的数据
}
```

## /api/admin-v1/hows/:uid

### POST

* 描述: 确认保存 `:uid` 对应的表
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

### DELETE

* 描述: 不保存 `:uid` 对应的表并删除临时存储的数据
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

## /api/admin-v1/images

### POST

* 描述: 上传习题和答案图片 (已经保存)
* 要求:
    * 发送的请求需带有 JWT
    * 只能上传 ZIP 且大小不超过 128 MB
    * 图片必须是 PNG 格式

请求体:

```
File: 习题和答案图片.zip
```

返回值:

```
[{
    "_id": "12345", // 题目辅助识别码
    "content": {
        "question": ["12345.png"], // 问题的图片
        "answer": ["12345D.png"], // 答案的图片
    },
}, {
    "_id": "12346", // 题目辅助识别码
    "content": {
        "question": ["12346-1.png", "12346-2.png"],
        "answer": ["12346D-1.png", "12346D-2.png"],
    },
},
...
]
```

## /api/admin-v1/images/:name

### GET

* 描述: 获取 `:name` 对应的图片
* 要求: 发送的请求需带有 JWT (`:name` 需带有扩展名)

返回值:

```
Image
```

### DELETE

* 描述: 删除 `:name` 对应的图片
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

## /api/admin-v1/contents

### POST

* 描述: 上传例题 Markdown 文档和图片 (文档没有真正保存但图片已保存)
* 要求:
    * 发送的请求需带有 JWT
    * 只能上传 ZIP 且大小不超过 128 MB
    * 图片必须是 PNG 格式

请求体:

```
File: 文档和图片.zip
```

返回值:

```
{
    "uid": "63362a46-0894-432a-8fc6-8b5594eccdfe", // UUID 用于标识这一操作

    "data": [{
        "_id": "12345",
        "content": {
            "question": {
                "text": "问题", // 题目的 Markdown 文本
                "images": ["12345.png"] // 题目中的图片
            },
            "answer": {
                "text": "答案", // 答案的 Markdown 文本
                "images": ["12345D.png"] // 答案中的图片
            },
        }
    }, {
        "_id": "12346",
        "content": {
            "question": {
                "text": "问题",
                "images": ["12346-1.png", "12346-2.png"]
            },
            "answer": {
                "text": "答案",
                "images": ["12346D-1.png", "12346D-2.png"]
            },
        }
    },
    ...
    ]
}
```

## /api/admin-v1/contents/:uid

### POST

* 描述: 确认保存 `:uid` 对应的数据
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

### DELETE

* 描述: 不保存 `:uid` 对应的文档并删除临时存储的数据
* 要求: 发送的请求需带有 JWT

返回值:

```
None
```

## /api/admin-v1/extremums

### POST

* 描述: 上传信息序号最值表 (已经保存)
* 要求: 发送的请求需带有 JWT

返回值:

```
{
    "columns": [],
    "data": {
        "章节序号最大值": [{
            "chapNum": "1", // 章序号
            "sectMax": "1", // 节序号MAX
        },
        ...
        ],
        "每学期章序号最值": [{
            "grade": "1", // 年级
            "semester": "1", // 学期
            "chapMin": "1", // 章序号MIN
            "chapMax": "5", // 章序号MAX
        },
        ...
        ],
        "每节知识块序号最大值": [{
            "chapNum": "1", // 章序号
            "sectNum": "1", // 节序号
            "blockMax": "1", // 知识块序号MAX
        },
        ...
        ],
    }
}
```
## /api/admin-v1/probchunk

### POST

* 描述: 上传题型知识块对应数据 (并没有真正保存)
* 要求: 发送的请求需带有 JWT

请求体:

```
file: 题型知识块对应表.xlsx
```

返回值:

```
{
    "uid": "63362a46-0894-432a-8fc6-8b5594eccdfe", // UUID 用于标识这一操作

    "columns": [{
      "title": "章序号",
      "key": "章序号",
      "dataIndex": "chapterNum"
    },
    {
      "title": "节序号",
      "key": "节序号",
      "dataIndex": "sectionNum"
    },
    ....
    ], // 预览表中的列

    "data": [{
      "id": "14f96c9ed73f90bc90b88df06a8d5118",
      "chapterNum": "14",
      "sectionNum": "2",
      "problemType": "常见构造全等三角形的技巧",
      "chunkNum": "-1",
      "chunkName": "-1",
      "key": "14.2"
    },
    ...
    ], // 预览表中的数据
}
```

## /api/admin-v1/probchunk/:uid

### POST

* 描述: 确认保存 `:uid` 对应的表
* 要求: 发送的请求需带有 JWT

返回值:

```
状态200  成功
状态204  无数据
```

### DELETE

* 描述: 不保存 `:uid` 对应的表并删除临时存储的数据
* 要求: 发送的请求需带有 JWT

返回值:

```
状态204 成功
```
