
* 获取`service`选择的`pod`名称：`kubectl get endpoints $服务名称 -o=jsonpath='{.subsets[*].addresses[*].ip}' | tr ' ' '\n' | kubectl get pods --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'`