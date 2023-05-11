import ApiClient from '@/lib/apiClient';

class AuthService {
  private apiClient: ApiClient;

  constructor() {
    this.apiClient = new ApiClient(`${process.env.API_BASE_URL!}/api`);
  }

  async authenticate(emailAndPassword: LoginInput): Promise<AuthResponse> {
    return this.apiClient
      .post('/v2/auth/login', { data: emailAndPassword })
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }

  async refreshToken(refresh_token: string): Promise<RefreshTokenResponse> {
    return this.apiClient
      .post('/v2/auth/refresh', { data: { refresh_token } })
      .then((response) => response.data.data);
  }
}

export default AuthService;
