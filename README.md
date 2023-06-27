# Feishu-Shell-Bot
功能很简单。就是跑一个shell命令客户端。跟飞书相关

配置请配置.feishu.env.sample 然后重命名到 .feishu.env
或者也可以直接在环境变量中注入相应的配置

```conf
FEISHU_APP_ID=飞书的AppId
FEISHU_APP_SECRET=飞书的AppSecret
FEISHU_ENCRYPT_KEY=飞书的EncryptKey
FEISHU_VERIFICATION_TOKEN=飞书的验证Token
FEISHU_BOT_PATH=监听服务的path
FEISHU_BOT_PORT=监听的端口
```

具体配置的含义请参考飞书后台开发者的文档
具体飞书怎么配置直接参考项目
https://github.com/whatwewant/chatgpt-for-chatbot-feishu

