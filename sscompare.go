package main

import (
	"os"
	"fmt"
   "flag"
   "time"
   "path/filepath"
   "github.com/dutchcoders/gossdeep" 
)

var fuzz, compare, compute, all bool
var file1, file2, hash1, hash2, string1, string2, dir1 string

var storeHashes bool = false
var hashes [][]string
var results_cache [][]string

var start time.Time

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
   flag.StringVar(&dir1, "dir1", "false", "[Conditional] Directory to run a full 1:1 comparison against.")
   flag.BoolVar(&all, "all", false, "[Optional] Output all files, including zero matches and duplicates.")
}

func compareHashes(hash1 string, hash2 string) {
   score, err := ssdeep.Compare(hash1, hash2)
   if err != nil {
      fmt.Fprintln(os.Stderr, "ERROR:", err)
      os.Exit(1)      
   }
   fmt.Fprintf(os.Stdout, "%d,%s,%s\n", score, hash1, hash2)
}

func fileExists(path string) (bool, os.FileInfo) {
   var fi os.FileInfo
   var exists bool = false   
   f, err := os.Open(path)
   if err != nil {
      fmt.Fprintln(os.Stderr, "ERROR:", err)
   } else {
      defer f.Close()
      fi, err = f.Stat()
      if err != nil {
         fmt.Fprintln(os.Stderr, "ERROR:", err)
         os.Exit(1)
      }
      exists = true
   }
   return exists, fi
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
   fmt.Fprintf(os.Stdout, "%d,%s,%s\n", score, str1, str2)
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
   f1, _ := fileExists(file1)
   f2, _ := fileExists(file2)
   if f1 && f2 {
      hash1 := createFileHash(file1)
      hash2 := createFileHash(file2)
      score, err := ssdeep.Compare(hash1, hash2)
      if err != nil {
         fmt.Fprintln(os.Stderr, "ERROR:", err)
         os.Exit(1)
      }
      fmt.Fprintf(os.Stdout, "%d,%s,%s\n", score, file1, file2)
   }
}

var rescache bool = false

func newVariant(str1 string, str2 string, results_cache [][]string) bool {
   add := false   
   exist1 := false
   exist2 := false
   for idx, _ := range results_cache {
      if len(results_cache[idx]) != 0 {
         file1 = results_cache[idx][0]
         file2 = results_cache[idx][1]
         if str1 == file1 || str1 == file2 {
            exist1 = true
         }
         if str2 == file1 || str2 == file2 {
            exist2 = true
            break
         }
      }
   }
   if exist1 != true || exist2 != true {   
      add = true  
   }
   return add
}

func handleComputeResults(score int, hash1 []string, hash2 []string, all bool) int {
   
   added := 0
   hfile1 := hash1[1]
   hfile2 := hash2[1]

   if all != true {
      row := []string{hfile1, hfile2}
      if rescache == false {
         results_cache = make([][]string, 512)
         results_cache = append(results_cache, row)
         fmt.Fprintf(os.Stdout, "%d,%s,%s\n", score, hfile1, hfile2)
         added = 1
         rescache = true
      } else {

         if newVariant(hfile1, hfile2, results_cache) {
            results_cache = append(results_cache, row)
            fmt.Fprintf(os.Stdout, "%d,%s,%s\n", score, hfile1, hfile2)
            added = 1
         }
      }
   } else {
      fmt.Fprintf(os.Stdout, "%d,%s,%s\n", score, hfile1, hfile2)
      added = 1
   }
   return added
}

func generateComparisonTable(hashes [][]string, all bool) {
   fmt.Fprintln(os.Stdout, "score,source,target")   
   total := 0
   x := len(hashes)
   for hash, _ := range hashes {
      x = x - 1
      if len(hashes[hash]) > 1 {    //we have (x) block slice, we may have empties
         hash1 := hashes[hash][0]
         found := false
         for h, _ := range hashes {
            if len(hashes[h]) > 1 {    //preferable to delete empty slices? 
               hash2 = hashes[h][0]
               score, err := ssdeep.Compare(hash1, hash2)
               if err == nil {
                  if score == 100 && found == false { //ignore first identical (itself)
                     found = true
                  } else {
                     if score != 0 {
                        total += handleComputeResults(score, hashes[hash], hashes[h], all)
                     }  //N.B. can opt to log zeroes, if we really care
                  }
               }         
            }
         }
      }
   }
   elapsed := time.Since(start)
   fmt.Fprintln(os.Stderr, total, elapsed)
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
         fmt.Fprintln(os.Stderr, "INFO: Something completely different.")
      }
   }
   return nil
}

func computeAll(path string, all bool) { 
   f1, fi := fileExists(path)
   if f1 {
      mode := fi.Mode()
      if mode.IsDir() {
         storeHashes = true
         hashes = make([][]string, 512)
         filepath.Walk(path, readFile)
      }
      if len(hashes) > 0 {
         generateComparisonTable(hashes, all)
      }
   }
}

func main() {

   flag.Parse()

   if flag.NFlag() < 2 {    // can access args w/ len(os.Args[1:]) too
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-fuzz] [-file1 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-fuzz] [-string1 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compare] [-file1 ...] [-file2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compare] [-string1 ...] [-string2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compare] [-hash1 ...] [-hash2 ...]")
      fmt.Fprintln(os.Stderr, "Usage:  sscompare [-compute] [-dir1 ...] [OPTIONAL] [-all]")
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

   if (compute == true && dir1 != "false") {
      start = time.Now()
      computeAll(dir1, all)
   }
}

