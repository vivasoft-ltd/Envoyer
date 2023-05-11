type CreateEventInput = {
  name: string;
  description: string;
  variables?: string[];
  app_id: number;
};

type EventResponse = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  name: string;
  description: string;
  variables?: string[];
  app_id: number;
};
