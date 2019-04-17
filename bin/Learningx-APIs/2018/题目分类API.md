## API

----------

使用 JSON 交换数据。默认当输入有错误时，返回 Status Code 422 (StatusUnprocessableEntity) 和详细错误信息。

### *注：下面所有URL均有前缀/api/admin-v1*

### *GET* /clusterByTarget/analyzeProblem

按目标点分析：分析一个题目

**filters**

```javascript
problemId: string,    // 题目识别码
subIdx: int    // 小问序号
```

**output**

```javascript
{
    typeName: string,  // 该题目题型名称
    totalNum: number,  // 题目数量
    problems: List<{
        problemId: string,      // 题目识别码
        subIdx: number,     // 小问序号
        similarity: float,  // 相似度
        detail: List<{
            step: string,    // 答案步骤
            origPoint: string,      // 原始知识点
            applyPoint: string,     // 应用知识点
            target: string,         // 目标点（小目标）
        }>,    // 题目具体信息（前端用表格形式展示，每一个元素就是一行）
    }>,
}
```

### *GET* /clusterByTarget/clusterProblem

按目标点：对未分类的题目进行分类展示分类结果

**filters**

```javascript
chapter: number,   // 章
section: number,   // 节
```

**output**

```javascript
{
    silhouetteIdx: float,  // 轮廓系数
    msg: string,  // 轮廓系数对应提示信息
    clusterResult: List<{
        typeName: string,  // 该题型名称
        priority: number,  // 该题型学习顺序
        totalNum: number,  // 题目数量
        radius: float,  // 聚类半径
        msg: string,  // 聚类半径对应的提示信息
        problems: List<{
            problemId: string,  // 题目识别码
            subIdx: number,  // 小问序号
            similarity: float,  // 相似度
            detail: List<{
                step: string,    // 答案步骤
                origPoint: string,      // 原始知识点
                applyPoint: string,     // 应用知识点
                target: string,         // 目标点（小目标）
            }>,    // 题目具体信息（前端用表格形式展示，每一个元素就是一行）
        }>,
    }>,     // 分类结果
}
```

### *POST* /clusterByTarget/clusterResult

提交分类结果

**input**

```javascript
List<{
    typeName: string,  // 该题型名称
    priority: number,  // 该题型学习顺序
    problems: List<{
        problemId: string,  // 题目识别码
        subIdx: number,  // 小问序号
    }>,
}>
```

**output**

```javascript
"Clustering succeeded."
```
