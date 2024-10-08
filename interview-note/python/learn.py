import functools


def repeat(num):
    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args,**kwargs):
            for i in range(0,num):
                func(*args,**kwargs)
        return wrapper
    return decorator


a =10
@repeat(a)
def myprint(msg):
    print(msg)
    

myprint("sss")