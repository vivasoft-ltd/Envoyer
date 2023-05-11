export const USER_ROLE = {
  SUPER_ADMIN: 'super_admin',
  ADMIN: 'admin',
  DEVELOPER: 'dev',
};

export const ProviderType = {
  TWILIO: 'twilio',
  VONAGE: 'vonage',
  SMTP: 'smtp',
  FIREBASE: 'fcm',
  WEBHOOK: 'webhook',
}

export const NotificationType = {
  EMAIL: 'email',
  SMS: 'sms',
  PUSH: 'push',
  WEBHOOK: 'webhook',
}

export const SMS_Policy = [
  {value:"receiver.country", label:"Receiver Country"},
  {value:"receiver.prefix", label:"Receiver Prefix"}
];
