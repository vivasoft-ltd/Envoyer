import ApiClient from '@/lib/apiClient';

class ProviderService {
  private apiClient: ApiClient;

  constructor() {
    this.apiClient = new ApiClient(`/api/proxy`);
  }

  async getAllProviderByAppId(appId: number): Promise<ProviderResponse[]> {
    return this.apiClient.get('/v2/provider/app_id/' + appId).then((response) => response.data.data);
  }

  async getProvidersByAppIdAndType(appId: number, type: string): Promise<ProviderResponse[]> {
    return this.apiClient.get('/v2/provider/app_id/' + appId + "/type/" + type).then((response) => response.data.data);
  }

  async getProvider(id: number): Promise<ProviderResponse> {
    return this.apiClient
      .get('/v2/provider/' + id)
      .then((response) => response.data.data)
      .catch((err) => {
        throw err.data;
      });
  }

  async createProvider(
    ProviderInput: CreateProviderInput
  ): Promise<ProviderResponse> {
    return this.apiClient
      .post('/v2/provider', {data: ProviderInput})
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }

  async deleteProvider(id: number): Promise<null> {
    return this.apiClient
      .delete("v2/provider/" + id)
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      })
  }

  async editProvider(
    editProvider: CreateProviderInput, Id: number | undefined
  ): Promise<ProviderResponse> {
    return this.apiClient
      .put('/v2/provider/' + Id, {data: editProvider})
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }

  async editProviderPriority(
    editProvider: Priority[], appId: number, type: string
  ): Promise<ProviderResponse[]> {
    let reqData = {priority: editProvider}
    return this.apiClient
      .put('/v2/provider/app_id/' + appId + "/type/" + type, {data: reqData})
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }
}

export default ProviderService;
