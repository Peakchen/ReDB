package students

import (
	// "LearnServer/services/students/problemSorting"
	// "LearnServer/services/students/problempdfs"
	// "LearnServer/services/students/tasks"
	// "LearnServer/utils"
	"LearnServer/services/students/problemSorting"
	"LearnServer/services/students/problempdfs"
	"LearnServer/services/students/tasks"
	"LearnServer/utils"
	"github.com/labstack/echo"
)

// RegisterStudentsApis 注册students下的路由
func RegisterStudentsApis(e *echo.Group) {

	groupPublic := e.Group("/students")
	{
		// pulic API
		groupPublic.POST("/login/", loginHandler)
	}

	groupPrivate := e.Group("/students/me", utils.JwtMiddleware())
	{
		// private API (只能给已经登录的用户访问)
		groupPrivate.POST("/logout/", logoutHandler)
		groupPrivate.PUT("/password/", changePasswordHandler)
		groupPrivate.GET("/schools/", getSchoolsHandler)
		groupPrivate.GET("/profile/", getProfileHandler)
		groupPrivate.PATCH("/profile/", updateProfileHandler)
		groupPrivate.GET("/books/", getBooksHandler)
		groupPrivate.GET("/problems/", getProblemsByPosHandler)
		groupPrivate.POST("/problems/", uploadProblemResultHandler)
		groupPrivate.GET("/info/", getInfoHandler)
		groupPrivate.GET("/problemsSortByType/", problemSorting.GetProblemsSortByTypeHandler)
		groupPrivate.GET("/problemsSortByTime/", problemSorting.GetProblemsSortByTimeHandler)
		groupPrivate.GET("/wrongProblemsInfo/", problempdfs.GetNewestWrongProblemsByChapSectHandler)
		groupPrivate.POST("/getProblemsFile/", problempdfs.GetProblemsFileHandler)
		groupPrivate.POST("/getAnswersFile/", problempdfs.GetAnswersFileHandler)
		groupPrivate.POST("/getPointsFile/", problempdfs.GetPointsFileHandler)
		groupPrivate.GET("/checkProblemsInfo/", problempdfs.GetCheckProblemsByChapSectHandler)
		groupPrivate.POST("/problemsChecked/", problempdfs.UploadProblemCheckedResultHandler)
		groupPrivate.POST("/newestWrongProblems/", problempdfs.GetNewestWrongProblemsByBookPageHandler)
		groupPrivate.POST("/onceWrongProblems/", problempdfs.GetOnceWrongProblemsByBookPageHandler)
		groupPrivate.POST("/checkProblemsForAll/", problempdfs.GetAllCheckProblemsByBookPageHandler)
		groupPrivate.POST("/checkProblemsForKnown/", problempdfs.GetKnownCheckProblemsByBookPageHandler)
		groupPrivate.POST("/checkProblemsForStillWrong/", problempdfs.GetNewestCheckProblemsByBookPageHandler)
		groupPrivate.GET("/markTasks/", tasks.ListMarkTasksHandler)
		groupPrivate.POST("/markTasks/", tasks.CreateMarkTaskHandler)
		groupPrivate.GET("/markTasks/:time/", tasks.GetMarkTaskHandler)
		groupPrivate.DELETE("/markTasks/:time/", tasks.DeleteMarkTaskHandler)
		groupPrivate.GET("/notMarkedPapers/", getNotMarkedPapersHandler)
		groupPrivate.GET("/markedPapers/", getMarkedPapersHandler)
		groupPrivate.GET("/paperProblems/", getProblemsByPaperIDHandler)
		groupPrivate.POST("/onceWrongPaperProblems/", problempdfs.GetOnceWrongPaperProblemsHandler)
		groupPrivate.POST("/newestWrongPaperProblems/", problempdfs.GetNewestWrongPaperProblemsHandler)
		groupPrivate.GET("/wrongProblems/", problempdfs.GetWrongProblemsForLearning)
		groupPrivate.GET("/problemRecords/", getProblemRecordsHandler)
		groupPrivate.POST("/problemFileState/", problempdfs.SetWrongProblemFileStateHandler)
		groupPrivate.GET("/problemFileState/", problempdfs.GetWrongProblemFileStateHandler)
		groupPrivate.GET("/lastFileURLs/", problempdfs.GetLastFileURLs)
		groupPrivate.GET("/lastWrongProblems/", problempdfs.GetLastWrongProblemsHandler)
		groupPrivate.POST("/learningPackage/", setLearningPackageHandler)
		groupPrivate.GET("/learningPackage/", getLearningPackageHandler)
		groupPrivate.GET("/products/", getProductsHandler)
	}
}
