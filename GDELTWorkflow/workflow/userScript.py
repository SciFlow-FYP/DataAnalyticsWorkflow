import os
from collections import OrderedDict
import sys

import pandas as pd

#Workflow

'''
Import df
select columnss
Drop missing value criteria
fill missing values
Convert categoricacl to numerical
feature sclaing - standardize/ normalize
split data - RF
'''
#input location
inputDataset = "/home/amanda/FYP/ds/combined.csv"

#specify output locatiion
outputDataset = "/home/amanda/FYP/testcsv/RFout.csv"

#wait time between each step of workflow (in milli seconds)
waitTime = 10000

#read csv to pandas df
inputDataFrame = pd.read_csv(inputDataset)

#######################		SELECTION	####################################

#select columns
#if "all" select everything. else give a list ["whatever1", "whatever2"]
selectColumns = ["SQLDATE", "Actor1Geo_Type", "Actor1Geo_CountryCode", "QuadClass", "GoldsteinScale", "NumMentions", "AvgTone"] #done
#selectColumns = ["SQLDATE"]



#select rows
selectFromRow = OrderedDict()
selectFromRow['Year'] = ["2018", "2019"] #doesnt work


#######################		CLEANING	####################################

#Run anyway - Drop unique columns #done

#user defined missing values
missingValues = ["n/a", "na", "--"]

#drop columns according to user defined empty value percentage
userDefinedColPercentage = 10 #done

#drop rows according to user defined empty value percentage
userDefinedRowPercentage = 10

#drop duplicate rows - run anyway

#Research how best to fill missing values
#mode for user defined columns
modeColumns = "all" #done
 

#######################		TRANSFORMATION	####################################

#Normalize
userDefinedNormalizeColumns = ["AvgTone"]

userDefinedEncodeColumns = ["Actor1Geo_CountryCode"]




