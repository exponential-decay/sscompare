package main

import (
   "fmt"
)

var err_sscomp = fmt.Errorf("ssdeep undefined comparison error on hashes.")
var err_sha1_file1 = fmt.Errorf("Cannot return sha1 for file1")
var err_sha1_file2 = fmt.Errorf("Cannot return sha1 for file2")
