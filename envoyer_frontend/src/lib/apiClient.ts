import HttpClient from './httpClient';

let signoutTriggered = false;

class ApiClient {
  private httpClient: HttpClient;

  constructor(baseUrl: string, timeout = 20000) {
    this.httpClient = new HttpClient(baseUrl, timeout);
    // this.httpClient.addResponseInterceptor(
    //   (response) => response,
    //   (err) => {
    //     if (err?.response?.status === 401 && !signoutTriggered) {
    //       signoutTriggered = true;
    //       signOut({
    //         callbackUrl: `${window.location.origin}/auth/login?reason=session_expired`,
    //       });
    //     }

    //     return Promise.reject(err);
    //   }
    // );
  }

  setToken(token: string | null) {
    if (token) {
      this.httpClient.addHeader('Authorization', `Bearer ${token}`);
      return this;
    }

    this.httpClient.removeHeader('Authorization');
    return this;
  }

  setUserInfo(userInfo: any) {
    this.httpClient.addHeader('userInfo', `${JSON.stringify(userInfo)}`);
    return this;
  }

  async get(
    url: string,
    { preRender, v2 }: { preRender?: boolean; v2?: boolean } = {
      preRender: false,
      v2: false,
    }
  ): Promise<any> {
    const requestUrl = () => {
      if (preRender) {
        if (v2) return process.env.API_BASE_URL + url;

        return process.env.API_BASE_URL + '/v1' + url;
      }
      return url;
    };
    return this.httpClient.get(requestUrl()).then((response) => {
      return response;
    });
  }

  async post(url: string, Request: any): Promise<any> {
    return this.httpClient.post(url, Request).then((response) => {
      return response;
    });
  }

  async patch(url: string, Request: any): Promise<any> {
    return this.httpClient.patch(url, Request).then((response) => {
      return response;
    });
  }

  async put(url: string, Request: any): Promise<any> {
    return this.httpClient.put(url, Request).then((response) => {
      return response;
    });
  }

  async delete(url: string, Request?: any): Promise<any> {
    return this.httpClient.delete(url, Request).then((response) => {
      return response;
    });
  }
}

export default ApiClient;
