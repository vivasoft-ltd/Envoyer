<!-- # Envoyer -->
<div align="center">
 <font size="7">Envoyer</font>
  <h2 align="center">Simplify Notification Management with Customizable Templates and Multi-Provider Integration</h2>
  <br>
  <br>
</div>

## ⭐️ Why Envoyer
Nowadays, notifications play a crucial role in every type of application, be it a web or mobile application. Our software offers a notification system that can be seamlessly integrated with multiple applications to send notifications via email, SMS, push, and more using different providers. This notification service makes it easier for developers to manage notifications, including sending emails, SMS, and push notifications. Here application owners and developers can effortlessly handle notifications, add multiple providers, save multiple templates, and manage them with ease.


## ✨ Features
  [Features](https://github.com/vivasoft-ltd/Envoyer/blob/main/envoyer_backend/doc/features.md)

## 📋 User Manual
  [User Manual](https://github.com/vivasoft-ltd/Envoyer/blob/main/envoyer_backend/doc/user_manual.md)

## Development environment

Go to `envoyer_backend` directory 

```shell
cd envoyer_backend
```

Copy `.env.template` and create `.env` file and write your local config in it.

```shell
cp .env.template .env
```

Run migration.

```shell
go run main.go migrate
```

Run Envoyer backend

Run rabitmq, mysql database and run the backend server with this command.

```shell
go run main.go server
```

Go to `envoyer_frontend` directory 

```shell
cd ..
cd envoyer_frontend
```

Run Envoyer frontend

Copy `.env.template` and create `.env` file and write your local config in it.

```shell
cp .env.template .env
```

Run the frontend server with this command.

```shell
yarn dev
```
Check http://localhost:8081/ping

Check http://localhost:3000

## Docker environment

Go to `envoyer_backend` directory.

```shell
cd envoyer_backend
```

Copy `.env.template` and create `.env` file and write your local config in it.

```shell
cp .env.template .env
```

Go to `envoyer_frontend` directory 

```shell
cd ..
cd envoyer_frontend
```

Copy `.env.template` and create `.env` file and write your local config in it.

```shell
cp .env.template .env
```

Go back to main directory

```shell
cd ..
```

Run Envoyer

```bash
docker-compose up
```

Check http://localhost:8081/ping

Check http://localhost:3000

Good Luck 👍


## Providers

Envoyer provides simple user interface for manage multiple notifications providers across multiple channels.

#### 💌 Email
- SMTP

#### 📞 SMS
- Twilio
- Nexmo - Vonage

#### 📱 Push
- FCM


## 📋 Read Our Code Of Conduct

Before you begin coding and collaborating, please read our [Code of Conduct](https://github.com/vivasoft-ltd/Envoyer/blob/main/CODE_OF_CONDUCT.md) thoroughly to understand the standards (that you are required to adhere to) for community engagement. As part of our open-source community, we hold ourselves and other contributors to a high standard of communication. As a participant and contributor to this project, you are agreeing to abide by our [Code of Conduct](https://github.com/vivasoft-ltd/Envoyer/blob/main/CODE_OF_CONDUCT.md).



## 🛡️ License

Envoyer is licensed under the Apache License 2.0 - see the [LICENSE](https://github.com/vivasoft-ltd/Envoyer/blob/main/LICENSE) file for details.

## 💪 Thanks To All Contributors

Thanks a lot for spending your time helping Novu grow. Keep rocking 🥂


<img src="https://contributors-img.web.app/image?repo=vivasoft-ltd/Envoyer" />

