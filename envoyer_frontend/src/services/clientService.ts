import ApiClient from '@/lib/apiClient';

class ClientService {
  private apiClient: ApiClient;

  constructor() {
    this.apiClient = new ApiClient(`/api/proxy`);
  }

  async getAllClients(appId: number): Promise<ClientResponse[]> {
    return this.apiClient.get('/v2/client/app_id/' + appId).then((response) => response.data.data);
  }

  async getClient(id: number): Promise<ClientResponse> {
    return this.apiClient
      .get('/v2/client/' + id)
      .then((response) => response.data.data)
      .catch((err) => {
        throw err.data;
      });
  }

  async createClient(
    ClientInput: CreateClientInput
  ): Promise<ClientResponse> {
    return this.apiClient
      .post('/v2/client', {data: ClientInput})
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }

  async deleteClient(id: number): Promise<null> {
    return this.apiClient
      .delete("v2/client/" + id)
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      })
  }

  async editClient(
    editClient: CreateClientInput, Id: number | undefined
  ): Promise<ClientResponse> {
    return this.apiClient
      .put('/v2/client/' + Id, {data: editClient})
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }
}

export default ClientService;
