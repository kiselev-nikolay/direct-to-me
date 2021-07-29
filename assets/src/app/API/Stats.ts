import axios from 'axios';
import token from './token';

export interface APIClick {
  direct: number;
  social: number;
}

export interface APIFail {
  notFound: number;
  databaseUnreachable: number;
  templateProcessFailed: number;
  clientContentProcessFailed: number;
}

export interface APIStats {
  clicks: APIClick;
  fails: APIFail;
}

export default async function GetStats() {
  let res = await axios.get("/api/stats", {
    headers: {
      'Authorization': `Bearer ${token}` 
    }
  });
  let data = new Map<string, APIStats>();
  for (const [key, value] of Object.entries(res.data.clicks)) {
    if (!data.has(key)) {
      data.set(key, { clicks: value as APIClick, fails: null });
      continue;
    }
    let v = data.get(key);
    data.set(key, { clicks: value as APIClick, fails: v.fails });
  }
  for (const [key, value] of Object.entries(res.data.fails)) {
    if (!data.has(key)) {
      data.set(key, { clicks: null, fails: value as APIFail });
      continue;
    }
    let v = data.get(key);
    data.set(key, { clicks: v.clicks, fails: value as APIFail });
  }
  return data;
}