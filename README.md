##

2018-09-14 日之前完成

## inclus

自动下载inclus.yaml中包含的版本及相关依赖软件

同时生成versions.yaml文件

### 使用

```
export github token

inclus g api v1.0.1 web v2.0 wiki v3.0

result api 1000 web 500 wiki 1100

v100.500.1100

```


### 上传versions.yaml

TOKEN=xxxxxx inclus u

将上传到github, token为github的密钥