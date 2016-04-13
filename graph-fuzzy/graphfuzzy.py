# -*- coding: utf-8 -*-

import os
import sys
import argparse
import networkx as nx
from CSVHandlerClass import genericCSVHandler

class GraphMe:

   def __init__(self):
      self.csv = genericCSVHandler()
      self.G = nx.Graph()

   def graphMLfromCSV(self, csvfile, filter):     
   
      filter_score = float(filter)
      if filter_score is False:
         filter_score = float(0.0)
         
      csv = self.csv.csvaslist(csvfile)      
      graphfile = csvfile.split('.', 1)[0] + "-graph.xml"      
      nodes = []      
      
      for line in csv:
         #we should have a correctly formatted CSV to work with
         if 'score' in line and 'source' in line and 'target' in line:
            source = line[u'source']
            target = line[u'target']
            score = float(line[u'score']) / 100
            if score >= filter_score:
               if source not in nodes:
                  self.G.add_node(source)
                  nodes.append(source)
               if target not in nodes:
                  self.G.add_node(target)
                  nodes.append(target)
               self.G.add_edge(source,target,weight=score)
               print score
      nx.write_graphml(self.G, graphfile)

def main():

   #	Usage: 	--csv [fuzzy report]
   #	Handle command line arguments for the script
   parser = argparse.ArgumentParser(description='Convert results of a fuzzy hash computation to a network graph, GraphML.\nOutputs using CSV filename with XML suffix.')
   parser.add_argument('--csv', '--results', help='CSV export from sscompare tool.')
   parser.add_argument('--score', help="Filter out results less than this match score.", default=False)
   
   if len(sys.argv)==1:
      parser.print_help()
      sys.exit(1)

   #	Parse arguments into namespace object to reference later in the script
   global args
   args = parser.parse_args()
   
   if args.csv:
      graph = GraphMe()
      graph.graphMLfromCSV(args.csv, args.score)
   else:
      sys.exit(1)

if __name__ == "__main__":
   main()
