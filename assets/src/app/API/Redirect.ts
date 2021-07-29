import axios from 'axios';
import token from './token';

interface APIParamsNewRedirect {
  from: string;
  after: string;

  to?: string;

  urlTemplate?: string;
  methodTemplate?: string;
  headersTemplate?: string;
  bodyTemplate?: string;
}


export class Redirect {
  FromURI: string;
  RedirectAfter: string;

  ToURL?: string;

  URLTemplate?: string;
  MethodTemplate?: string;
  HeadersTemplate?: string;
  BodyTemplate?: string;

  serialize(): APIParamsNewRedirect {
    let redirectData: APIParamsNewRedirect;
    if (this.ToURL === undefined || this.ToURL === null || this.ToURL === "") {
      redirectData = {
        from: this.FromURI,
        after: this.RedirectAfter,
        urlTemplate: this.URLTemplate,
        methodTemplate: this.MethodTemplate,
        headersTemplate: this.HeadersTemplate,
        bodyTemplate: this.BodyTemplate,
      };
    } else {
      redirectData = {
        from: this.FromURI,
        after: this.RedirectAfter,
        to: this.ToURL,
      };
    }
    return redirectData;
  }

  async createNew() {
    let resp = await axios.post("/api/new", this.serialize(), {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    if (resp.status === 200) {
      return null;
    }
    return resp.data.status;
  }
}
