# todo

## V2.0更新

V2.0将成为第一个能够正式使用的版本.将支持以下功能:

1. 用户的登录注册
2. 用户每日任务的邮件推送
3. 用户能够添加一次性任务,定期任务(每日,每周,每月任务)
4. 用户能够给每个任务添加

## 目录结构:

```
-- api     // API服务
 |- controllers // api的controllers
 |- models      // 数据模型
 |- libs        // common libs
 |- routes      // API路由表
 |- config      // API配置服务
-- front   // 前端程序
 |- controllers // 前端控制器
 |- models      // 请求API的model
 |- dist        // 生成的文件,静态文件存放处
   |- js        // JS文件
     |- app.min.js  // app运行js 
     |- vendor.min.js // app运行的第三方库
   |- css       // css存放路径
   |- fonts     // 字体文件存放路径
   |- imgs      // 图片文件存放路径
-- tools
 |- tables.sql // app需要的mysql配置
 |- cmd        // 一些命令行脚本文件  
```
