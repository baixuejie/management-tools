import axios from 'axios';

const apiClient = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// POST /api/login → {token}
export const authAPI = {
  login: (username, password, rememberMe = false) =>
    apiClient.post('/login', { username, password, remember_me: rememberMe }),
};

// Key specs CRUD
export const keySpecAPI = {
  list: () => apiClient.get('/key-specs'),
  get: (id) => apiClient.get(`/key-specs/${id}`),
  create: (data) => apiClient.post('/key-specs', data),
  update: (id, data) => apiClient.put(`/key-specs/${id}`, data),
  delete: (id) => apiClient.delete(`/key-specs/${id}`),
};

// Keys management
export const keyAPI = {
  // GET /api/keys?spec_id=X&only_unused=true&limit=20&offset=0 → {keys, total}
  list: (specId, { onlyUnused = false, limit = 50, offset = 0 } = {}) =>
    apiClient.get('/keys', {
      params: { spec_id: specId, only_unused: onlyUnused, limit, offset },
    }),
  // POST /api/keys/batch → {message}
  batchUpload: (specId, keys) =>
    apiClient.post('/keys/batch', { spec_id: specId, keys }),
  // GET /api/keys/available/:spec_id → {id, spec_id, key_value, is_used, used_at}
  getAvailable: (specId) => apiClient.get(`/keys/available/${specId}`),
  // PUT /api/keys/:id/use → {message}
  markUsed: (id) => apiClient.put(`/keys/${id}/use`),
  // DELETE /api/keys/:id
  delete: (id) => apiClient.delete(`/keys/${id}`),
};

// Config
export const configAPI = {
  // GET /api/config/copy-template → {template}
  getTemplate: () => apiClient.get('/config/copy-template'),
  // PUT /api/config/copy-template → {message}
  updateTemplate: (template) =>
    apiClient.put('/config/copy-template', { template }),
};

export default apiClient;
