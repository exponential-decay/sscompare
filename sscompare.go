package main

import (
	"os"
	"fmt"
   "flag"
   "time"
   "path/filepath"
)

var (
   fuzz, compare, compute, all, storeHashes bool
   file1, file2, hash1, hash2, string1, string2, dir string
   results_cache, hashes [][]string
   start time.Time
)

func init() {
   flag.BoolVar(&fuzz, "fuzz", false, "Generate a fuzzy hash for a file or string.")
   flag.BoolVar(&compare, "compare", false, "Compare two hashes and return the percentage (%) familiarity.")
   flag.BoolVar(&compute, "compute", false, "Compare all file hashes and output all comparisons for a given directory.")
   flag.StringVar(&file1, "file1", "false", "[Conditional] File or string to generate and/or compare a hash for.")
   flag.StringVar(&file2, "file2", "false", "[Conditional] File to compare file1 to.")
   flag.StringVar(&string1, "string1", "false", "[Conditional] File or string to generate and/or compare a hash for.")
   flag.StringVar(&string2, "string2", "false", "[Conditional] String to compare string1 to.")
   flag.StringVar(&hash1, "hash1", "false", "[Conditional] Hash to run a comparison against. The needle.")
   flag.StringVar(&hash2, "hash2", "false", "[Conditional] Hash to compare a hash1 to. The haystack.")
   flag.StringVar(&dir, "dir", "false", "[Conditional] Directory to run a full 1:1 comparison against.")
   flag.BoolVar(&all, "all", false, "[Optional] Output all files, including zero matches and duplicates.")
}

func fileExists(path string) (bool, os.FileInfo) {
   f, err := os.Open(path)
   defer f.Close()
   if err != nil {
      fmt.Fprintln(os.Stderr, "Error:", err)
      os.Exit(1)
   } 
   fi, err := f.Stat()
   if err != nil {
      fmt.Fprintln(os.Stderr, "Error:", err)
      os.Exit(1)
   }
   return true, fi
}

func computeAll(path string, all bool) { 
   f1, fi := fileExists(path)
   if f1 {
      switch mode := fi.Mode(); {
      case mode.IsDir():
         storeHashes = true
         hashes = make([][]string, 512)
         filepath.Walk(path, readFile)
      }
   }

   if len(hashes) > 0 {
      generateComparisonTable(hashes, all)
   }
}

func readFile(path string, fi os.FileInfo, err error) error {
   f1, _ := fileExists(path) 
   if f1 {
      switch mode := fi.Mode(); {
      case mode.IsRegular():
         hash := createFileHash(path)
         if storeHashes == true {
            row := []string{hash, path}
            hashes = append(hashes, row)
         } else {
            fmt.Fprintf(os.Stderr, "%s,%s\n", path, hash)
         }
      case mode.IsDir():
         //fmt.Fprintln(os.Stderr, "INFO:", fi.Name(), "is a directory.")      
      default: 
         //fmt.Fprintln(os.Stderr, "INFO: Something completely different.")
      }
   }
   return nil
}

func main() {

   flag.Parse()

   if flag.NFlag() < 2 {    // can access args w/ len(os.Args[1:]) too
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-fuzz] [-file1 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-fuzz] [-string1 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compare] [-file1 ...] [-file2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compare] [-string1 ...] [-string2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compare] [-hash1 ...] [-hash2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compute] [-dir ...] [OPTIONAL] [-all]")
      fmt.Fprintln(os.Stderr, "Output: [CSV] 'file1','hash'")
      fmt.Fprintln(os.Stderr, "Output: [CSV] 'similarity','file1 | string1 | hash1','file2 | string1 | hash2',")
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

   if (compute == true && dir != "false") {
      start = time.Now()
      computeAll(dir, all)
   }
}

