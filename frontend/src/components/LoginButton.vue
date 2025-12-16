<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue';

const auth0 = useAuth0();
const login = async (): Promise<void> => auth0.loginWithRedirect();
</script>
<template>
  <button @click="login" className="btn btn-primary">Login</button>
</template>
