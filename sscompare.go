package main

import (
	"os"
	"fmt"
   "flag"
   "path/filepath"
   "github.com/dutchcoders/gossdeep" 
)

var fuzz bool
var compare bool
var file1 string
var file2 string
var hash1 string
var hash2 string
var string1 string
var string2 string

func init() {
   flag.BoolVar(&fuzz, "fuzz", false, "Generate a fuzzy hash for a file or string.")
   flag.BoolVar(&compare, "compare", false, "Compare two hashes and return the percentage (%) familiarity.")
   flag.StringVar(&file1, "file1", "false", "[Conditional] File or string to generate and/or compare a hash for.")
   flag.StringVar(&file2, "file2", "false", "[Conditional] File or string to generate and/or compare a hash for.")
   flag.StringVar(&string1, "string1", "false", "[Conditional] File or string to generate and/or compare a hash for.")
   flag.StringVar(&string2, "string2", "false", "[Conditional] File or string to generate and/or compare a hash for.")
   flag.StringVar(&hash1, "hash1", "false", "[Conditional] Hash to run a comparison against. The needle.")
   flag.StringVar(&hash2, "hash2", "false", "[Conditional] Hash to compare a hash1 to. The haystack.")
}

func compareHashes(hash1 string, hash2 string) {
   score, err := ssdeep.Compare(hash1, hash2)
   if err != nil {
      fmt.Fprintln(os.Stderr, "ERROR:", err)
      os.Exit(1)      
   }
   fmt.Fprintln(os.Stdout, hash1, ",", hash2, ",", score)
}

func fileExists(path string) bool {
   var exists bool = false   
   _, err := os.Open(path)
   if err != nil {
      fmt.Fprintln(os.Stderr, "ERROR:", err)
   } else {
      exists = true
   }
   return exists
}

func hashString(str string) string {
   hash, err := ssdeep.HashString(str)
   if err != nil {
      fmt.Fprintln(os.Stderr, "ERROR:", err)
      os.Exit(1)
   }
   return hash
}

func compareStrings(str1 string, str2 string) {
   hash1 := hashString(str1)
   hash2 := hashString(str2)
   score, err := ssdeep.Compare(hash1, hash2)
   if err != nil {
      fmt.Fprintln(os.Stderr, "ERROR:", err)
      os.Exit(1)
   }
   fmt.Fprintln(os.Stdout, str1, ",", str2, ",", score)
}

func createFileHash(path string) string {

   var hash string
   f, err := os.Open(path)
   if err != nil {
      fmt.Fprintln(os.Stderr, "ERROR:", err)
   } else {
      f.Close()
      hash, err = ssdeep.HashFilename(path)   //confusing function title
      if err != nil {
         fmt.Fprintln(os.Stderr, "ERROR:", err)
         os.Exit(1)
      }
   }
   return hash
}

func compareFiles(file1 string, file2 string) {
   if fileExists(file1) && fileExists(file2) {
      hash1 := createFileHash(file1)
      hash2 := createFileHash(file2)
      score, err := ssdeep.Compare(hash1, hash2)
      if err != nil {
         fmt.Fprintln(os.Stderr, "ERROR:", err)
         os.Exit(1)
      }
      fmt.Fprintln(os.Stdout, file1, ",", file2, ",", score)
   }
}

func handleHash(hash string) {
   fmt.Println(hash)
}

func readFile(path string, fi os.FileInfo, err error) error {

   if fileExists(path) {
      switch mode := fi.Mode(); {
      case mode.IsRegular():
         fmt.Println(createFileHash(path))
      case mode.IsDir():
         fmt.Fprintln(os.Stderr, "INFO:", fi.Name(), "is a directory.")      
      default: 
         fmt.Fprintln(os.Stderr, "INFO: Something completely different.")
      }
   }
   return nil
}

func main() {

   flag.Parse()

   if flag.NFlag() < 2 {    // can access args w/ len(os.Args[1:]) too
      fmt.Fprintln(os.Stderr, "Usage:  ssdir [-fuzz] [-file1 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  ssdir [-fuzz] [-string1 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  ssdir [-compare] [-file1 ...] [-file2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  ssdir [-compare] [-string1 ...] [-string2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  ssdir [-compare] [-hash1 ...] [-hash2 ...]")
      fmt.Fprintln(os.Stderr, "Output: [CSV] 'file1','hash'")
      fmt.Fprintln(os.Stderr, "Output: [CSV] 'file1 | hash1','file2 | hash2','similarity'")
      flag.Usage()
      os.Exit(0)
   }
   
   if (fuzz == true && file1 != "false") {
      filepath.Walk(file1, readFile)
   }

   if (fuzz == true && string1 != "false") {
      fmt.Println(hashString(string1))
   }

   if (compare == true && file1 != "false" && file2 != "false") {
      compareFiles(file1, file2)      
   }
   
   if (compare == true && string1 != "false" && string2 != "false") {
      compareStrings(string1, string2)
   }

   if (compare == true && hash1 != "false" && hash2 != "false") {
      compareHashes(hash1, hash2)
   } 
}

