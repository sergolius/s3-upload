# AWS S3 Upload
CLI Tool what upload folder to AWS S3

### Usage
Configured from environment variables
```
S3_BUCKET="" S3_REGION="" S3_ID="" S3_SECRET="" S3_LOG="" s3-upload {folderPath1} ... {folderPathN}
```

#### Log Level (S3_LOG):
- warn (default)
- debug
- trace, info

#### License
[MIT](https://github.com/sergolius/s3-upload/blob/master/LICENSE)