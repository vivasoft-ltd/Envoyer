import ApiClient from '@/lib/apiClient';

class AppServices {
  private apiClient: ApiClient;

  constructor() {
    this.apiClient = new ApiClient(`/api/proxy`);
  }

  async getAllApps(): Promise<ListAppResponse[]> {
    return this.apiClient.get('/v2/app').then((response) => response.data.data);
  }

  async getApp(id: number): Promise<ListAppResponse> {
    return this.apiClient
      .get('/v2/app/' + id)
      .then((response) => response.data.data)
      .catch((err) => {
        throw err.data;
      });
  }

  async createApp(
    emailAndPassword: CreateAppInput
  ): Promise<CreateAppResponse> {
    return this.apiClient
      .post('/v2/app', {data: emailAndPassword})
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }

  async deleteApp(id: number): Promise<null> {
    return this.apiClient
      .delete("v2/app/" + id)
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      })
  }

  async editApp(
    editApp: CreateAppInput, Id: number | undefined
  ): Promise<CreateAppResponse> {
    return this.apiClient
      .put('/v2/app/' + Id, {data: editApp})
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }
}

export default AppServices;
