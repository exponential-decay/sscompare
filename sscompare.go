//ssdeep c binding note: https://godoc.org/github.com/dutchcoders/gossdeep#HashFilename
package main

import (
	"os"
	"fmt"
   "log"
   "flag"
   "path/filepath"
   "github.com/michielbuddingh/spamsum"
)

var fuzz bool
var compare bool
var file1 string
var file2 string
var hash1 string
var hash2 string

func init() {
   flag.BoolVar(&fuzz, "fuzz", true, "Generate a fuzzy hash for a file or string.")
   flag.BoolVar(&compare, "compare", true, "Compare two hashes and return the percentage (%) familiarity.")
   flag.StringVar(&file1, "file1", "false", "[Conditional] File or string to generate and/or compare a hash for.")
   flag.StringVar(&file2, "file2", "false", "[Conditional] File or string to generate and/or compare a hash for.")
   flag.StringVar(&hash1, "hash1", "false", "[Conditional] Hash to run a comparison against. The needle.")
   flag.StringVar(&hash2, "hash2", "false", "[Conditional] Hash to compare a hash1 to. The haystack.")
}

func createHash(fp *os.File) {

   if stat, err := fp.Stat(); err != nil {
      log.Fatal(err)   
   }  else {
      hash, err := spamsum.HashReadSeeker(fp, stat.Size())
      if err == nil {
         fmt.Println(hash)
      }
   }
}

func readFile(path string, fi os.FileInfo, err error) error {

   f, err := os.Open(path)
   if err != nil {
      fmt.Fprintln(os.Stderr, "ERROR:", err)
      os.Exit(1)  //should only exit if root is null, consider no-exit
   }

   switch mode := fi.Mode(); {
   case mode.IsRegular():
      createHash(f)
   case mode.IsDir():
      fmt.Fprintln(os.Stderr, "INFO:", fi.Name(), "is a directory.")      
   default: 
      fmt.Fprintln(os.Stderr, "INFO: Something completely different.")
   }
   return nil
}

func main() {

   flag.Parse()

   if flag.NFlag() < 2 {    // can access args w/ len(os.Args[1:]) too
      fmt.Fprintln(os.Stderr, "Usage:  ssdir [-fuzz] [-file1 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  ssdir [-fuzz] [-file1 ...] [-file2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  ssdir [-compare] [-hash1 ...] [-hash2 ...]")
      fmt.Fprintln(os.Stderr, "Output: [CSV] 'file1','hash'")
      fmt.Fprintln(os.Stderr, "Output: [CSV] 'file1 | hash1','file2 | hash2','similarity'")
      flag.Usage()
      os.Exit(0)
   }
   
   if (fuzz == true && file1 != "false") {
      filepath.Walk(file1, readFile)
   }

}

