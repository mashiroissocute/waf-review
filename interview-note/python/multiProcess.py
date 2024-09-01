from threading import Lock
from concurrent.futures import ThreadPoolExecutor,as_completed,ProcessPoolExecutor
from multiprocessing import Queue
from functools import reduce
import time


class bank:
    def __init__(self):
        self.balance = 0
        self.lock = Lock()
    def add(self,num):
        with self.lock:
            self.balance += num

def multi_thread():
    b = bank()
    with ThreadPoolExecutor(max_workers=16) as pool:
        for i in range(0,1000):
            pool.submit(b.add,1)
    print(b.balance)
    

class is_prime:
    def __init__(self,num,queue:Queue):
        self.num = num
        self.queue = queue
    
    def check_prime(self,start,end):
        print(f'check {start}-{end}')
        for i in range(start,end):
            if self.num % i == 0:
                return 1
        return 0
    
    def multi_thread(self):
        with ThreadPoolExecutor(max_workers=16) as pool:
            
            start,end = 2,int(self.num**0.5)
            period = int((end - start) / 16)
            
            future_list = []
            
            while start < end:
                future = pool.submit(self.check_prime,start,start+period)
                future_list.append(future)
                start = start + period
                
                
            if reduce(lambda x,y:x+y , [future.result() for future in as_completed(future_list)]) == 0 :
                print(f'{self.num} is prime')
                
    def multi_process(self):
        with ProcessPoolExecutor(max_workers=2) as pool:
            start,end = 2,int(self.num**0.5)
            period = int((end - start) / 16)
            
            future_list = []
            
            while start < end:
                future = pool.submit(self.check_prime,start,start+period)
                future_list.append(future)
                start = start + period
                
                
            if reduce(lambda x,y:x+y , [future.result() for future in as_completed(future_list)]) == 0 :
                print(f'{self.num} is prime')

    
    def multi_process_queue(self):
        
        def num_producer():
            for i in range(2,self.num**0.5+1):
                self.queue.put(i)
            return
            
        def num_consumer():
            while True:
                i = self.queue.get()
                if self.num % i == 0:
                    return 1
                if i == self.num**0.5:
                    self.queue.put(i)
                    return 0
                
        with ProcessPoolExecutor() as pool:
            producr_future = pool.submit(num_producer)
            consumer_future_list = [pool.submit(num_consumer) for i in range(0,4)]
            
        
            if reduce(lambda x,y:x+y , [future.result() for future in as_completed(consumer_future_list)]) == 0 :
                print(f'{self.num} is prime')


        
            
            
if __name__ == '__main__':
    
    isP = is_prime(1000000,Queue())
    start = time.time()
    # isP.multi_thread()
    end = time.time()
    print(f'multi thread cost {end-start}')
    start =time.time()
    # isP.multi_process()
    end = time.time()
    print(f'multi process cost {end - start}')
    
    isP.multi_process_queue()
    