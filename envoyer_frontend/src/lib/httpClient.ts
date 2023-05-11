import axios, {
  AxiosInstance,
  AxiosResponse,
  Method,
  ResponseType,
} from 'axios';

interface HttpClientRequest {
  method?: Method;
  url?: string;
  headers?: any;
  params?: any;
  data?: any;
  timeout?: number;
  responseType?: ResponseType;
  onUploadProgress?: (progressEvent: any) => void;
}

interface HttpClientResponse<TData> {
  status: number;
  data: TData;
}

export class HttpClientError<TError = any> extends Error {
  status: number;
  data: TError;

  constructor(message: string, status: number, data: TError) {
    super(message);
    this.status = status;
    this.data = data;
  }
}

class HttpClient {
  private axiosInstance: AxiosInstance;

  constructor(baseUrl = '', timeout = 20000) {
    this.axiosInstance = axios.create({
      withCredentials: true,
      baseURL: baseUrl,
      timeout,
      headers: {
        Accept: 'application/json',
      },
    });
  }

  private async request<TData>(
    config: HttpClientRequest
  ): Promise<HttpClientResponse<TData>> {
    return this.axiosInstance
      .request(config)
      .then((response) => {
        return {
          status: response.status,
          data: response.data,
        };
      })
      .catch((err) => {
        const status = err?.response?.status || 0;
        const data = err?.response?.data || null;
        console.log(err.message, status, data);

        throw new HttpClientError(err.message, status, data);
      });
  }

  async get<TData = any>(url: string, config: HttpClientRequest = {}) {
    return this.request<TData>({ method: 'get', url, ...config });
  }

  async post<TData = any>(url: string, config: HttpClientRequest = {}) {
    console.log('post', url);
    return this.request<TData>({ method: 'post', url, ...config });
  }

  async put<TData = any>(url: string, config: HttpClientRequest = {}) {
    return this.request<TData>({ method: 'put', url, ...config });
  }

  async patch<TData = any>(url: string, config: HttpClientRequest = {}) {
    return this.request<TData>({ method: 'patch', url, ...config });
  }

  async delete<TData = any>(url: string, config: HttpClientRequest = {}) {
    return this.request<TData>({ method: 'delete', url, ...config });
  }

  addHeader(name: string, value: string) {
    this.axiosInstance.defaults.headers.common[name] = value;
  }

  removeHeader(name: string) {
    delete this.axiosInstance.defaults.headers.common[name];
  }

  addResponseInterceptor = (
    successCallback: (
      value: AxiosResponse<any>
    ) => AxiosResponse<any> | Promise<AxiosResponse<any>>,
    errorCallback: (err: any) => any
  ) => {
    this.axiosInstance.interceptors.response.use(
      successCallback,
      errorCallback
    );
  };
}

export default HttpClient;
