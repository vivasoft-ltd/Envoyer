type ErrorLog = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  app_id?: number;
  client_key?: string;
  event_name?: string;
  provider_id?: number;
  message?: string;
  data?: object;
  request?: object;
  type?: string;
  date?: string;
  is_requeue?: boolean;
}