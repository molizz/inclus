## inclus

自动下载inclus.yaml中包含的版本及相关依赖软件

## 使用

### 下载

下载 v2.0 版本的所有依赖(默认从当前命令程序路径读取 inclus.yaml 文件):
```
inclus version v2.0
```

下载 v2.0 版本的所有依赖, 指定配置文件 inclus.yaml:
```
inclus version v2.0 inclus.yaml
```


### 提交inclus.yaml

将inclus.yaml提交到仓库中.

需要在 https://github.com/settings/tokens/new 新增 access token. 至少需要 repo:repo_deployment 权限

配置好后, 通过命令

```
inclus commit
```

将inclus.yaml提交到仓库

或指定配置文件

```
inclus commit inclus.yaml
```