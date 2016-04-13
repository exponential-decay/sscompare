import os
import unicodecsv

class genericCSVHandler():

   BOM = False
   BOMVAL = '\xEF\xBB\xBF'

   def __init__(self, BOM=False):
      self.BOM = BOM

   def __getCSVheaders__(self, csvcolumnheaders):
      header_list = []
      for header in csvcolumnheaders:      
         header_list.append(header)
      return header_list

   def trimheaders(self, row):
      head = []      
      for h in row:
         head.append(h.strip())
      return head

   # returns list of rows, each row is a dictionary
   # header: value, pair. 
   def csvaslist(self, csvfname):
      columncount = 0
      csvlist = None
      if os.path.isfile(csvfname): 
         csvlist = []
         with open(csvfname, 'rb') as csvfile:
            if self.BOM is True:
               csvfile.seek(len(self.BOMVAL))
            csvreader = unicodecsv.reader(csvfile)
            for row in csvreader:
               if csvreader.line_num == 1:		# not zero-based index
                  header_list = self.__getCSVheaders__(self.trimheaders(row))
                  columncount = len(header_list)
               else:
                  csv_dict = {}
                  #for each column in header
                  #note: don't need ID data. Ignoring multiple ID.
                  for i in range(columncount):
                     csv_dict[header_list[i]] = row[i].strip()
                  csvlist.append(csv_dict)
      return csvlist