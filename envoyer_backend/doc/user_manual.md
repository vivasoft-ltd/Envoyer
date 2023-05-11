# Envoyer

## Getting Started

### 1. Create an app
Only super admin can create a new app with name and optional description.
 App name can only have letter, digits, underscores and 
hyphens. An app can be active or inactive. If an app is active
that means a consumer is running on a thread for 
consuming the messages that are published on that app. 
So for unused apps it is recommended to set the app inactive.

### 2. Create users
Only super admins can create users. All user need to be 
unique. There are two type of user associated for an app. 
Developer and application owner. Developers can 
create events, clients and application owner can only
create or update templates and providers.


## App Setup
User with developer role can set up the application with 
necessary configurations. No need to involve the super admin.

### 1. Create Events
Developers can create events with multiple variables. 
While triggering these events, event name and variables
along with their values need to be provided in the api call.
A variable name can only have letter, digits, underscores and hyphens.

### 2. Create Templates
Developers and application owners can create templates 
for an event. He can use the variables associated with
the event. They can create multiple templates for an event 
and set one active at a time. They can also update or 
delete them.

### 3. Create Providers
Developers and application owners can add multiple 
providers for email, sms, push etc. They can set them active as
well as prioritize them. So if one provider fails to
deliver the notification next active provider will do the task.

## Trigger Event
Developers can easily trigger an event from their 
application with an api call.
### 1. Create Client
For trigger an events we need a client of the app. 
Developers can create clients with 
a name and optional description. It can be frontend, backend or 
any service that will trigger the event. 
After creating a client 
we can get app key and client key that are needed to
trigger an event.

### 2. Trigger an Event
After creating a client, adding an active provider and adding 
an event that has active template, we can trigger that event.
To trigger the event we just need one api call. 

The api call will be a POST request like this:

`http://{url}/api/publish/{type}`

Here the `url` will be the url of the envoyer and `type` can
be `email`, `sms`, `push` or `webhook`. The body of the request will be like this:
```json
{
    "app_key" : "your_app_key",
    "client_key" : "your_client_key",
    "event_name": "your_event_name",
    "receivers": [
        "example@example.com"
    ],
    "variables": [
        {
            "name": "{{.FirstName}}",
            "value": "John"
        },
        {
            "name": "{{.LastName}}",
            "value": "Doe"
        }
    ]
}
```

The `receivers` can be emails, phone numbers (with country codes) etc.
according to the `type`. It is an array, 
so we can send to multiple receivers if we want. 

If we want to send to multiple receivers with 
individual variables values, the body will be like this:
```json
{
    "app_key" : "your_app_key", 
    "client_key" : "your_client_key",
    "event_name": "your_event_name",
    "receivers_with_variables": [
        {
            "receiver": "johndoe@example.com",
            "variables": [
                {
                   "name": "{{.FirstName}}",
                   "value": "John"
                },
               {
                  "name": "{{.LastName}}",
                  "value": "Doe"
               }
            ]
        },
        {
            "receiver": "janedoe@example.com",
            "variables": [
               {
                  "name": "{{.FirstName}}",
                  "value": "Jane"
               },
               {
                  "name": "{{.LastName}}",
                  "value": "Doe"
               }
            ]
        }
    ]
}
```

For every event we can have an active template  for 
each language. So for an event we can have an active 
template in English (en), French (fr), Bengali (bn) etc. 
While triggering the event we can specify which language
template we want to use. For this we need to add 
`language` field in the request body like this:

```json
{
    "app_key" : "your_app_key",
    "client_key" : "your_client_key",
    "event_name": "your_event_name",
    "receivers": [
        "example@example.com"
    ],
    "variables": [
        {
            "name": "{{.FirstName}}",
            "value": "John"
        },
        {
            "name": "{{.LastName}}",
            "value": "Doe"
        }
    ],
    "language": "bn"
}
```

If we don't give `language` field in the body, by default it 
will use the template with English (en) language.

#### Email
For email, we can override the sender in the body for AWS SES.

```json
    "sender": "example@example.com"
```

We can also add cc and bcc if needed only in this format:
```json
{
    "app_key": "your_app_key",
    "client_key": "your_client_key",
    "event_name": "your_event_name",
    "receivers": [
        "example@example.com"
    ],
    "variables": [
        {
            "name": "{{.FirstName}}",
            "value": "John"
        },
        {
            "name": "{{.LastName}}",
            "value": "Doe"
        }
    ],
    "cc": [
        "example.cc1@example.com",
        "example.cc2@example.com"
    ],
    "bcc": [
        "example.bcc1@example.com",
        "example.bcc2@example.com"
    ]
}
```


#### Push
For push notification the body will look like this:

```json
{
    "app_key" : "your_app_key",
    "client_key" : "your_client_key",
    "event_name": "your_event_name",
    "receivers": [
        "registration_token_of_receiver"
    ],
    "variables": [
        {
            "name": "{{.FirstName}}",
            "value": "John"
        },
        {
            "name": "{{.LastName}}",
            "value": "Doe"
        }
    ], 
    "language": "bn"
}
```

We can override the image url with `image_url` field, and we
can also use variables in it.  

```json
    "image_url": "https://picsum.photos/200"
```
We can also use `topic` or `condition` instead of `receivers` 
as well as additional `data` (as key value pair)  in the
body for firebase fcm.

```json
    "topic": "your_topic",
    "data": {
        "key": "value",
        "key2": "value2"
    }
```
or,
```json
    "condition": "your_condition"
```

But you can not use `topic` or `condition` in this format.

```json
{
    "app_key" : "your_app_key", 
    "client_key" : "your_client_key",
    "event_name": "your_event_name",
    "receivers_with_variables": [
        {
            "receiver": "registration_token_of_receiver_1",
            "variables": [
                {
                   "name": "{{.FirstName}}",
                   "value": "John"
                },
               {
                  "name": "{{.LastName}}",
                  "value": "Doe"
               }
            ]
        },
        {
            "receiver": "registration_token_of_receiver_2",
            "variables": [
               {
                  "name": "{{.FirstName}}",
                  "value": "Jane"
               },
               {
                  "name": "{{.LastName}}",
                  "value": "Doe"
               }
            ]
        }
    ],
    "image_url": "https://picsum.photos/200",
    "data": {
        "key": "value"
    }
}
```

#### Webhook
We can add a webhook provider with an url and optional 
bearer token. We can send any data to that url like this:

```json
{
    "app_key" : "your_app_key",
    "client_key" : "your_client_key",
    "data": {
        "message": "Hello Mr John Doe",
        "info": {
            "text": "This is a text",
            "key": "value"
        }
    }
}
```
Here the `data` can be any json object, and it will be 
sent to the url with a POST request. It will use bearer 
token for authorization if it is given in the provider 
config.

Success response will be like this with status code 200.
```json
{
    "data": "Successfully published to the queue",
    "status": "ok"
}
```



