import ApiClient from '@/lib/apiClient';

class EventService {
  private apiClient: ApiClient;

  constructor() {
    this.apiClient = new ApiClient(`/api/proxy`);
  }

  async getAllEvents(appId: number): Promise<EventResponse[]> {
    return this.apiClient.get('/v2/event/app_id/' + appId).then((response) => response.data.data);
  }

  async getEvent(id: number): Promise<EventResponse> {
    return this.apiClient
      .get('/v2/event/' + id)
      .then((response) => response.data.data)
      .catch((err) => {
        throw err.data;
      });
  }

  async createEvent(
    EventInput: CreateEventInput
  ): Promise<EventResponse> {
    return this.apiClient
      .post('/v2/event', {data: EventInput})
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }

  async deleteEvent(id: number): Promise<null> {
    return this.apiClient
      .delete("v2/event/" + id)
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      })
  }

  async editEvent(
    editEvent: CreateEventInput, Id: number | undefined
  ): Promise<EventResponse> {
    return this.apiClient
      .put('/v2/event/' + Id, {data: editEvent})
      .then((response) => {
        return response.data;
      })
      .catch((err) => {
        throw err.data;
      });
  }
}

export default EventService;
