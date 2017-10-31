package mysql

import (
	"galaxy-s3-gateway/db"
)

type ID struct {
	ID    int64 `gorm:"AUTO_INCREMENT"`
	Value int   `gorm:"type:int;unsigned;not null"`
}

type MysqlBucket struct {
	BucketID   int64 `gorm:"type:bigint(20);unsigned;not null"`
	KeyMd5High int64 `gorm:"type:bigint(20);not null"`
	db.Bucket
}

func toMysqlBucket(id int64, orig *db.Bucket) *MysqlBucket {
	return &MysqlBucket{
		BucketID:   id,
		KeyMd5High: MD5High(MD5Sum(orig.BucketName)),
		Bucket: db.Bucket{
			BucketName: orig.BucketName,
			UserID:     orig.UserID,
			ACL:        orig.ACL,
			CreateTime: orig.CreateTime,
			// 使用mysql存储默认不支持Version功能
			VersionEnabled: false,
		}}
}

func toBucket(orig *MysqlBucket) *db.Bucket {
	return &db.Bucket{
		BucketName: orig.BucketName,
		UserID:     orig.UserID,
		ACL:        orig.ACL,
		CreateTime: orig.CreateTime,
		// 使用mysql存储默认不支持Version功能
		VersionEnabled: false,
	}
}

type MysqlObject struct {
	db.Object
	ObjectID      int64  `gorm:"type:bigint(20);unsigned;NOT NULL;index:IDX_ObjectID"`
	KeyMd5High    int64  `gorm:"type:bigint(20);NOT NULL;primary_key"`
	KeyMd5Low     int64  `gorm:"type:bigint(20);NOT NULL;primary_key"`
	// Md5High       int64  `gorm:"type:bigint(20);NOT NULL;index:Md5High"`
	// Md5Low        int64  `gorm:"type:bigint(20);NOT NULL"`
	ConflictFlag  int8   `gorm:"type:tinyint(3);unsigned;NOT NULL;default:0;primary_key"`
	BucketID      int64  `gorm:"type:bigint(20);unsigned;NOT NULL;index:IDX_BucketID_LastModified"`
	DigestVersion int8   `gorm:"type:smallint(5);NOT NULL"`
}

func toMysqlObject(id, bucketId int64, orig *db.Object) *MysqlObject {
	return &MysqlObject{
		ObjectID:   id,
		KeyMd5High: MD5High(MD5Sum(orig.ObjectName)),
		KeyMd5Low:  MD5Low(MD5Sum(orig.ObjectName)),
		// Md5High:    MD5High([]byte(orig.Etag)),
		// Md5Low:     MD5Low([]byte(orig.Etag)),
		DigestVersion: 0,

		BucketID:   bucketId,
		Object: db.Object{
			ObjectName: orig.ObjectName,
			Bucket:     orig.Bucket,
			Size:       orig.Size,
			Etag:       orig.Etag,
			LastModified: orig.LastModified,
			Fid: orig.Fid,
			Meta: orig.Meta,
			MultipartUpload: orig.MultipartUpload,
			PartSize: orig.PartSize,
			// mysql存储暂不支持version
			Version:  "",
			DeleteMarker: false,
		}}
}

func toObject(orig *MysqlObject) *db.Object {
	return &db.Object {
		Fid:        orig.Fid,
		Meta:       orig.Meta,
		Size:       orig.Size,
		Etag:       orig.Etag,
		Bucket:     orig.Bucket,
		PartSize:   orig.PartSize,
		// mysql存储暂不支持version
		Version:    "",
		DeleteMarker: false,
		ObjectName:   orig.ObjectName,
		LastModified: orig.LastModified,
		MultipartUpload: orig.MultipartUpload,
	}
}

type HistoryObject struct{
	db.Object
	ObjectID      int64  `gorm:"type:bigint(20);unsigned;NOT NULL;primary_key”`
	KeyMd5High    int64  `gorm:"type:bigint(20);NOT NULL"`
	KeyMd5Low     int64  `gorm:"type:bigint(20);NOT NULL"`
	// Md5High       int64  `gorm:"type:bigint(20);NOT NULL"`
	// Md5Low        int64  `gorm:"type:bigint(20);NOT NULL"`
	BucketID      int64  `gorm:"type:bigint(20);unsigned;NOT NULL"`
	DigestVersion int8   `gorm:"type:smallint(5);NOT NULL"`
	HistoryTime   int64  `gorm:"type:bigint(20);unsigned;NOT NULL;index:IDX_HTime"`
	DeleteHint    int8   `gorm:"type:tinyint(3);unsigned;not null"`
}

type ObjectList struct {
	BucketID   int64  `gorm:"type:bigint(20);unsigned;NOT NULL;primary_key"`
	ObjectName string `gorm:"type:varchar(512);COLLATE utf8_bin;NOT NULL;primary_key"`
}

type MysqlUploadInfo struct {
	db.UploadInfo
	BucketID       int64  `gorm:"type:bigint(20);unsigned;NOT NULL"`
}

func toMysqlUpload(bucketId int64, orig *db.UploadInfo) *MysqlUploadInfo {
	return &MysqlUploadInfo {
		BucketID: bucketId,
		UploadInfo: db.UploadInfo{
			UploadID:   orig.UploadID,
			StartTime:  orig.StartTime,
			Bucket:     orig.Bucket,
			Object:     orig.Object,
			UserID:     orig.UserID,
			IsAbort:    orig.IsAbort,
			Meta:       orig.Meta,
		},
	}
}

func toUpload(orig *MysqlUploadInfo) *db.UploadInfo {
	return &db.UploadInfo {
		UploadID:   orig.UploadID,
		StartTime:  orig.StartTime,
		Bucket:     orig.Bucket,
		Object:     orig.Object,
		UserID:     orig.UserID,
		IsAbort:    orig.IsAbort,
		Meta:       orig.Meta,
	}
}

type AbandonUploadPart struct {
	DocID string `gorm:"type:varchar(64);not null;primary_key"`
	CTime int64  `gorm:"type:bigint(20);unsigned;NOT NULL"`
}
