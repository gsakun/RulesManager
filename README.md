# Ruler规则接口

## ruler数据格式

```json
{
	"groups": [{         
		"name": "string",   //规则分组名称
		"rules": [{
			"alert": "string",   //告警规则名称
			"expr": "string",    // 告警判断表达式
			"for": "string",     //触发告警时间
			"labels": [
                {
                    "severity": "string" //告警等级 
                }
            ], //自定义label severity必填
			"annotations": "map[string]string" //自定义注释
		}]
	}]
}
```

## 示例yaml文件

```yaml
groups:
- name: example
  rules:
  - alert: HighRequestLoad
    expr: rate(http_request_total{pod="p1"}[5m]) > 1000
    for: 1m
    labels:
      severity: warning
    annotations:
      info: High Request Load
```

## 解释说明

在一个规则文件中可以指定若干个group，每个group内可以指定多条告警规则。一般来说，一个group中的告警规则之间会存在某种逻辑上的联系，但即使它们毫无关联，对后续的流程也不会有任何影响。而一条告警规则中包含的字段及其含义如下：

1. `alert`: 告警名称
2. `expr`: 告警的触发条件，本质上是一条promQL查询表达式，Prometheus Server会定期（一般为15s）对该表达式进行查询，若能够得到相应的时间序列，则告警被触发
3. `for`: 告警持续触发的时间，因为数据可能存在毛刺，Prometheus并不会因为在`expr`第一次满足的时候就生成告警实例发送到AlertManager。比如上面的例子意为名为"p1"的Pod，每秒接受的HTTP请求的数目超过1000时触发告警且持续时间为一分钟，若告警规则每15s评估一次，则表示只有在连续四次评估该Pod的负载都超过1000QPS的情况下，才会真正生成告警实例。
4. `labels`: 用于附加到告警实例中的标签，Prometheus会将此处的标签和评估`expr`得到的时间序列的标签进行合并作为告警实例的标签。告警实例中的标签构成了该实例的唯一标识。事实上，告警名称最后也会包含在告警实例的label中，且key为"alertname"。
5. `annotations`: 用于附加到告警实例中的额外信息，Prometheus会将此处的annotations作为告警实例的annotations，一般annotations用于指定告警详情等较为次要的信息

需要注意的是，一条告警规则并不只会生成一类告警实例，例如对于上面的例子，可能有如下多条时间序列满足告警的触发条件，即n1和n2这两个namespace下名为p1的pod的QPS都持续超过了1000：

```
http_request_total{namespace="n1", pod="p1"}
http_request_total{namespace="n2", pod="p1"}
```

最终生成的两类告警实例为：

```
# 此处只显示实例的label
{alertname="HighRequestLoad", severity="warning", namespace="n1", pod="p1"}
{alertname="HighRequestLoad", severity="warning", namespace="n2", pod="p1"}
```

因此，例如在K8S场景下，由于Pod具有易失性，我们完全可以利用强大的promQL语句，定义一条Deployment层面的告警，只要其中任何的Pod满足触发条件，都会产生对应的告警实例。

## 接口

### 查询

```json
uri : /ruler/querylist  //查询规则列表
method: GET
payload: {
	"groups": [{         
		"name": "string",   //规则分组名称
		"rules": [{
			"alert": "string",   //告警规则名称
			"expr": "string",    // 告警判断表达式
			"for": "string",     //触发告警时间
			"labels": "map[string]string", //自定义label
			"annotations": "map[string]string" //自定义注释
		}]
	}]
}

uri : /ruler/querydetail //查询规则详情 
method: POST
payload: {
	"groupname": string,
    "rulername": string
}
```

### 添加

```
url: /ruler/add
method: POST
payload : {
	"groups": [{         
		"name": "string",   //规则分组名称
		"rules": [{
			"alert": "string",   //告警规则名称
			"expr": "string",    // 告警判断表达式
			"for": "string",     //触发告警时间
			"labels": "map[string]string", //自定义label
			"annotations": "map[string]string" //自定义注释
		}]
	}]
}
```

### 更新

```
url: /ruler/update
method: POST
payload : {
	"groups": [{         
		"name": "string",   //规则分组名称
		"rules": [{
			"alert": "string",   //告警规则名称
			"expr": "string",    // 告警判断表达式
			"for": "string",     //触发告警时间
			"labels": "map[string]string", //自定义label
			"annotations": "map[string]string" //自定义注释
		}]
	}]
}
```



### 删除

```
url: /ruler/delete
method: POST
payload: {
	"groups": [{         
		"name": "string",   //规则分组名称
		"rules": [{
			"alert": "string",   //告警规则名称
		}]
	}]
}
```

## 常用规则参考链接

https://awesome-prometheus-alerts.grep.to/rules.html