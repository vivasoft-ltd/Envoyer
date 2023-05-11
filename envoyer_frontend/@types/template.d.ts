type CreateTemplateInput = {
  type: string;
  description?: string;
  message?: string;
  email_subject?: string;
  markup?: string;
  email_rendered_html?: string;
  event_id: number;
  active?: boolean;
  title?: string;
  link?: string;
  file?: string;
  language?: string;
};

type UpdateTemplateInput = CreateTemplateInput & {
  id: number;
};

type TemplateResponse = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  type: string;
  description?: string;
  message?: string;
  email_subject?: string;
  markup?: string;
  email_rendered_html?: string;
  event_id: number;
  active?: boolean;
  title?: string;
  link?: string;
  file?: string;
  language?: string;
};
