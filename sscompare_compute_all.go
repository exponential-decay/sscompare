//Functions specifically for the compute all capability of the tool.

package main

import "github.com/dutchcoders/gossdeep"

var rescache []result

type pathhash struct {
   fuzzy string
   path  string
}

//compare every result to each other...
func createcomparisontable(hashes []pathhash, all bool) (error) {
   //header for the CSV generated in compute all function
   outputpathheader()
   //compare n*n values and output
   for _, v1 := range hashes {
      for _, v2 := range hashes {
         if v1.path != v2.path {
            var err error
            var r result
            r.paths = true
            r.s1 = v1.path
            r.s2 = v2.path
            r.score, err = ssdeep.Compare(v1.fuzzy, v2.fuzzy)
            if err != nil {
               return err
            }
            //don't record zeros, n.b. log not recording of zeros
            if r.score != 0 {
               collectresults(r)
            }
         }
      }
   }
   return nil
}

//is one result struct different from the other
func different(newr, oldr result) bool {
   if newr.score == oldr.score {
      if (newr.s1 == oldr.s1 || newr.s1 == oldr.s2) && 
         (newr.s2 == oldr.s1 || newr.s2 == oldr.s2) {
         return false
      }
   }
   return true
}

//do we have a new result to output?
func newresult(newr result) bool {
   if len(rescache) > 0 {
      for x := range(rescache) {
         if !different(newr, rescache[x]) {
            return false
         }
      }   
   } else {
      rescache = append(rescache, newr)
   }
   rescache = append(rescache, newr)
   return true
}

//function to allow us to begin outputting results
func collectresults(r result) {
   if !all {
      //n*n comparison, we will filter out results matching themselves
      if newresult(r) {
         outputresults(r)
      }
   } else {
      outputresults(r)
   }
}

