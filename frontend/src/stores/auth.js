import { defineStore } from 'pinia';
import { authAPI } from '../api';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    username: localStorage.getItem('username') || null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
  },

  actions: {
    async login(username, password, rememberMe = false) {
      const response = await authAPI.login(username, password, rememberMe);
      this.token = response.data.token;
      this.username = username;
      localStorage.setItem('token', this.token);
      localStorage.setItem('username', username);
    },

    logout() {
      this.token = null;
      this.username = null;
      localStorage.removeItem('token');
      localStorage.removeItem('username');
    },
  },
});
