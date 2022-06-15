## 体外部署prometheus

### 部署

docker部署：
```bash
$ docker run -d --name pm \
  -p 9090:9090 \
  -v $(pwd)/config:/config \
  prom/prometheus:v2.30.0 --web.enable-lifecycle --config.file=/config/prometheus.yml
```

重载配置 `POST http://主机IP:9090/-/reload`

访问`主机IP:9090`

### 1. 配置k8s基础组件指标

如`pod`,`svc`,`configmap`等资源的指标

https://github.com/kubernetes/kube-state-metrics

查看并切换到适合k8s版本的`tag`(当前v2.3.0)，进入`examples/standard`，拷贝文件至本地`kube-state-metrics`

修改了两个地方`deployment.yaml`内`- image: bitnami/kube-state-metrics:2.3.0` 方便国内网络

另一个`service.yaml`内，允许外部访问:

```yaml
spec:
#  clusterIP: None
  type: NodePort
  ports:
    - name: http-metrics
      port: 8080
      targetPort: http-metrics
      nodePort: 32280
```

进入`kube-state-metrics`文件夹，部署 `kb apply -f .`


#### 手动配置获取基础组件指标

在 `config/prometheus.yml`内新增：

```yaml
global:
  scrape_interval: 5s
scrape_configs:
  - job_name: "prometheus-state-metrics"
    static_configs:
      - targets: ["$(k8sIP):32280"]
```

`targets`参数地址为`nodePort`暴露的地址

重载配置`curl -X POST http://主机IP:9090/-/reload`

查看添加结果：`http://主机IP:9090/service-discovery`

查看指标： 打开`http://主机IP:9090/classic/graph`，选择复选框内的指标选项就能查看值

### 2. 配置节点指标

主要内容为CPU，内存，硬盘，I/O等

部署 `kb apply -f node-exporter/install.yaml`

查看部署结果 `http://{k8sIP}:9100`

#### 手动配置获取节点指标

在 `config/prometheus.yml`内新增：

```yaml
scrape_configs:
  - job_name: "node-exporter"
    static_configs:
      - targets: ["$(k8sIP):9100"]
```

重载配置`curl -X POST http://主机IP:9090/-/reload`

### 自动发现

上面获取都是通过手动的方式，可以配置自动发现，两种都可以使用，只是自动发现更灵活

#### 自动获取节点指标

创建一个集群角色`kb apply -f config/rbac.yaml`

因为是体外部署`prometheus`，所以要在设置`sa token`来请求，提取token: 

`kubectl -n kube-system describe secret $(kubectl -n kube-system describe sa my-prometheus | grep 'Mountable secrets' | cut -f 2- -d ":" | tr -d " ") | grep -E '^token' | cut -f2 -d':' | tr -d '\t'`

将内容保存至`sa.token`

新增`prometheus.yml`配置

```yaml
  - job_name: "k8s-node"
  # 查看配置文件
  ...
```

此时获取的内容与手动配置`job_name: "node-exporter"`的内容相同，可以注释掉手动配置，避免重复获取

重载配置`curl -X POST http://主机IP:9090/-/reload`

#### 自动获取pod指标

`kubelet`默认集成了`cAdvisor`单节点内部容器和节点资源使用统计

实例 `kb get --raw "/api/v1/nodes/{nodename}/proxy/metrics/cadvisor"`

查看配置文件`config/prometheus.yml`内`  - job_name: "k8s-kubelet"` 的配置