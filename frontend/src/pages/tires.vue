<script setup lang="ts">
import TireManager from '@/components/TireManager.vue';
import { useAuth0 } from '@auth0/auth0-vue';
import { useTitle } from '@vueuse/core';

useTitle('タイヤ');
const { isLoading, isAuthenticated } = useAuth0();
</script>

<template>
  <div class="w-full" v-if="isLoading">
    <div class="w-full loading loading-spinner loading-xl" />
  </div>
  <div v-else-if="!isAuthenticated">ログインしてください</div>
  <div v-else>
    <h1 class="font-bold text-lg mt-8">タイヤ</h1>
    <TireManager />
  </div>
</template>
