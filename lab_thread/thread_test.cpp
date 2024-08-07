// g++ -Wall -no-pie -pthread -lpthread -o thread_test thread_test.cpp

#include <thread>
#include <iostream>
#include <vector>

#define NUM_THREADS 10

// global for keeping count
volatile int32_t count = 0;

// function to be run in threads
void run_count(void){

    // increment and decrement count
    // should have NO NET EFFECT
    for(uint32_t i=0;i<50000; ++i){
        ++count; 
        --count;
    }
    return;
}

// main
int main(void){

    std::vector<std::thread> thread_vec;

    // launch some threads and store in vector
    for(unsigned int i = 0; i < NUM_THREADS; ++i){
        std::thread thd(run_count);
        thread_vec.push_back(move(thd));
    }
    
    // join threads
    for(unsigned int i = 0; i < NUM_THREADS; ++i){
        thread_vec[i].join();
    }

    // display result
    std::cout << "Count = " << count << std::endl;

    return 0;
}
