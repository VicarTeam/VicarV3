import axios from "axios";
import {getSession, setSession} from "./session.ts";
import {useMainStore} from "@/stores/main.ts";

/**
 * A time in milliseconds which will be subtracted from the session's expiration time to determine if the session is expired.
 * It is used to prevent the session from expiring while the user is still using the application.
 */
const EXPIRE_THRESHOLD = 10_000;

const client = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  withCredentials: true,
});
const refreshClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  withCredentials: true,
});

function extractTokenFromQuery(): string|null {
  const query = new URLSearchParams(window.location.search);
  return query.get('act');
}

client.interceptors.request.use(async (config) => {
  const store = useMainStore();
  const deviceName = await store.getDeviceName();
  if (deviceName) {
    config.headers['X-Device-Name'] = deviceName;
  }

  const act = extractTokenFromQuery();
  if (act) {
    config.headers.Authorization = `Bearer ${act}`;
    return config;
  }

  const session = getSession();
  if (session) {
    if (session.expires - EXPIRE_THRESHOLD < Date.now()) {
      const data = (await refreshClient.post<{token: string, expiresIn: number}>('/auth/refresh', {}, {
        headers: {
          'X-Device-Name': deviceName,
        }
      })).data;
      setSession(data.token, data.expiresIn);
      session.token = data.token;
    }

    config.headers.Authorization = `Bearer ${session.token}`;
  }

  return config;
});

function handleRestError(error: any) {
  if (error.response) {
    if (error.response.data.error) {
      console.error(`REST[${error.response.status}-${error.response.statusText}]: ${error.response.data.error}`, error.response.data);
    } else {
      console.error(`REST[${error.response.status}-${error.response.statusText}]: `, error.response.data);
    }

    return;
  }

  console.error(error);
}

export async function GET_AND_THROW<T = any>(url: string, params?: {[key: string]: any}): Promise<T> {
  try {
    const response = await client.get<T>(url, { params });
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw new RestError(error.response.statusText, error.response.status, error.response.data);
    }

    throw error;
  }
}

export async function GET<T = any>(url: string, params?: {[key: string]: any}, defVal?: T): Promise<T> {
  try {
    return await GET_AND_THROW<T>(url, params);
  } catch (error: any) {
    handleRestError(error);
    return (defVal || null) as any;
  }
}

export type PaginationOptions = {search?: string, page?: number, itemsPerPage?: number};
export type PaginationResult<T> = {total: number, items: T[]};
export async function GET_PAGINATED<T = any>(url: string, opts: PaginationOptions = {}): Promise<PaginationResult<T>> {
  const {search, page, itemsPerPage} = opts;
  const params: {[key: string]: any} = {};
  if (search) {
    params.search = search;
  }
  if (page) {
    params.page = page;
  }
  if (itemsPerPage) {
    params.itemsPerPage = itemsPerPage;
  }

  try {
    return await GET_AND_THROW<PaginationResult<T>>(url, params);
  } catch (error: any) {
    handleRestError(error);
    return {total: 0, items: []};
  }
}

export async function POST_AND_THROW<T = any>(url: string, data?: any, params?: {[key: string]: any}): Promise<T> {
  try {
    const response = await client.post<T>(url, data, { params });
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw new RestError(error.response.statusText, error.response.status, error.response.data);
    }

    throw error;
  }
}

export async function POST<T = any>(url: string, data?: any, params?: {[key: string]: any}, defVal?: T): Promise<T> {
  try {
    const res = await POST_AND_THROW<T>(url, data, params);
    if (res) {
      return res;
    }
    return (defVal || true) as any;
  } catch (error: any) {
    handleRestError(error);
    return (defVal || null) as any;
  }
}

export async function PUT_AND_THROW<T = any>(url: string, data?: any, params?: {[key: string]: any}): Promise<T> {
  try {
    const response = await client.put<T>(url, data, { params });
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw new RestError(error.response.statusText, error.response.status, error.response.data);
    }

    throw error;
  }
}

export async function PUT<T = any>(url: string, data?: any, params?: {[key: string]: any}, defVal?: T): Promise<T> {
  try {
    const res = await PUT_AND_THROW(url, data, params);
    if (res) {
      return res;
    }
    return (defVal || true) as any;
  } catch (error: any) {
    handleRestError(error);
    return (defVal || null) as any;
  }
}

export async function DELETE_AND_THROW<T = any>(url: string, params?: {[key: string]: any}): Promise<T> {
  try {
    const response = await client.delete<T>(url, { params });
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw new RestError(error.response.statusText, error.response.status, error.response.data);
    }

    throw error;
  }
}

export async function DELETE<T = any>(url: string, params?: {[key: string]: any}, defVal?: T): Promise<T> {
  try {
    const res = await DELETE_AND_THROW(url, params);
    if (res) {
      return res;
    }
    return (defVal || true) as any;
  } catch (error: any) {
    handleRestError(error);
    return (defVal || null) as any;
  }
}

export async function PATCH_AND_THROW<T = any>(url: string, data?: any, params?: {[key: string]: any}): Promise<T> {
  try {
    const response = await client.patch<T>(url, data, { params });
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw new RestError(error.response.statusText, error.response.status, error.response.data);
    }

    throw error;
  }
}

export async function PATCH<T = any>(url: string, data?: any, params?: {[key: string]: any}, defVal?: T): Promise<T> {
  try {
    const res = await PATCH_AND_THROW<T>(url, data, params);
    if (res) {
      return res;
    }
    return (defVal || true) as any;
  } catch (error: any) {
    handleRestError(error);
    return (defVal || null) as any;
  }
}

export async function GET_FILE(url_: string, fileName: string): Promise<boolean> {
  try {
    const response = await client.get(url_, { responseType: 'blob' });
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', fileName);
    document.body.appendChild(link);
    link.click();
    return true;
  } catch (error: any) {
    handleRestError(error);

    return false;
  }
}

export interface ErrorResponse {
  error: string;
}

export interface SuccessResponse {
  message: string;
}

export interface PagniatedResponse<T> {
  total: number;
  items: T[];
}

export class RestError<T = ErrorResponse> extends Error {
  constructor(message: string, public status: number, public data: T) {
    super(message);
  }
}