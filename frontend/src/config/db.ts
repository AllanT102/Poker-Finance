const API_URL = process.env.API_BASE_URL || ''
const API_VERSION = process.env.API_VERSION || ''
const API_PORT = process.env.API_BASE_PORT || 0

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

interface FetchOptions {
  endpoint: string;
  method: HttpMethod;
  data?: any;
}

interface ApiResponse<T> {
  status: number;
  data: T;
}

async function fetchBackend<T>({
  endpoint,
  method,
  data,
}: FetchOptions): Promise<ApiResponse<T>> {
  const headers: HeadersInit = {
    Accept: 'application/json',
    'Content-Type': 'application/json'
  };

  const options: RequestInit = {
    method,
    headers,
  };

  if (data) {
    options.body = JSON.stringify(data);
  }

  const response = await fetch(`${API_URL}:${API_PORT}${API_VERSION}${endpoint}`, options);
  const responseData = await response.json();

  if (!response.ok) {
    throw new Error(responseData.message || 'API request failed');
  }

  return {
    status: response.status,
    data: responseData
  };
}

export { fetchBackend };
