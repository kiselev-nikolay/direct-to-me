import axios from 'axios';

interface APIRedirect {
  from: string;
  after: string;
  to: string;
  urlTemplate: string;
  methodTemplate: string;
  headersTemplate: string;
  bodyTemplate: string;
}

export class Redirect {
  FromURI: string;
  RedirectAfter: string;

  ToURL?: string;

  URLTemplate?: string;
  MethodTemplate?: string;
  HeadersTemplate?: string;
  BodyTemplate?: string;

  isAdvanced: boolean;

  deserialize(data: APIRedirect) {
    this.FromURI = data.from;
    this.RedirectAfter = data.to;
    this.isAdvanced = data.to === "";
    if (this.isAdvanced) {
      this.URLTemplate = data.urlTemplate;
      this.MethodTemplate = data.methodTemplate;
      this.HeadersTemplate = data.headersTemplate;
      this.BodyTemplate = data.bodyTemplate;
    } else {
      this.ToURL = data.to;
    }
  }
}

export default async function GetRedirects() {
  let res = await axios.get("/api/list");
  return res.data.redirects.map((i: APIRedirect) => { let r = new Redirect(); r.deserialize(i); return r; });
}