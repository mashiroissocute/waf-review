Python中的set是一种无序且不重复的集合数据类型，它的主要用法包括：

1. 创建set

创建set有两种方式：

- 使用花括号 `{}` 包含元素来创建set，例如：

```python
my_set = {1, 2, 3}
```

- 使用 `set()` 函数来创建set，例如：

```python
my_set = set([1, 2, 3])
```

注意：空的花括号 `{}` 会创建一个空的字典对象而不是set，因此要创建一个空的set，应该使用 `set()` 函数。

2. 添加和删除元素

可以使用 `add()` 方法向set中添加元素，例如：

```python
my_set.add(4)
```

可以使用 `remove()` 或 `discard()` 方法从set中删除元素，例如：

```python
my_set.remove(3)  # 如果元素不存在会抛出 KeyError 异常
my_set.discard(3)  # 如果元素不存在不会抛出异常
```

3. 集合运算

Python中的set可以进行多种集合运算，包括：

- 并集（Union）：使用 `union()` 方法或 `|` 运算符，例如：

```python
set1 = {1, 2, 3}
set2 = {2, 3, 4}
print(set1.union(set2))  # 输出 {1, 2, 3, 4}
print(set1 | set2)  # 输出 {1, 2, 3, 4}
```

- 交集（Intersection）：使用 `intersection()` 方法或 `&` 运算符，例如：

```python
print(set1.intersection(set2))  # 输出 {2, 3}
print(set1 & set2)  # 输出 {2, 3}
```

- 差集（Difference）：使用 `difference()` 方法或 `-` 运算符，例如：

```python
print(set1.difference(set2))  # 输出 {1}
print(set1 - set2)  # 输出 {1}
```

- 对称差集（Symmetric Difference）：使用 `symmetric_difference()` 方法或 `^` 运算符，例如：

```python
print(set1.symmetric_difference(set2))  # 输出 {1, 4}
print(set1 ^ set2)  # 输出 {1, 4}
```

4. 判断元素是否在set中

可以使用 `in` 或 `not in` 关键字来判断一个元素是否在set中，例如：

```python
print(2 in my_set)  # 输出 True
print(4 not in my_set)  # 输出 True
```

5. 其他常用方法

- `len()` 函数可以获取set中元素的个数，例如：

```python
print(len(my_set))  # 输出 3
```

- `clear()` 方法可以清空set中的所有元素，例如：

```python
my_set.clear()
print(my_set)  # 输出 set()
```

- `copy()` 方法可以复制一个set，例如：

```python
new_set = my_set.copy()
print(new_set)  # 输出 {1, 2, 3}
```

这些是Python中set的基本用法，希望对你有所帮助！