package problempdfs

import "time"

// BasicProblemForCreatingFiles 用于生成文档的题目基础类型
type BasicProblemForCreatingFiles struct {
	Type      string `json:"type" bson:"type"`           // 题目类型
	Book      string `json:"book" bson:"book"`           // 书本名称
	Page      int64  `json:"page" bson:"page"`           // 页码
	Column    string `json:"column" bson:"column"`       // 栏目名称
	Idx       int    `json:"idx" bson:"idx"`             // 题目在原书中的题目序号
	ProblemID string `json:"problemID" bson:"problemID"` // 题目识别码
	SubIdxs   []int  `json:"subIdxs" bson:"subIdxs"`     // 该题生成文档时需要用到的小问序号，没有小问的题目是 [-1]
	Full      bool   `json:"full" bson:"full"`           // 是否需要完成整一道题的所有小问
	How       string `json:"how" bson:"how"`             // 出题方式
	Reason    string `json:"reason" bson:"reason"`       // 选题依据
}

// ProblemForCreatingFiles 用于生成文档的题目类型
type ProblemForCreatingFiles struct {
	BasicProblemForCreatingFiles `json:",inline" bson:",inline"`
	CheckProblems                []BasicProblemForCreatingFiles `json:"checkProblems" bson:"checkProblems"` // 检验题目
}

type DetailedProblem struct {
	Book               string `json:"book" db:"book"`
	Page               int64  `json:"page" db:"page"`
	Column             string `json:"column" db:"column"`
	Idx                int    `json:"idx" db:"idx"`
	ProblemID          string `json:"problemID" db:"problemID"`
	SubIdx             int    `json:"subIdx" db:"subIdx"`
	Index              int    `json:"index"`
	Full               bool   `json:"full"`
	Type               string `json:"type" db:"type"`
	Reason             string `json:"reason"`
	sortType           int    // 获取错题时候提交的排序方式
	how                string
	howCode            int
	newestState        bool
	markTimes          int
	wrongTimes         int
	sourceType         int
	newestMarkDuration time.Duration
}

// for convenience, don't want to change every detailedProblem to DetailedProblem
type detailedProblem = DetailedProblem

type problem struct {
	ProblemID string `json:"problemID" bson:"problemID" db:"problemID"`
	SubIdx    int    `json:"subIdx" bson:"subIdx" db:"subIdx"`
}

type typeInfo struct {
	Type     string `db:"typename"`
	Category string `db:"category"`
	Priority int    `db:"priority"`
	Chapter  int    `db:"typeChapter"`
	Section  int    `db:"typeSection"`
}

type detailedProblemOfTypeInfo struct {
	Type     typeInfo
	Problems []detailedProblem
}

type DetailedProblemsType struct {
	Type     string            `json:"type"`
	Problems []detailedProblem `json:"problems"`
}

type detailedProblemsType = DetailedProblemsType

type paperType struct {
	PaperID string `json:"paperID"`
	Name    string `json:"name"`
}
