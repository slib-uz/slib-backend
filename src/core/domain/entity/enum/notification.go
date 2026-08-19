package enum

type NotificationSegment int

const (
	SegmentWeb    NotificationSegment = 1
	SegmentMobile NotificationSegment = 2
)

type NotificationTopic string

const (
	TopicNews NotificationTopic = "NEWS"
)
