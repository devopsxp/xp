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