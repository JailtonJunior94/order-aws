package messaging

import "time"

type (
	EventNotifications struct {
		Records []Record `json:"Records"`
	}

	Record struct {
		EventVersion      string            `json:"eventVersion"`
		EventSource       string            `json:"eventSource"`
		AwsRegion         string            `json:"awsRegion"`
		EventTime         time.Time         `json:"eventTime"`
		EventName         string            `json:"eventName"`
		UserIdentity      UserIdentity      `json:"userIdentity"`
		RequestParameters RequestParameters `json:"requestParameters"`
		ResponseElements  ResponseElements  `json:"responseElements"`
		S3                S3                `json:"s3"`
		GlacierEventData  GlacierEventData  `json:"glacierEventData"`
	}

	UserIdentity struct {
		PrincipalID string `json:"principalId"`
	}

	RequestParameters struct {
		SourceIPAddress string `json:"sourceIPAddress"`
	}

	ResponseElements struct {
		XAmzRequestID string `json:"x-amz-request-id"`
		XAmzID2       string `json:"x-amz-id-2"`
	}

	S3 struct {
		S3SchemaVersion string `json:"s3SchemaVersion"`
		ConfigurationID string `json:"configurationId"`
		Bucket          Bucket `json:"bucket"`
		Object          Object `json:"object"`
	}

	Bucket struct {
		Name          string       `json:"name"`
		OwnerIdentity UserIdentity `json:"ownerIdentity"`
		Arn           string       `json:"arn"`
	}

	Object struct {
		Key       string `json:"key"`
		Size      int64  `json:"size"`
		ETag      string `json:"eTag"`
		VersionID string `json:"versionId"`
		Sequencer string `json:"sequencer"`
	}

	GlacierEventData struct {
		RestoreEventData RestoreEventData `json:"restoreEventData"`
	}

	RestoreEventData struct {
		LifecycleRestorationExpiryTime time.Time `json:"lifecycleRestorationExpiryTime"`
		LifecycleRestoreStorageClass   string    `json:"lifecycleRestoreStorageClass"`
	}
)
