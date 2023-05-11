import ApiClient from '@/lib/apiClient';

class PublishService {
  private apiClient: ApiClient;

  constructor() {
    this.apiClient = new ApiClient('');
  }

  async publishInQueueCustom(
    custom: string,
    publishInput: PublishInput
  ): Promise<unknown> {
    return this.apiClient
      .post(`${process.env.API_BASE_URL!}/api/v2/publish/` + custom, {
        data: publishInput,
      })
      .then((response) => {
        return response;
      })
      .catch((err) => {
        throw err;
      });
  }
}

export default PublishService;
