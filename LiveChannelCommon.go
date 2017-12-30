package OssLiveChannel

type LiveChannelConfiguration struct {
	Description string `xml:"Description,omitempty"`
	Status      string //enabled、disabled
	Target      Target
}
type Target struct {
	Type         string //HLS
	FragDuration string //当Type为HLS时，指定每个ts文件的时长（单位：秒），取值范围为[1, 100]的整数。
	FragCount    string //当Type为HLS时，指定m3u8文件中包含ts文件的个数，取值范围为[1, 100]的整数。
	PlayListName string `xml:"PlayListName,omitempty"`
}

type CreateLiveChannelResult struct {
	PublishUrls PublishUrls
	PlayUrls    PlayUrls
}

type PublishUrls struct {
	Url string
}
type PlayUrls struct {
	Url string
}

type ListLiveChannelResult struct {
	Prefix, Marker, MaxKeys, IsTruncated, NextMarker string
	LiveChannel                                      LiveChannel
}

type LiveChannel struct {
	Name, Description, Status, LastModified string
	PublishUrls                             PublishUrls
	PlayUrls                                PlayUrls
}

type LiveChannelHistory struct {
	LiveRecord []LiveRecord
}

type LiveRecord struct {
	StartTime, EndTime, RemoteAddr string
}
type LiveRespError struct {
	Code, Message, RequestId, HostId, ChannelId string
}
