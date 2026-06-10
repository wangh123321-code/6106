const API_BASE = '/api/v1';

function getToken() {
  if (typeof window === 'undefined') return '';
  return localStorage.getItem('token') || '';
}

async function request(method, path, data = null) {
  const headers = { 'Content-Type': 'application/json' };
  const token = getToken();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const opts = { method, headers };
  if (data && method !== 'GET') {
    opts.body = JSON.stringify(data);
  }

  const res = await fetch(`${API_BASE}${path}`, opts);
  const json = await res.json();

  if (res.status === 401) {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
  }

  return json;
}

export const api = {
  get: (path) => request('GET', path),
  post: (path, data) => request('POST', path, data),
  put: (path, data) => request('PUT', path, data),
  del: (path) => request('DELETE', path),
};

export function isLoggedIn() {
  if (typeof window === 'undefined') return false;
  return !!localStorage.getItem('token');
}

export function getCurrentUser() {
  if (typeof window === 'undefined') return null;
  const raw = localStorage.getItem('user');
  if (!raw) return null;
  try { return JSON.parse(raw); } catch { return null; }
}

export function logout() {
  if (typeof window === 'undefined') return;
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  window.location.href = '/login';
}
