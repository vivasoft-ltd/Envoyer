import ApiClient from '@/lib/apiClient';

class ErrorLogService {
  private apiClient: ApiClient;

  constructor() {
    this.apiClient = new ApiClient(`/api/proxy`);
  }

  async getAllErrorLogs(appId: number): Promise<ClientResponse[]> {
    return this.apiClient.get('/v2/log/app_id/' + appId).then((response) => response.data.data);
  }

  async deleteLog(id: number): Promise<null> {
    return this.apiClient
      .delete("v2/log/" + id)
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      })
  }

}

export default ErrorLogService;
