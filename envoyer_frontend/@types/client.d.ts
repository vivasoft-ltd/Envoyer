type CreateClientInput = {
  name: string;
  description?: string;
  app_id: number;
};

type ClientResponse = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  name: string;
  description?: string;
  app_id: number;
  client_key: string;
};
