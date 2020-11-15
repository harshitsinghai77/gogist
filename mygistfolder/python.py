def myFunc(x):
    return x*2

lst = [1,2,3,4,5,6,10]

my_array = [myFunc(i) for i in lst]
print(my_array)
