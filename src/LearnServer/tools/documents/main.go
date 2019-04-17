package documents

// requestMessageType 发送给 channel 的数据类型
type requestMessageType struct {
	PostData    interface{}                // 请求发送的数据
	RequestType int                        // 请求的类型，参见 requests.go 中的枚举值
	CallBack    func(DocumentURLType, int) // 请求完成的回调函数
}

// requestChannelList 全局变量，channel 数组，下标越小的 channel 优先级越高
var requestChannelList = [11]chan requestMessageType{}

// notificationChannel 用于触发控制器对 requestChannelList 的检查
var notificationChannel = make(chan bool)

func init() {
	for i := range requestChannelList {
		requestChannelList[i] = make(chan requestMessageType, 100000)
	}

	go func() {
		// 生成文档请求的调度器
		for {
			select {
			case <-notificationChannel:
				// 通过 notificationChannel 来触发检查，方便之后先处理 requestChannelList 中的哪个channel，实现优先级
				allRequestsDone := false
				for !allRequestsDone {
					sentRequest := false
					for i := range requestChannelList {
						if len(requestChannelList) > 0 {
							requestMsg := <-requestChannelList[i]
							requests[requestMsg.RequestType](requestMsg.PostData, requestMsg.CallBack)
							sentRequest = true
							// 发送一个请求之后，确保重新从第一个 channel 去获取数据，保证按照优先级顺序
							break
						}
					}
					if !sentRequest {
						allRequestsDone = true
					}
				}
			}
		}
	}()
}

// PutToRequestQueue 将数据放入待请求队列(channel)中，等待发送
func PutToRequestQueue(postData interface{}, requestType int, priority int, callBack func(DocumentURLType, int)) {
	go func() {
		requestMsg := requestMessageType{
			PostData:    postData,
			RequestType: requestType,
			CallBack:    callBack,
		}

		requestChannelList[priority] <- requestMsg
		notificationChannel <- true
	}()
}

// GetChannelWaitingCount 获取特定 channel 正在等待的个数
func GetChannelWaitingCount(priority int) int {
	return len(requestChannelList[priority])
}
