type CreateProviderInput = {
  app_id: number;
  type: string;
  provider_type: string;
  name?: string;
  description?: string;
  config: object;
  priority?: number;
  active?: boolean;
  policy?: object;
};

type ProviderResponse = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  app_id: number;
  type: string;
  provider_type: string;
  name: string;
  description: string;
  config: object;
  priority: number;
  active: boolean;
  policy?: object;
};

type Priority = {
  id: number;
  priority: number;
}

type TwilioSmsConfig = {
  account_sid: string;
  auth_token: string;
  sender_id: string;
}

type VonageSmsConfig = {
  api_key: string;
  api_secret: string;
  sender_id: string;
}

type SmtpConfig = {
  smtp_host: string;
  smtp_port: number;
  smtp_user_name: string;
  smtp_password: string;
  sender: string;
}

type WebhookConfig = {
  url: string;
  token?: string;
}