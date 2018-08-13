# s3explorer
Terminal Based S3 File Explorer

### Building From Source


**Dependencies**
  - github.com/aws/aws-sdk-go
  - github.com/gizak/termui


```bash
$> go get github.com/aws/aws-sdk-go
$> go get github.com/gizak/termui
$> go get github.com/tinyzimmer/s3explorer
$> go install github.com/tinyzimmer/s3explorer
```

### Usage

```bash
$> s3explorer <-d [debug file]>
```

The program will list all of the S3 Buckets you have access to and present them in a file explorer format. You can descend into the buckets and directories therein with your keyboard. 
