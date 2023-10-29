打包之前准备工作：

1. 配置env文件：

   新建一个.env文件，填入一下参数：

   ```shell
   QQ_BOT_APPID=botAppID
   QQ_BOT_TOKEN=机器人令牌
   
   BAIDU_CLIENT_ID=APIkey
   BAIDU_CLIENT_SECRET=Secret Key
   ```

   qq的参数在https://q.qq.com/bot/#/developer/developer-setting获取

   百度的在千帆大模型平台的应用接入 获取

2. 打包

   mac用户

   ```bash
   GOOS=linux GOARCH=amd64 go build -o qqbot
   ```

   win用户

   ```bash
   SET GOOS=linux
   SET GOARCH=amd64
   go build -o qqbot
   ```

3. 服务器部署

   上传到服务器，并在文件所在目录执行

   ```bash
   nohup ./qqbot
   ```

   

