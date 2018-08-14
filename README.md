# s3explorer
Terminal Based S3 File Explorer

### Download

There are pre-compiled binaries in the [releases](https://github.com/tinyzimmer/s3explorer/releases) section.

#### AWS Credentials

   - Refer to the AWS documentation to configure your credentials.
   `s3explorer` loads credentials in the following order:
     - Environment credentials
     - Shared credentials file (e.g. `$HOME/.aws/credentials`)
     - EC2 Instance Profile

### Building From Source


**Dependencies**
  - github.com/aws/aws-sdk-go
  - github.com/gizak/termui


```bash
# Using `go get`
$> go get github.com/aws/aws-sdk-go
$> go get github.com/gizak/termui
$> go get github.com/tinyzimmer/s3explorer

# From git or source code tarball
$> git clone https://github.com/tinyzimmer/s3explorer # or download archive
$> cd s3explorer
$> go build .
```

### Usage

```bash
$> s3explorer <-d [debug file]>
```

The program will list all of the S3 Buckets you have access to and present them in a file explorer format. You can descend into the buckets and directories therein with your keyboard.
