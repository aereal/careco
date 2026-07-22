<script setup lang="ts">
import BottomNav from '@/components/BottomNav.vue';
import LoginButton from '@/components/LoginButton.vue';
import LogoutButton from '@/components/LogoutButton.vue';
import { useAuth0 } from '@auth0/auth0-vue';
import { RouterView } from 'vue-router';

const { isLoading, isAuthenticated } = useAuth0();
</script>

<template>
  <header
    class="app-header navbar bg-primary text-primary-content fixed top-0 inset-x-0 z-10 p-0"
  >
    <div
      class="max-w-2xl mx-auto w-full px-4 flex items-center justify-between"
    >
      <span class="font-bold text-lg">Careco</span>
      <div v-if="!isLoading">
        <LogoutButton v-if="isAuthenticated" />
        <LoginButton v-else />
      </div>
    </div>
  </header>
  <main class="app-content">
    <div class="max-w-2xl mx-auto p-4">
      <RouterView />
    </div>
  </main>
  <BottomNav />
</template>

<style scoped></style>
