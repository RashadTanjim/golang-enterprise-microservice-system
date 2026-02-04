const API_BASE = import.meta.env.VITE_API_BASE || '/api/v1';
const HEALTH_BASE = import.meta.env.VITE_HEALTH_BASE || '/health';

const defaultHeaders = {
  'Content-Type': 'application/json'
};

const allowedHealthServices = new Set(['user', 'order', 'audit-log']);

export class ApiError extends Error {
  constructor(message, status = 0, code = 'API_ERROR', details = null) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
    this.details = details;
  }
}

const isNonEmptyString = (value) => typeof value === 'string' && value.trim().length > 0;

const ensureToken = (token) => {
  if (!isNonEmptyString(token)) {
    throw new ApiError('Authorization token is required.', 0, 'VALIDATION');
  }
};

const ensureHealthService = (service) => {
  if (!allowedHealthServices.has(service)) {
    throw new ApiError('Unknown health service target.', 0, 'VALIDATION');
  }
};

const parseJson = async (response) => {
  const contentType = response.headers?.get?.('content-type') || '';
  if (!contentType.includes('application/json')) {
    return null;
  }
  try {
    return await response.json();
  } catch {
    return null;
  }
};

const request = async (path, options = {}) => {
  const response = await fetch(path, options);
  const payload = await parseJson(response);

  if (!response.ok) {
    const message =
      payload?.error?.message ||
      payload?.message ||
      `Request failed with status ${response.status}.`;
    const code = payload?.error?.code || 'HTTP_ERROR';
    throw new ApiError(message, response.status, code, payload);
  }

  if (payload && payload.success === false) {
    const message = payload?.error?.message || 'Request failed.';
    const code = payload?.error?.code || 'API_ERROR';
    throw new ApiError(message, response.status, code, payload);
  }

  return payload;
};

export async function fetchHealth(service) {
  ensureHealthService(service);
  return request(`${HEALTH_BASE}/${service}`);
}

export async function listUsers(token) {
  ensureToken(token);
  const payload = await request(`${API_BASE}/users`, {
    headers: {
      ...defaultHeaders,
      Authorization: `Bearer ${token}`
    }
  });

  if (!payload || !Array.isArray(payload.data)) {
    throw new ApiError('Invalid users response payload.', 0, 'INVALID_RESPONSE', payload);
  }

  return payload;
}

export async function listOrders(token) {
  ensureToken(token);
  const payload = await request(`${API_BASE}/orders`, {
    headers: {
      ...defaultHeaders,
      Authorization: `Bearer ${token}`
    }
  });

  if (!payload || !Array.isArray(payload.data)) {
    throw new ApiError('Invalid orders response payload.', 0, 'INVALID_RESPONSE', payload);
  }

  return payload;
}
