<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue';

const auth0 = useAuth0();
const logout = async (): Promise<void> =>
  auth0.logout({ logoutParams: { returnTo: window.location.href } });
</script>
<template>
  <button @click="logout" class="btn btn-ghost">Logout</button>
</template>
