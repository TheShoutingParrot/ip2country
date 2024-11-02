# ip2country

**ip2country** is a golang package to find out IP's origin country. It uses [db-ip.com's csv](https://db-ip.com/db/download/country) file to provide answers.

This fork supports both IPv4 and IPv6 addresses. This forks is not compatible with the original package!

## Original repo by Mostafa Asgari

The [original repo's](https://github.com/mostafa-asg/ip2country) was made by Moustafa Asgari. It's under the MIT license.

The original license notice is the following:

```
MIT License

Copyright (c) 2018 Mostafa Asgari

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## Install

```
go get -u github.com/theshoutingparrot/ip2country
```

## Usage
```Go
package main

import (
	"github.com/mostafa-asg/ip2country"
)

func main() {
	ips := ip2country.Load( PATH_TO_CSV FILE )
	println(ips.GetCountry("2.179.6.12"))
	println(ips.GetCountry("172.217.18.14"))
	println(ips.GetCountry("217.160.123.58"))
}

```
