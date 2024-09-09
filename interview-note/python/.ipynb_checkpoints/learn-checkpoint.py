from queue import Queue

q = Queue()
q.put(1)  # 向队列中添加元素
q.put(2)  # 向队列中添加元素

item = q.get()  # 从队列中获取元素
print(item)