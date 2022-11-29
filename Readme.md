# jsonpath 解析器

最近要用到JSON PATH 解析器， 那就写一个吧， 

未完待续， 突然发现json path有好多符号啊。竟然有人跟我说json path不好用， 科科。

## quick start

```go
func main() {
	origin := `{"data": {"text": "人生短短几个秋", "type": "彩虹屁"}}`
	var data any
	err := json.Unmarshal([]byte(origin), &data)
	if err != nil {
		log.Fatal(err)
	}
	val, err := core.Read(data, "$.data.text")
	log.Println(val, err)
}
```

## todo

- [x] 支持 $
- [x] 支持 .
- [x] 支持 ['key']
- [] 支持各种奇怪的符号

## jsonpath规则

- 操作符

| 符号 | 	描述
| :-: | :-
|$ |	查询的根节点对象，用于表示一个json数据，可以是数组或对象
|@ |	过滤器（filter predicate）处理的当前节点对象
|* |	获取所有节点
|. |	获取子节点
|.. |	递归搜索，筛选所有符合条件的节点
|?() |	过滤器表达式，筛选操作
|[start:end] |	数组片段，区间为[start,end),不包含end
|[A]或[A,B] |	迭代器下标，表示一个或多个数组下标

- 函数。可以在JsonPath表达式执行后进行调用，其输入值为表达式的结果。

| 名称 |	描述
| :-: | :-
|min() |	获取数值类型数组的最小值
|max() |	获取数值类型数组的最大值
|length() |	获取数值类型数组的长度，例如$.data.length()
|... |	...


- 过滤器。过滤器是用于过滤数组的逻辑表达式。

| 操作符 |	描述
| :-: | :-
|== |	等于
|!= |	不等于
|< |	小于
|=~ |	判断是否符合正则表达式，例如[?(@.name =~ /foo.*?/i)]
|in |	所属符号，例如[?(@.type in ["小雨","中到大雨"])]
|nin |	排除符号
|... |	....