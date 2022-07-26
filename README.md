# qinchy-scheduler
## 编译指南
* 先克隆kubernetes到本地，切换到需要适配的版本的分支。
* 修改go.mod中replace中的链接路径，即本地kubernetes远吗路径
* 执行go mod tidy，引用项目需要的依赖增加到go.mod文件。 去掉go.mod文件中项目不需要的依赖。
* 执行go mod vendor，将依赖拷贝到vendor目录
* 执行go build -o qinchy-scheduler，编译得到目标可执行文件。
* 执行docker build -t ${TAG} .，构建镜像，并推送到镜像仓库。
* 在k8s集群中使用`qinchy-scheduler-clusterrole.yml`创建`ClusterRole`,
  使用`qinchy-scheduler-clusterrolebinding.yml`创建`ClusterRoleBinding`,
  使用`qinchy-scheduler-scheduler-config.yaml`文件创建`configmap`,
  使用`qinchy-scheduler-serviceaccount.yml`创建`ServiceAccount`,
* 使用`qinchy-scheduler-deployment.yml`创建调度器的Deployment

  