// Code generated by xgen. DO NOT EDIT.

package repository

// WorkerRequest ...
type WorkerRequest struct {
	RequestId  string `xml:"RequestId"`
	Hash       string `xml:"Hash"`
	Alphabet   string `xml:"Alphabet"`
	MaxLength  int    `xml:"MaxLength"`
	PartNumber int    `xml:"PartNumber"`
	PartCount  int    `xml:"PartCount"`
}
