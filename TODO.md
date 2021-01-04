[docker化pipeline]

1. -w /tmp/workspace #重新制定工作目录
2. --volumes-from a713nlasd #共享运行中的docker app
3. 构建缓存 使用-v m2:/root/.m2进行

[git流]

创建git_role,确定项目名 传递文件夹路径

[global]

设置宿主机目录 

[k8s调度]

1. 使用job或pod，利用contained共享磁盘空间的功能进行代码共享
2. 可用利用pvc进行/root/.m2 /root/.sonar等构建缓存，进行加速构建

[中间构建物]

1. 保存
2. 记录
3. 下载

[C/S模式]

1. RPC Agent改造
2. 日志Call_Back机制
    1）logs call back to server
    2) logs call back to es
    3) logs call back to kafka
    4) logs call back to logstash

[frontend]

1. pipeline yaml记录
2. 日志展示
3. 实时日志+pipeline显示（drone）
4. pipeline执行记录展示

[功能加强]
1.微服务化开发
2.多云环境支持：aws、asure、gcp等
3.多种中间件支持，融合jenkins pipeline触发或者触发jenkins pipeline
4.支持定时任务
5.支持并发任务
6.支持暂停等待项目
7.支持DAG流向图编排
8.支持内置变量
9.支持结果状态接口返回+redis记录
10.支持流水线度量
11.支持代码扫描质量度量和统计