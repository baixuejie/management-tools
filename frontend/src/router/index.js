import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';

import Login from '../views/Login.vue';
import Ledger from '../views/Ledger.vue';
import KeyManagement from '../views/KeyManagement.vue';
import More from '../views/More.vue';
import CostRecords from '../views/CostRecords.vue';
import Customers from '../views/Customers.vue';
import KeySpecList from '../views/KeySpecList.vue';
import Config from '../views/Config.vue';
import Profile from '../views/Profile.vue';

const routes = [
  { path: '/', redirect: '/ledger' },
  { path: '/login', name: 'Login', component: Login, meta: { requiresAuth: false } },

  { path: '/ledger', name: 'Ledger', component: Ledger, meta: { requiresAuth: true } },
  { path: '/keys', name: 'Keys', component: KeyManagement, meta: { requiresAuth: true } },
  { path: '/more', name: 'More', component: More, meta: { requiresAuth: true } },
  { path: '/more/costs', name: 'CostRecords', component: CostRecords, meta: { requiresAuth: true } },
  { path: '/more/customers', name: 'Customers', component: Customers, meta: { requiresAuth: true } },
  { path: '/more/specs', name: 'KeySpecList', component: KeySpecList, meta: { requiresAuth: true } },
  { path: '/more/template', name: 'TemplateConfig', component: Config, meta: { requiresAuth: true } },
  { path: '/more/profile', name: 'Profile', component: Profile, meta: { requiresAuth: true } },

  // Legacy routes kept for compatibility.
  { path: '/key-specs', redirect: '/more/specs' },
  { path: '/config', redirect: '/more/template' },
  {
    path: '/keys/:specId',
    redirect: (to) => `/keys?specId=${to.params.specId}`,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore();
  const requiresAuth = to.meta.requiresAuth !== false;

  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login');
    return;
  }

  if (to.path === '/login' && authStore.isAuthenticated) {
    next('/ledger');
    return;
  }

  if (requiresAuth && authStore.isAuthenticated && !authStore.user) {
    try {
      await authStore.fetchMe();
    } catch {
      authStore.logout();
      next('/login');
      return;
    }
  }

  next();
});

export default router;
