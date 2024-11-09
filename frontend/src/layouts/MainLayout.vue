<script setup lang="ts">
import Mail from '../components/Mail.vue'
import { mockAccounts, mockMails } from '../data/mock'
</script>

<template>
  <div class="h-screen">
    <Mail
      :accounts="mockAccounts"
      :mails="mockMails"
      :nav-collapsed-size="4"
    />
  </div>
</template>
