# 5th-simple-tiktok
第五届字节青训营项目
# simple-demo

## 抖音项目服务端简单示例

具体功能内容参考飞书说明文档

工程无其他依赖，直接编译运行即可

```shell
go build && ./simple-demo

换：./5th-simple-tiktok
```

### 功能说明

接口功能不完善，仅作为示例

* 用户登录数据保存在内存中，单次运行过程中有效
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可

### 测试

test 目录下为不同场景的功能测试case，可用于验证功能实现正确性

其中 common.go 中的 _serverAddr_ 为服务部署的地址，默认为本机地址，可以根据实际情况修改

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试

### 试图拆分成微服务
用户模块：  
    - 用户注册  
    - 用户登录  
    - 用户信息  

点赞，评论，收藏功能的实现思路：  
https://www.cnblogs.com/xiaoyantongxue/p/15758271.html  
https://blog.csdn.net/zhizhengguan/article/details/87264601  
https://blog.csdn.net/shachao888/article/details/117129285  

Corn定时任务的使用  
https://www.imooc.com/article/46466  
http://www.zyiz.net/tech/detail-141215.html  

Go操作Redis
https://www.cnblogs.com/itbsl/p/14198111.html

Redis批量模糊删除 Key  
https://blog.csdn.net/qianyer/article/details/106383423   

https://blog.csdn.net/qq171563857/article/details/107406409  
在缓存评论内容的时候只缓存不变的内容,比如评论ID,评论时间,评论内容  
点赞数和回复数都另外用Redis计数器处理,读取缓存时同时读取计数器缓存  

        