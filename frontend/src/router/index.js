import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import Login from '../views/Login.vue';
import KeySpecList from '../views/KeySpecList.vue';
import KeyManagement from '../views/KeyManagement.vue';
import Config from '../views/Config.vue';

const routes = [
  {
    path: '/',
    redirect: '/key-specs',
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false },
  },
  {
    path: '/key-specs',
    name: 'KeySpecList',
    component: KeySpecList,
    meta: { requiresAuth: true },
  },
  {
    path: '/keys/:specId',
    name: 'KeyManagement',
    component: KeyManagement,
    meta: { requiresAuth: true },
  },
  {
    path: '/config',
    name: 'Config',
    component: Config,
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const requiresAuth = to.meta.requiresAuth !== false;

  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login');
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/key-specs');
  } else {
    next();
  }
});

export default router;
