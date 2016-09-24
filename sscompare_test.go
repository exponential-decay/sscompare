package main

import (
   "testing"
)

var goodhashcompare = map[[2]string]int {
   [2]string{"3:hZIlbWKEaFgcR2LBOmuF3E1MC0LGFhL7muUSsehn:hZIlbW5aJR2cmu+1MzyD7muUzehn", "3:hZIlbWKEaFgcR2LBOmuF3E1MC0LGFhL7muUSsehn:hZIlbW5aJR2cmu+1MzyD7muUzehn"}: 100,
}

var badhashcompare = map[[2]string]error {
   [2]string{"abc", "abc"}: sscomp_err,   
}

var ExportCompareHashes = comparehashes

func TestCompareGoodHashes(t *testing.T) {
   for k, v := range(goodhashcompare) {
      r, err := ExportCompareHashes(k[0], k[1])
      if err != nil {
         t.Errorf("Unexpected error from comparison: %v\n", err)
         break
      }
      if r.score != v {
         t.Errorf("Expected score %d doesn't match actual score %d\n", v, r.score)         
      } 
   }   
}

func TestCompareBadHashes(t *testing.T) {
   for k, v := range(badhashcompare) {
      _, err := ExportCompareHashes(k[0], k[1])
      if err != v {
         t.Errorf("Unexpected error from comparison:\nReceived: %v\nExpected: %v\n", err, v)
         break
      }
   }   
}
