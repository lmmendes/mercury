<script setup lang="ts">
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { cn } from '@/lib/utils'
import { Icon } from '@iconify/vue'
import { computed, ref } from 'vue'
import type { Account } from '../types'

interface Props {
  isCollapsed: boolean
  accounts: Account[]
}

const props = defineProps<Props>()

const selectedEmail = ref<string>(props.accounts[0].email)
const selectedAccount = computed(() => props.accounts.find(item => item.email === selectedEmail.value))
</script>

<template>
  <Select v-model="selectedEmail">
    <SelectTrigger
      aria-label="Select account"
      :class="cn(
        'flex items-center gap-2 [&>span]:line-clamp-1 [&>span]:flex [&>span]:w-full [&>span]:items-center [&>span]:gap-1 [&>span]:truncate [&_svg]:h-4 [&_svg]:w-4 [&_svg]:shrink-0',
        { 'flex h-9 w-9 shrink-0 items-center justify-center p-0 [&>span]:w-auto [&>svg]:hidden': isCollapsed },
      )"
    >
      <SelectValue placeholder="Select an account">
        <div class="flex items-center gap-3">
          <Icon class="size-4" :icon="selectedAccount!.icon" />
          <span v-if="!isCollapsed">
            {{ selectedAccount!.name }}
          </span>
        </div>
      </SelectValue>
    </SelectTrigger>
    <SelectContent>
      <SelectItem v-for="account of accounts" :key="account.email" :value="account.email">
        <div class="flex items-center gap-3 [&_svg]:size-4 [&_svg]:shrink-0 [&_svg]:text-foreground">
          <Icon class="size-4" :icon="account.icon" />
          {{ account.email }}
        </div>
      </SelectItem>
    </SelectContent>
  </Select>
</template>
