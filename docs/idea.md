# FOFA_PK
我想要开发一个程序，完成FOFA规则的PK的游戏，判断两个参赛者提供的查询语法谁的更加全面且准确，最终得分高的参赛者赢。

## FOFA简介
FOFA是一个网络空间测绘平台，它的核心功能就是根据用户输入的查询语法进行网络空间的检索。查询语法有很多个关键字，比如ip、host、title、body、domain、cert、banner、body等等。针对不同的任务输入的查询语法不通，返回的结果就有好有坏。

## FOFA的API与命令行介绍
FOFA提供了API，通过`github.com/LubyRuffy/gofofa`作为库来调用：
```golang
cli, err := gofofa.NewClient()
cli.HostSearch("domain=fofa.info", 10, []string{"host", "ip", "domain", "certs_subject_org"})
```
返回
```text
https://g.fofa.info,106.75.10.35,fofa.info,
https://jw.fofa.info,106.75.10.35,fofa.info,
jw.fofa.info,106.75.10.35,fofa.info,
s.fofa.info:9093,120.52.152.23,fofa.info,
https://www.fofa.info,113.96.142.1,fofa.info,
www.fofa.info,113.96.142.1,fofa.info,
s.fofa.info:2181,120.52.152.23,fofa.info,
https://fofa.info,101.36.123.64,fofa.info,
https://en.fofa.info,101.36.123.64,fofa.info,
https://fofa.info,101.36.123.64,fofa.info,
https://staticbeta.fofa.info,106.75.56.231,fofa.info,
beta.fofa.info,106.75.56.231,fofa.info,
https://beta.fofa.info,106.75.56.231,fofa.info,
enbeta.fofa.info,106.75.56.231,fofa.info,
apibeta.fofa.info,106.75.56.231,fofa.info,
https://api.gamma.fofa.info,117.50.16.112,fofa.info,
gamma.fofa.info,117.50.16.112,fofa.info,
https://gamma.fofa.info,117.50.16.112,fofa.info,
api.gamma.fofa.info,117.50.16.112,fofa.info,
static.fofa.info,104.21.26.154,fofa.info,
https://static.fofa.info,172.67.137.15,fofa.info,"Cloudflare, Inc."
https://staticbeta.fofa.info,106.75.56.231,fofa.info,
fofa.info,106.75.95.206,fofa.info,
hub.fofa.info,43.193.0.238,fofa.info,
https://hub.fofa.info,43.193.0.238,fofa.info,
fofa.info,106.75.95.206,fofa.info,
https://store.fofa.info,106.75.10.17,fofa.info,
store.fofa.info,106.75.10.17,fofa.info,
https://xgbeta.fofa.info,152.32.214.56,fofa.info,
https://staticbeta1.fofa.info,106.75.56.231,fofa.info,
https://apibeta1.fofa.info,106.75.56.231,fofa.info,
https://static.beta.fofa.info,106.75.56.231,fofa.info,
sywx7xh5at8jwbjd2314.fofa.info:22703,123.58.250.22,fofa.info,
octra.fofa.info,106.75.96.150,fofa.info,
https://octra.fofa.info,106.75.96.150,fofa.info,
```

## 程序的整体流程
- 两个参赛者分别写出FOFA的查询语句，提供给系统；
- 系统调用FOFA的API接口，返回所有的结果；
- 根据两个结果集和，提取差异性，比如domain的差异，ip的差异；
- 每个差异询问用户判断是否有效，可以选择如下三个选择之一：有效、无效、未知。对于有效就加分，无效就减分，未知就不影响任何参赛者的评分。
- 遍历完成后，计算谁的分值更高；

## 代码输出
### 前端代码
提供一个页面，要求如下：
- 分为三个功能区：
    - 最上面是两个文件上传的区块
    - 中间是对比和评分的区块
    - 最下面是日志区块
- 最上面的区块有一个按钮“Load Diff”，点击后调用后端的`/api/v1/add_diff_task`接口，得到taskid
    - 使用taskid和ts循环请求`/api/v1/task_info`接口，一直到data.finished为true
        - 把过程的日志实时的打印出来到下面的日志区块
    - 显示q1多出来的数据xx条，q2多出来的数据xx条，相同的数据xx条
- 中间的对比和评分的区块在diff完成后，显示差异的数据，但是不显示是哪个参赛者的
    - 每次随机从a1和a2中提取未处理的一条差异项
    - 提供三个评价按钮，分别是：有效、无效、未知
    - 用户点击后，更新参赛者的计分，加载下一个差异项
    - 一直到结束

### 后端代码
提供几个API，要求如下：
- 提供/api/v1/query  GET 用于测试单个查询的有效性
    - 参数
        - q 代表fofa的原始查询语句
    - 返回json格式
        - code int 状态码
        - error string 错误提示
        - data [][]string 结果

- 提供/api/v1/add_diff_task  POST 用户返回两个查询语句的不同
    - 参数
        - q1 代表第一个参赛者提供的fofa的原始查询语句
        - q2 代表第二个参赛者提供的fofa的原始查询语句
    - 返回json格式
        - code int 状态码
        - error string 错误提示
        - data {"taskid": 1234} 结果 为创建的异步任务的taskid

- 提供/api/v1/task_info  GET 用户返回两个查询语句的不同，注意每个参赛者的语句都只提取三个月之内的数据（通过控制fofa的after参数来实现）
    - 参数
        - taskid int 任务id
        - ts string 时间戳， 默认为空，未完成的情况下取进度结果的last_ts参数作为请求参数
    - 返回json格式
        - code int 状态码
        - error string 错误提示
        - data
            - 如果没有完成，返回进度：{"finished":false,"progress":50.00,"logs":["正在查询参赛者1的FOFA结果。。。","正在查询参赛者2的FOFA结果。。。"],"last_ts":"最后更新的时间戳"}
            - 如果已经完成，返回结果： {"finished":true,"progress":100.00,"logs":["完成所有操作"],"last_ts":"最后更新的时间戳","a1"{"domain":["a.com"],"ip":["1.1.1.1"]}, "a2":{"domain":["c.com"],"ip":["2.2.2.2"]}} 结果，a1代表q1多出来的结果，a2代表q2多出来的结果



====================
请完成上面的代码输出，包含完整的前端代码（vue3）和后端代码（golang）。