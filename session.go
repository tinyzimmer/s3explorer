/**
This file is part of s3explorer.

s3explorer is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

s3explorer is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with s3explorer.  If not, see <https://www.gnu.org/licenses/>.
**/

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Session struct {
	S3Service      *s3.S3
	ConfigProvider *client.ConfigProvider
	Buckets        []*s3.Bucket
}

type BucketWithDisplay struct {
	bucket        *s3.Bucket
	displayString string
	region        string
}

func (s S3Session) DownloadObject(bucket BucketWithDisplay, node *Node, dest string) (err error) {

	log.Printf("\nDownload Call:\n\tBucket: %+v\n\tNode: %+v\n\tDestination: %s\n", bucket, node, dest)

	// sanity check

	if node.S3Object == nil {
		err = errors.New(fmt.Sprintf("No s3 object associated with node: %+v\n", node))
		log.Println(err)
		return
	}

	// Check if dest already exists

	if FileExists(dest) {
		log.Println("Removing pre-existing file")
		os.Remove(dest)
	}

	// Recursively create needed directories

	path, _ := filepath.Split(dest)
	log.Printf("Recursively creating directory: %s\n", path)
	err = os.MkdirAll(path, DEFAULT_DIRECTORY_MODE)
	if err != nil {
		return
	}

	log.Printf("Creating destination path: %s\n", dest)

	// Open a file

	file, err := os.Create(dest)
	if err != nil {
		return
	}
	defer file.Close()

	log.Printf("Getting downloader for region: %s\n", bucket.region)

	// Create a downloader with the s3 client and custom options

	downloader := s3manager.NewDownloaderWithClient(s.S3Service, func(d *s3manager.Downloader) {
		d.PartSize = 64 * 1024 * 1024 // 64MB per part
	})

	n, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: bucket.bucket.Name,
		Key:    node.S3Object.Key,
	})

	if err != nil {
		log.Printf("failed to download file: %v\n", err)
	}

	log.Printf("file downloaded, %d bytes\n", n)
	return
}

func (s S3Session) GetBucketObjects(bucket BucketWithDisplay) (objects []*s3.Object, err error) {

	// For a given bucket, retrieve a list of all its objects

	log.Printf("Listing Objects for Bucket: %s\n", bucket.displayString)
	err = s.S3Service.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: bucket.bucket.Name,
	},
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			objects = append(objects, page.Contents...)
			return true
		})
	return
}

func (s S3Session) GetBucketWithDisplayStrings() (bucketStrings []BucketWithDisplay, err error) {

	// Get all buckets, and attach a display string to it

	buckets, err := s.GetBucketListing()
	if err != nil {
		return
	}
	for _, bucket := range buckets {
		region, err := s.GetBucketRegion(bucket)
		if err != nil {
			RenderError(err.Error())
		}
		displayBucket := BucketWithDisplay{
			bucket:        bucket,
			displayString: fmt.Sprintf("%s (%s)", *bucket.Name, region),
			region:        region,
		}
		bucketStrings = append(bucketStrings, displayBucket)
	}
	return
}

func (s S3Session) GetBucketRegion(bucket *s3.Bucket) (region string, err error) {

	// Get the region for a bucket

	log.Printf("Retrieving region for bucket: %s\n", *bucket.Name)
	ctx := context.Background()
	region, err = s3manager.GetBucketRegionWithClient(ctx, s.S3Service, *bucket.Name)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			region = "unknown"
		} else {
			log.Printf("Unknown Error: %s\n", err.Error())
		}
	}
	return
}

func (s S3Session) GetBucketListing() (buckets []*s3.Bucket, err error) {

	log.Println("Listing Buckets")
	resp, err := s.S3Service.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return
	}
	buckets = resp.Buckets
	log.Println("Retrieved bucket list")
	return
}

func (s S3Session) RefreshBucketListing() (err error) {
	buckets, err := s.GetBucketListing()
	log.Println("Refreshed Bucket Listing")
	if err != nil {
		return
	}
	s.Buckets = buckets
	return
}

func InitSession(region string) (s3session S3Session, err error) {

	creds, err := getCreds()
	if err != nil {
		return
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))

	s3session.S3Service = s3.New(sess)
	log.Printf("Connected to S3 in Region: %s\n", region)
	return
}

func getCreds() (creds *credentials.Credentials, err error) {
	sess := session.Must(session.NewSession())
	creds = credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		})
	_, err = creds.Get()
	if err != nil {
		return
	}
	log.Println("Got AWS Credentials")
	return
}
