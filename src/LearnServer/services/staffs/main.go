package staffs

import (
	"LearnServer/services/staffs/dataArchive"
	"LearnServer/services/staffs/errorRateAnalysis"
	"LearnServer/services/staffs/scoreAnalysis"
	"LearnServer/services/staffs/personCenter"
	"LearnServer/utils"
	"github.com/labstack/echo"
)

// RegisterStaffsApis 注册staffs下的路由
func RegisterStaffsApis(e *echo.Group) {

	groupPublic := e.Group("/staffs")
	{
		// pulic API
		groupPublic.POST("/login/", loginHandler)
	}

	groupPrivate := e.Group("/staffs", utils.JwtMiddleware())
	{
		// private API (只能给已经登录的用户访问)
		groupPrivate.POST("/archiveData/", dataArchive.ArchiveDataHandler)

		groupPrivate.POST("/me/logout/", logoutHandler)
		groupPrivate.GET("/me/testLogin/", testLoginHandler)
		groupPrivate.PUT("/me/password/", changePasswordHandler)
		groupPrivate.GET("/me/profile/", retriveStaffProfileHandler)
		groupPrivate.GET("/schools/", getSchoolsHandler)
		groupPrivate.POST("/schools/", uploadSchoolHandler)
		groupPrivate.POST("/students/", uploadStudentsHandler)
		groupPrivate.POST("/students/addOne/", addOneNewStudentHandler)
		groupPrivate.DELETE("/students/:learnID/", deleteStudentHandler)
		groupPrivate.GET("/students/:learnID/markedPapers/", getMarkedPapersHandler)
		groupPrivate.GET("/students/:learnID/notMarkedPapers/", getNotMarkedPapersHandler)

		groupPrivate.POST("/students/:learnID/getWrongProblems/", getWrongProblemsHandler)
		groupPrivate.POST("/students/:learnID/documents/", createDocumentHandler)

		groupPrivate.POST("/students/markTasks/", createMarkTasksHandler)
		groupPrivate.POST("/students/getDocumentZip/", getPackedFileHandler)

		groupPrivate.POST("/studentFile/", previewStudentFileHandler)
		groupPrivate.DELETE("/studentFile/:uid/", deleteStudentTmpFileHandler)
		groupPrivate.GET("/books/search/", searchBookHandler)
		groupPrivate.GET("/papers/search/", searchPaperHandler)
		groupPrivate.POST("/classes/addBooks/", addBookHandler)
		groupPrivate.POST("/classes/addPapers/", addPaperHandler)
		groupPrivate.GET("/classes/books/", getBooksOfClassHandler)
		groupPrivate.GET("/classes/papers/", getPapersOfClassHandler)
		groupPrivate.GET("/classes/students/", getStudentsHandler)
		groupPrivate.POST("/classes/deleteBooks/", deleteBookFromClassHandler)
		groupPrivate.POST("/classes/deletePapers/", deletePaperFromClassHandler)

		groupPrivate.GET("/students/:learnID/paperProblems/", getPaperProblemsHandler)
		groupPrivate.GET("/students/:learnID/bookProblems/", getBookProblemsHandler)
		groupPrivate.POST("/students/:learnID/problems/", uploadProblemResultHandler)
		groupPrivate.GET("/students/:learnID/markTasks/", listMarkTasksHandler)
		groupPrivate.GET("/students/:learnID/markTasks/:time/", getMarkTaskHandler)
		groupPrivate.DELETE("/students/:learnID/markTasks/:time/", deleteMarkTaskHandler)

		groupPrivate.POST("/students/getProblemRecords/", getProblemRecordsHandler)

		groupPrivate.GET("/students/:learnID/", retriveStudentDetailHandler)
		groupPrivate.PUT("/students/:learnID/productID/", updateStudentProductIDHandler)
		groupPrivate.PUT("/students/:learnID/", updateStudentInfoHandler)

		groupPrivate.POST("/products/", uploadProductHandler)
		groupPrivate.GET("/products/:productID/", retriveProductHandler)
		groupPrivate.PUT("/products/:productID/", updateProductHandler)
		groupPrivate.PUT("/products/:productID/status/", updateProductStatusHandler)
		groupPrivate.GET("/products/", listProductHandler)

		groupPrivate.PUT("/classes/students/level/", updateStudentLevelHandler)
		groupPrivate.GET("/classes/productID/", getClassProductsHandler)
		groupPrivate.PUT("/classes/productID/", updateClassProductsHandler)
		groupPrivate.GET("/classes/totalLevel/", getClassTotalLevelHandler)
		groupPrivate.POST("/classes/totalLevel/", updateClassTotalLevelHandler)

		groupPrivate.POST("/batchDownloads/", createBatchDownloadTask)
		groupPrivate.DELETE("/batchDownloads/:batchID/", deleteBatchDownloadTask)
		groupPrivate.GET("/batchDownloads/", listBatchDownloadTasks)

		groupPrivate.POST("/classes/getErrorRateAnalysis/", errorRateAnalysis.GetErrorRateHandler)
		// groupPrivate.POST("/classes/getPracticeProblems/", errorRateAnalysis.GetPracticeProblemsHandler)
		// groupPrivate.POST("/classes/practiceProblems/getProblemsFile/", errorRateAnalysis.GetProblemsFileHandler)
		// groupPrivate.POST("/classes/practiceProblems/getAnswersFile/", errorRateAnalysis.GetAnswersFileHandler)
		// groupPrivate.POST("/classes/practiceProblems/getProbsAnsFilesZip/", errorRateAnalysis.GetPackedFileHandler)

		groupPrivate.GET("/info/chapsSects/", getChaptersSectionsHandler)
		groupPrivate.GET("/classes/targets/", getClassTargetsHandler)
		groupPrivate.POST("/classes/targets/", addClassTargetsHandler)
		groupPrivate.DELETE("/classes/targets/", deleteClassTargetsHandler)

		groupPrivate.GET("/classes/bookProblems/", getUnassignedBookProblemsHandler)
		groupPrivate.POST("/classes/assignments/", uploadAssignmentsHandler)

		groupPrivate.GET("/classes/markTasks/", listClassMarkTasksHandler)
		groupPrivate.DELETE("/classes/markTasks/", deleteBundleMarkTasksHandler)

		groupPrivate.POST("/templates/", uploadTemplateHandler)
		groupPrivate.GET("/templates/", listTemplatesHandler)
		groupPrivate.GET("/templates/:templateID/", retriveTemplateHandler)
		groupPrivate.PUT("/templates/:templateID/", updateTemplateHandler)
		groupPrivate.DELETE("/templates/:templateID/", deleteTemplateHandler)
		groupPrivate.GET("/templates/:templateID/preview/", previewTemplateHandler)

		groupPrivate.POST("/classes/examScores/", uploadExamScoresHandler)
		groupPrivate.POST("/classes/scoreFile/", uploadScoreFileHandler)
		groupPrivate.GET("/classes/examScores/", getExamScoresHandler)
		groupPrivate.GET("/classes/papersForMarkScore/", getPapersOfClassForMarkingScoreHandler)

		groupPrivate.GET("/info/problemTypes/", getProblemTypesHandler)
		groupPrivate.GET("/info/problems/", getProblemsOfTypeHandler)
		groupPrivate.POST("/info/getProblemsZip/", getProblemDocsZipHandler)

		groupPrivate.GET("/info/bookProblems/", getBookProblemsInfoHandler)
		groupPrivate.GET("/info/paperProblems/", getPaperProblemsInfoHandler)
		groupPrivate.POST("/classes/problemsLearned/methodOne/", uploadProblemsLearnedMethodOneHandler)
		groupPrivate.POST("/classes/problemsLearned/methodTwo/", uploadProblemsLearnedMethodTwoHandler)

		groupPrivate.GET("/info/knowledgePoint/", getKnowledgePointOfChapSectHandler)
		groupPrivate.POST("/classes/knowledgeLearned/", uploadKnowledgeLearnedHandler)

		groupPrivate.POST("/classes/semester/", uploadSemesterHandler)
		groupPrivate.GET("/classes/semester/", getSemesterHandler)

		groupPrivate.GET("/classes/examAnalysis/average/", scoreAnalysis.GetAverageAnalysisHandler)
		groupPrivate.GET("/classes/examAnalysis/rankingLevelAverage/", scoreAnalysis.GetRankingLevelAverageAnalysisHandler)
		groupPrivate.GET("/classes/examAnalysis/scoreProportion/", scoreAnalysis.GetScoreProportionAnalysisHandler)
		groupPrivate.GET("/classes/examAnalysis/studentScore/", scoreAnalysis.GetStudentScoreAnalysisHandler)
		groupPrivate.GET("/classes/examAnalysis/studentRanking/", scoreAnalysis.GetStudentRankingAnalysisHandler)
		groupPrivate.POST("/classes/examAnalysis/thoughts/", scoreAnalysis.UploadExamThoughtsHandler)

		groupPrivate.PUT("/me/telphone/", personCenter.BindTelPhoneNumber)
		groupPrivate.PUT("/me/manageClasses/:classIndex/Update/", personCenter.UpdateClassInfo)
		groupPrivate.DELETE("/me/manageClasses/:classIndex/Delete/", personCenter.DeleteClassInfo)
		groupPrivate.POST("/me/manageClasses/Add/", personCenter.AddClassInfo)
		groupPrivate.PUT("/me/baseInfo/", personCenter.UpdatePersonInfo)
		
	}
}
