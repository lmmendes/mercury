<script setup lang="ts">
import Mail from '../components/Mail.vue'
import { mockProjects, mockMails } from '../data/mock'
</script>

<template>
  <div class="h-screen">
    <Mail
      :projects="mockProjects"
      :mails="mockMails"
      :nav-collapsed-size="4"
    />
  </div>
</template>
