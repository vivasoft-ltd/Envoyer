import ApiClient from '@/lib/apiClient';

class TemplateService {
  private apiClient: ApiClient;

  constructor() {
    this.apiClient = new ApiClient(`/api/proxy`);
  }

  async getAllTemplates(eventId: number): Promise<TemplateResponse[]> {
    return this.apiClient
      .get('/v2/template/event_id/' + eventId)
      .then((response) => response.data.data);
  }

  async getTemplate(id: number): Promise<TemplateResponse> {
    return this.apiClient
      .get('/v2/template/' + id)
      .then((response) => response.data.data)
      .catch((err) => {
        throw err.data;
      });
  }

  async createTemplate(
    TemplateInput: CreateTemplateInput
  ): Promise<TemplateResponse> {
    return this.apiClient
      .post('/v2/template', { data: TemplateInput })
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }

  async deleteTemplate(id: number): Promise<null> {
    return this.apiClient
      .delete('v2/template/' + id)
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }

  async editTemplate(
    editTemplate: UpdateTemplateInput
  ): Promise<TemplateResponse> {
    return this.apiClient
      .put('/v2/template/' + editTemplate.id, { data: editTemplate })
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }
}

export default TemplateService;
