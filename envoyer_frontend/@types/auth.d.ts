type LoginInput = {
  username: string;
  password: string;
};

type AuthResponse = {
  data: {
    expired_at: number;
    access_token: string;
    refresh_token: string;
    id: number;
    role: string;
    app_id: number;
  };
};

type RefreshTokenResponse = {
  access_token: string;
  refresh_token: string;
  token_type: string;
  status?: number;
};
