import { createRouter, createWebHistory } from 'vue-router'
import {logout} from "@/rest/api/auth.ts";
import {destroySession, isSessionExisting} from "@/rest/session.ts";
import type {User} from "@/@types/user.ts";
import {useMainStore} from "@/stores/main.ts";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '',
      name: 'home',
      component: () => import('@/views/HomeView.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue')
    },
    {
      path: '/characters/:id',
      name: 'character-view',
      component: () => import('@/views/CharacterView.vue')
    },
    {
      path: '/characters/:id/edit',
      name: 'character-edit',
      component: () => import('@/views/CharacterEditorView.vue')
    },
    {
      path: '/logout',
      name: 'logout',
      component: () => null!,
    }
  ],
});

router.beforeEach(async (to, _, next) => {
  if (to.path === "/login") {
    next();
    return;
  }

  const store = useMainStore();

  if (to.path === "/logout") {
    await _logout(store);
    next("/login");
    return;
  }

  if (!store.user && isSessionExisting()) {
    if (!await store.loadUser()) {
      await _logout(store);
      next("/login?redirect=" + btoa(to.fullPath));
      return;
    }
  }

  if (!store.user) {
    next("/login?redirect=" + btoa(to.fullPath));
    return;
  }

  next();
});

async function _logout(store: any) {
  await logout();
  destroySession();

  store.user = null as unknown as User;

  window.location.reload();
}

export default router
