//Functions specifically for the compute all capability of the tool.

package main

import (
   "os"
   "fmt"
   "time"
   "github.com/dutchcoders/gossdeep"
)

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
