## 体外部署

### 部署prometheus

docker部署：
```bash
$ docker run -d --name pm \
  -p 9090:9090 \
  -v $(pwd)/config:/config \
  prom/prometheus:v2.30.0 --web.enable-lifecycle --config.file=/config/prometheus.yml
```

重载配置 `POST http://主机IP:9090/-/reload`

访问`主机IP:9090`

### 配置k8s基础组件指标

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


### 获取基础组件指标

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

### 配置node指标

主要内容为CPU，内存，硬盘，I/O等




