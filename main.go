package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
)

type s3File struct {
	name   string
	size   uint64
	reader io.Reader
}

type s3Client interface {
	save(file s3File) (err error)
}

type fakeS3Client struct {
}

func (c fakeS3Client) save(file s3File) (err error) {
	fmt.Println("Name: " + file.name)
	fmt.Printf("Size: %d \n", file.size)
	fmt.Print("Content: ")
	_, err = io.CopyN(os.Stdout, file.reader, 4)
	fmt.Println()
	return err
}

type fsS3Client struct {
}

func (c fsS3Client) save(file s3File) (err error) {
	fmt.Println("Name: " + file.name)
	fmt.Printf("Size: %d \n", file.size)
	fmt.Print("Content: ")
	_, err = io.CopyN(os.Stdout, file.reader, 4)
	fmt.Println()
	return err
}

//ZipProcessor reads zip files
type ZipProcessor struct {
	file     string
	s3Client s3Client
}

func (z ZipProcessor) process() error {
	// Open a zip archive for reading.
	r, err := zip.OpenReader(z.file)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
			break
		}
		defer rc.Close()
		s3File := s3File{
			name:   f.Name,
			reader: rc,
			size:   f.UncompressedSize64}

		err = z.s3Client.save(s3File)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	s3 := fakeS3Client{}
	z := ZipProcessor{file: "c:/tmp/arch.zip", s3Client: s3}

	err := z.process()
	if err != nil {
		log.Fatal(err)
	}
}
