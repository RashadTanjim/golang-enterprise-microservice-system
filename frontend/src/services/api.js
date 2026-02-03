const API_BASE = import.meta.env.VITE_API_BASE || '/api/v1';
const HEALTH_BASE = import.meta.env.VITE_HEALTH_BASE || '/health';

const defaultHeaders = {
  'Content-Type': 'application/json'
};

export async function fetchHealth(service) {
  const response = await fetch(`${HEALTH_BASE}/${service}`);
  if (!response.ok) {
    throw new Error('Health check failed');
  }
  return response.json();
}

export async function listUsers(token) {
  const response = await fetch(`${API_BASE}/users`, {
    headers: {
      ...defaultHeaders,
      Authorization: `Bearer ${token}`
    }
  });
  if (!response.ok) {
    throw new Error('Failed to list users');
  }
  return response.json();
}

export async function listOrders(token) {
  const response = await fetch(`${API_BASE}/orders`, {
    headers: {
      ...defaultHeaders,
      Authorization: `Bearer ${token}`
    }
  });
  if (!response.ok) {
    throw new Error('Failed to list orders');
  }
  return response.json();
}
