# Envoyer

## Overview
In today's world, notifications have become an integral part of every application, whether it
is a web application or a mobile application. Our software provides a notification system
that can be integrated with multiple applications to send notifications via email, SMS, push etc.
using various providers. It is a notification service that will help developers to handle the 
notifications like sending email, sms, push etc. 
Multiple applications can connect to this service to manage 
their notifications. Application owners or developers can
add multiple providers and can save multiple notification 
templates and easily manage them.

## Features
* Supports email, SMS, and push notifications.
* Also supports webhook feature by which one can 
 send data to other internal services as well.
* Provides the option to create multiple templates with multiple
  variables.
* Can save multiple templates for different language.
* Can add multiple email, sms or push providers and
  prioritize them. So if one provider fails to
  send notification, another provider will do the task.
* Can send notifications using various providers like smtp, firebase fcm, twilio, vonage etc.
* Admin dashboard by which application owner can modify the templates or modify providers as needed.
* From dashboard developers can easily add new events with variables.
* Simple way to trigger the events with just an api call.
* It will let you trigger notifications for various events
  (sign-ups, new purchases etc.) with various
  channels (emails, SMS etc.) all in one place.
* One can send a bulk number of notifications to a large audience.
* One can send a notification with some time delay.
* It uses message queues so that it can handle multiple notifications
 without losing any data. Also, it can do retry in case of failure.


## Goals 
* Offers a flexible and scalable solution for sending notifications.
* Developers can easily create events, variables to handle notifications with UI.
* Can easily integrate new apps.
* Save developer time.
* Application owners can modify templates as they need. No need to consult
developers.
* Can integrate various providers, and use them with the same interface.

## Technology Used
* RabbitMQ as the message broker 
* NextJS as the front-end framework 
* Go as the back-end language and 
* MySQL as the database


