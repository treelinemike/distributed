# imports
import urllib.request
from time import time
import numpy as np
import string
import matplotlib.pyplot as plt
import threading
from multiprocessing import Process

source_url = "https://www.gutenberg.org/cache/epub/100/pg100.txt" # complete works
#source_url = "https://www.gutenberg.org/cache/epub/1513/pg1513.txt" # romeo and juliet

letter_counts = np.zeros((26,1))

def p_count_task(task_url,start_line,num_lines_to_read):

    global letter_counts

    print("Starting process")
    
    # open url
    this_file = urllib.request.urlopen(task_url)
    
    # discard all lines ahead of our starting poitn
    for line_idx in range(0,start_line):
        next(this_file)

    # iterate through lines in the desired range
    for line_idx in range(start_line,start_line + num_lines_to_read):
        
        # read line, but fail if we've hit the end of the file iterator
        try:
            #with sourcelock:
            this_line = next(this_file)   # note: uses iterator not readline()
        except StopIteration:
            return
        
        this_line = this_line.decode('utf-8').lower()
        for i in range(0,26):
            this_count = this_line.count(chr(i+97))
            #print(this_count)
            #with listlock:
            letter_counts[i] += this_count  # this is the critical line!


# try this with threads and locks
if __name__ == '__main__':
    NUM_PROCESSES = 1
    
    all_processes = list()
    
    # open file from url and count lines
    rawdata = urllib.request.urlopen(source_url)
    numlines = sum(1 for l in rawdata)
    lines_per_process = np.ceil(numlines/NUM_PROCESSES).astype(int)
    
    # launch threads
    start_time = time()
    for process_num in range(0,NUM_PROCESSES):
        start_line = np.ceil((process_num)*lines_per_process).astype(int)
        end_line = np.ceil(start_line + (lines_per_process - 1)).astype(int)
        print('Process {0:d} -> start: {1:d}, stop: {2:d}'.format(process_num,start_line,end_line))
        p = Process(target=p_count_task,args=(source_url,start_line,lines_per_process,))
        p.start()
        all_processes.append(p)
    
    for p in all_processes:
        p.join()
    
    end_time = time()
    print("Processing time: {:.4f}".format(end_time - start_time))
    print(letter_counts)