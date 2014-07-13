# aws - A client for various AWS services, written in Golang.

# Disclaimer

Currently uses `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` ENV variables as id/secret sources.

TODO: remove this fetching of ENV variables from this library - that should be in the domain of the user application!

# Install

	$ go get github.com/bwilkins/aws

## Use

NB: Application will require `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` to be set in the environment (currently).

	// NOTE: This project is still a work in progress.  For those of you who
	// see where I'm going with it and want to help, please do!

	package main

	import (
		"github.com/bwilkins/aws/opsworks"
		"log"
	)

	// Assumes AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are set in ENV.
	func main() {
                request := opsworks.DescribeInstancesRequest{StackId: "your-stack-id"}
	        response, err := opsworks.DescribeInstances(request)

		log.Printf("%v", response)
	}


## LICENCES

Copyright (C) 2014 by Brett Wilkins (@bjmaz) `<brett@brett.geek.nz>`

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

