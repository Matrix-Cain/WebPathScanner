# pathmap

## 📖Introduction

[ZH-CN/中文](https://github.com/zwh-china/WebPathScanner/blob/main/README_ZH-CN.md)

A simple but effective web path scanner tool, an implementation based on [Dirmap](https://github.com/H4ckForJob/dirmap)

![image-20220504012723565](https://github.com/zwh-china/WebPathScanner/blob/main/README/banner.png)

## ⚙️ About Pathmap

+ Concurrent Support using [ants](https://github.com/panjf2000/ants)

+ CLI Feature using [cobra](https://github.com/spf13/cobra)
+ Human Readable Config File using [TOML](https://github.com/toml-lang/toml)
+ Logging using [Logrus](https://github.com/sirupsen/logrus)

## 🛠 How to use

Just follow the help message

+ Basic scan for an ordinary website

> pathmap.exe scan url http://www.target.com

+ Fuzz mode scan scenario

> pathmap.exe scan url "http://www.target.com/flag.{ext}" --mode 1

![image-20220504030122604](https://github.com/zwh-china/WebPathScanner/blob/main/README/demo1.png)PS. 

--mode 

0 (Default as vintage dict mode) 

1 (Fuzz mode for particular prefix or suffix)  

+ Specify proxy to view outbound traffic

> pathmap.exe scan url "https://target.com/flag.{ext}" --mode 1

+ Using target file for scanning multiple targets

> pathmap.exe scan file target.txt

The target format in target.txt should be like

```
http://target0.com
http://target1.com
http://target2.com
http://target3.com
```

## 🚀 Features

+ Faster than original Dirmap
+ Custom Dictionary supports
+ Auto Handle fake 404
+ Random Sleep support
+ Custom User-Agent supported

## 🚧TODO

+ Save scanning results to file

Since this is a small tool I write just for practice Golang. There is surely many bugs to fix. And the scanner's functions are not rich. So any suggestions and bug reports are more than welcome. Fell free to PR if you have some good idea.

## ⚔️Disclaimer

This software has been created purely for the purposes of academic research and for the development of effective defensive techniques, and is not intended to be used to attack systems except where explicitly authorized. Project maintainers are not responsible or liable for misuse of the software. Use responsibly.