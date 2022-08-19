# pathmap

## 📖介绍

[英语/EN](https://github.com/zwh-china/WebPathScanner/blob/main/README_ZH-CN.md)

一个简单高效的web路径扫描工具, 基于 [Dirmap](https://github.com/H4ckForJob/dirmap)

![image-20220504012723565](https://github.com/zwh-china/WebPathScanner/blob/main/README/banner.png)

## ⚙️ 关于pathmap

+ 使用 [ants](https://github.com/panjf2000/ants)作并发支持

+ 使用 [cobra](https://github.com/spf13/cobra)作为CLI命令框架
+ 采用人类可读性比较好的[TOML](https://github.com/toml-lang/toml)作为config文件格式
+ 使用 [Logrus](https://github.com/sirupsen/logrus)作为log模块

## 🛠 使用方法

使用参考help命令的输出信息即可

+ 基础网站路径扫描

> pathmap.exe scan url http://www.target.com

+ Fuzz模式场景下扫描

> pathmap.exe scan url "http://www.target.com/flag.{ext}" --mode 1

![image-20220504030122604](https://github.com/zwh-china/WebPathScanner/blob/main/README/demo1.png)PS. 

--mode 

0 (默认值，代表使用传统的vintage字典扫描) 

1 (fuzz模式用于对文件前缀或者后缀名进行爆破)  

+ 指定proxy代理来对扫描的url进行监测

> pathmap.exe scan url "https://target.com/flag.{ext}" --mode 1

+ 使用文件的形式指定多个目标进行扫描

> pathmap.exe scan file target.txt

扫描目标的文件的格式应该如下所示:

```
http://target0.com
http://target1.com
http://target2.com
http://target3.com
```

## 🚀 特性

+ 比原版Dirmap更快效率更高
+ 支持使用自定义字典
+ 自动处理目标网站的假性404界面
+ 支持随机请求延时
+ 支持自定义UA头
+ 保存扫描结果为文件

## 🚧TODO

这个项目只是我为了学习和实践Golang写的，因此肯定会有一些bug， 同时作为一个扫描器它的功能还不完善。欢迎通过PR和提issue对本项目做出贡献。

## ⚔️声明

软件只作为开发和学习使用，任何未授权的扫描和攻击行为与开发者无关。