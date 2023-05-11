<!-- # Envoyer -->
<div align="center">
<img src="https://github.com/vivasoft-ltd/Envoyer/blob/main/envoyer_backend/doc/logo.png" width="300" alt="Logo"/>
<h3 align="center">Simplify Notification Management with Customizable Templates and Multi-Provider Integration</h3>
</div>

## ⭐️ Why Envoyer
Nowadays, notifications play a crucial role in every type of application, be it a web or mobile application. Our software offers a notification system that can be seamlessly integrated with multiple applications to send notifications via email, SMS, push, and more using different providers. This notification service makes it easier for developers to manage notifications, including sending emails, SMS, and push notifications. Here application owners and developers can effortlessly handle notifications, add multiple providers, save multiple templates, and manage them with ease.

## 📚 Table Of Contents

- [Features](#-features)
- [Providers](#-providers)
  - [Email](#-email)
  - [SMS](#-sms)
  - [Push](#-push)
 - [Demo](#-demo)
- [Getting Started](#-getting-started)
    - [Local](#-local-environment)
    - [Docker](#-docker-environment)
- [User Manual](#-user-manual)
- [Code Of Conduct](#-read-our-code-of-conduct)
- [License](#-license)
- [Contributions](#-thanks-to-all-contributors)

## ✨ Features
- 📨 Supports email, SMS, and push notifications, as well as webhooks to send data to internal services.
- 📝 Provides the option to create multiple templates with multiple variables, and save multiple templates for different languages.
- 🌐 Can add multiple email, SMS or push providers, and prioritize them to ensure notification delivery even if one provider fails.
- 🚀 Supports various providers like SMTP, Firebase FCM, Twilio, Vonage, and more.
- 👨‍💼 Admin dashboard for application owners to modify templates or providers as needed, and for developers to easily add new events with variables.
- 📲 Simple way to trigger events with just an API call, and to send notifications for various events (sign-ups, new purchases, etc.) with various channels (emails, SMS, etc.) all in one place.
- 👥 Can send bulk notifications to a large audience and send notifications with a time delay.
- 🧩 Uses message queues to handle multiple notifications without losing any data, and can retry in case of failure.

[See more details](https://github.com/vivasoft-ltd/Envoyer/blob/main/envoyer_backend/doc/features.md)

## 📡 Providers
Envoyer provides simple user interface for manage multiple notifications providers across multiple channels.

#### :mailbox: Email
- SMTP

#### :inbox_tray: SMS
- Twilio
- Nexmo - Vonage

#### 📱 Push
- FCM

## 🎨 Demo
Integrate different provider

<img src="https://github.com/vivasoft-ltd/Envoyer/blob/main/envoyer_backend/doc/twilio.png"/>   <img src="https://github.com/vivasoft-ltd/Envoyer/blob/main/envoyer_backend/doc/smtp.png"/>

Add template


<img src="https://github.com/vivasoft-ltd/Envoyer/blob/main/envoyer_backend/doc/add_template.png"/>
<img src="https://github.com/vivasoft-ltd/Envoyer/blob/main/envoyer_backend/doc/template.png"/>

Save multiple template and manage them easily.

<img src="https://github.com/vivasoft-ltd/Envoyer/blob/main/envoyer_backend/doc/multi_template.png"/>

## 🚀 Getting Started

## 💻 Local environment

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

## 🐳 Docker environment

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

## 📋 User Manual
  Click [here](https://github.com/vivasoft-ltd/Envoyer/blob/main/envoyer_backend/doc/user_manual.md) to see the details user manual of envoyer.

Good Luck 👍


## 📜 Read Our Code Of Conduct

Before you begin coding and collaborating, please read our [Code of Conduct](https://github.com/vivasoft-ltd/Envoyer/blob/main/CODE_OF_CONDUCT.md) thoroughly to understand the standards (that you are required to adhere to) for community engagement. As part of our open-source community, we hold ourselves and other contributors to a high standard of communication. As a participant and contributor to this project, you are agreeing to abide by our [Code of Conduct](https://github.com/vivasoft-ltd/Envoyer/blob/main/CODE_OF_CONDUCT.md) and [Contribution Guidelines](https://github.com/vivasoft-ltd/Envoyer/blob/main/CONTRIBUTING.md).



## 📝 License

Envoyer is licensed under the Apache License 2.0 - see the [LICENSE](https://github.com/vivasoft-ltd/Envoyer/blob/main/LICENSE) file for details.

## 🤝 Thanks To All Contributors

Thanks a lot for spending your time helping Envoyer grow. Keep rocking 👍


<img src="https://contributors-img.web.app/image?repo=vivasoft-ltd/Envoyer" />

