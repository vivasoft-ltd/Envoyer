type CreateUserInput = {
  user_name: string;
  password: string;
  app_id: number;
  role: string;
};

type UserResponse = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  user_name: string;
  password: string;
  app_id: number;
  role: string;
};
