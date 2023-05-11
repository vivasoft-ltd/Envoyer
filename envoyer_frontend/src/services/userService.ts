import ApiClient from '@/lib/apiClient';

class UserService {
  private apiClient: ApiClient;

  constructor() {
    this.apiClient = new ApiClient(`/api/proxy`);
  }

  async getAllUsers(appId: number): Promise<UserResponse[]> {
    return this.apiClient.get('/v2/user/app_id/' + appId).then((response) => response.data.data);
  }

  async getUser(id: number): Promise<UserResponse> {
    return this.apiClient
      .get('/v2/user/' + id)
      .then((response) => response.data.data)
      .catch((err) => {
        throw err.data;
      });
  }

  async createUser(
    UserInput: CreateUserInput
  ): Promise<UserResponse> {
    return this.apiClient
      .post('/v2/user', {data: UserInput})
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }

  async deleteUser(id: number): Promise<null> {
    return this.apiClient
      .delete("v2/user/" + id)
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      })
  }

  async editUser(
    editUser: CreateUserInput, Id: number | undefined
  ): Promise<UserResponse> {
    return this.apiClient
      .put('/v2/user/' + Id, {data: editUser})
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }
}

export default UserService;
