#QKC监视器
```
    用来监控节点的运行情况
    使用方法：go run main.go ./cofig.json
```
##config
```
    {
      "RulerList": [{
        "Module": "BlockTime",     //检查类型
        "Hosts": ["http://13.228.159.171:38391","http://52.194.81.124:38391"],     //要检查的IpList
        "Interval": 10,     //检查的频率
        "WeChatDetail": {     //微信公众号参数
          "AppID": "wx900ca629d0906a34",
          "AppSecret": "1d51ae95382ce4a4b81885b785f469fa",
          "TemplateId": "oVzLVIFdm_joCZyWFrPEwDnNcxGyrU6I_UNQiNpSYs8"
        },
        "AlertLevel": "P0",     //预警级别
        "Email": null,     //要发送到的邮箱列表
        "Extra": {     //每个Module私有的配置
          "Interval": 10,
          "MaxBlock": 1
        }
      },{
        "Module": "PeerNumber",
        "Hosts": ["http://13.228.159.171:38391","http://52.194.81.124:38391"],
        "Interval": 10,
        "WeChatDetail": {
          "AppID": "wx900ca629d0906a34",
          "AppSecret": "1d51ae95382ce4a4b81885b785f469fa",
          "TemplateId": "oVzLVIFdm_joCZyWFrPEwDnNcxGyrU6I_UNQiNpSYs8"
        },
        "AlertLevel": "P1",
        "Email": null,
        "Extra": {
          "Interval": 10,
          "MinPeer": 100
        }
      }
      ]
    }
``` 

##模块
```
    type RulerI interface {
    	Check() []string
    	PreCheck() error
    }
    如果需要增加一个模块，需要实现这两个方法
             Check: 定期检查函数，返回错误类型，
             PreCheck: 程序启动时的预先检查，检查节点访问情况等
```