package handler

var Errors map[string]string

func init() {
	Errors = make(map[string]string, 16)
	Errors["InvalidBucketName"] = "The specified bucket is not valid"
	Errors["InvalidGetBucketVersion"] = "The get bucket version is not supported"
	Errors["InvalidParameter"] = "The specified parameter is not valid"
	Errors["NoSuchBucket"] = "The specified bucket does not exist"
	Errors["BucketNotEmpty"] = "The bucket you tried to delete is not empty"
	Errors["NoSuchBucket"] = "The specified bucket does not exist"
	Errors["InternalError"] = "We encountered an internal error. Please try again"
	Errors["TooManyBuckets"] = "You have attempted to create more buckets than allowed"
	Errors["BucketAlreadyExists"] = "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again"
	Errors["NoSuchKey"] = "The specified key does not exist"
	Errors["BadRequest"] = "Bad request"
	Errors["NoSuchUpload"] = "The specified multipart upload does not exist. The upload ID might be invalid, or the multipart upload might have been aborted or completed"
	Errors["UploadAborted"] = "The specified multipart upload has been aborted"
	Errors["InvalidPartOrder"] = "The list of parts was not in ascending order.Parts list must specified in order by part number"
	Errors["InvalidPart"] = "One or more of the specified parts could not be found. The part might not have been uploaded, or the specified entity tag might not have matched the part's entity tag"
}

func ValidError(key string) bool {
	_, ok := Errors[key]
	return ok
}

func ErrorMessage(key string) (string, bool) {
	msg, ok := Errors[key]
	return msg, ok
}
