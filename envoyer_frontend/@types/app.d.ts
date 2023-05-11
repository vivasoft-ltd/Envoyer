type CreateAppInput = {
  name: string;
  description: string;
  active: boolean;
};

type CreateAppResponse = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  name: string;
  description: string;
  app_key: string;
  active: boolean;
};

type ListAppResponse = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  name: string;
  description: string;
  app_key: string
  active: boolean;
}
