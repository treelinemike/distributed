# imports
import urllib.request
from time import time
import numpy as np
import string
import matplotlib.pyplot as plt
import threading

source_url = "https://www.gutenberg.org/cache/epub/100/pg100.txt" # complete works
#source_url = "https://www.gutenberg.org/cache/epub/1513/pg1513.txt" # romeo and juliet

a = 0
while(1):
    if(a%1):
        a = a + 1
    else:
        a = a - 1
        
def count_task():
    while(1):
        try:
            this_line = next(rawdata)   # note: uses iterator not readline()
        except StopIteration:
            return
        this_line = this_line.decode('utf-8').lower()
        for i in range(0,26):
            letter_counts[i] += this_line.count(chr(i+97))

# try this sequentially
rawdata = urllib.request.urlopen(source_url)
letter_counts = np.zeros((26,1))
start_time = time()
count_task()
end_time = time()
print("Processing time: {:.4f}".format(end_time - start_time))
print(letter_counts)