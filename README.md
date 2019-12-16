# Seaman

Seaman是Kubecloud的一个Gitops工具，用来同步git配置仓库与k8s集群的状态。

## 部署Seaman

### 依赖条件

1. Git配置仓库,例如gitlab，github等
2. kubecloud 0.1.0版本


### 使用kubernets部署

    $ kubectl apply -f deploy/seaman.yaml


### 使用Docker镜像部署

你可以使用以下命令启动seaman:
    
    $ mkdir -p $GOPATH/src/github.com/seaman
    $ cd $GOPATH/src/github.com/seaman
    $ git clone https://github.com/ZhongAnTech/seaman.git
    $ cd seaman
    $ make docker
    $ docker run -d --name seaman -p 8080:8080 seaman


### 编译源码

    $ mkdir -p $GOPATH/src/github.com/seaman
    $ cd $GOPATH/src/github.com/seaman
    $ git clone https://github.com/ZhongAnTech/seaman.git
    $ cd seaman
    $ make build
    $ ./seaman

## License

Apache License 2.0, see [LICENSE](LICENSE).
