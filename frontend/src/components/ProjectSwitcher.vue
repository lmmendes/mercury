<script setup lang="ts">
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { cn } from '@/lib/utils'
import { Icon } from '@iconify/vue'
import { computed, ref } from 'vue'
import type { Project } from '../types'

interface Props {
  isCollapsed: boolean
  projects: Project[]
}

const props = defineProps<Props>()

const selectedEmail = ref<string>(props.projects[0].email)
const selectedProject = computed(() => props.projects.find(item => item.email === selectedEmail.value))
</script>

<template>
  <Select v-model="selectedEmail">
    <SelectTrigger
      aria-label="Select project"
      :class="cn(
        'flex items-center gap-2 [&>span]:line-clamp-1 [&>span]:flex [&>span]:w-full [&>span]:items-center [&>span]:gap-1 [&>span]:truncate [&_svg]:h-4 [&_svg]:w-4 [&_svg]:shrink-0',
        { 'flex h-9 w-9 shrink-0 items-center justify-center p-0 [&>span]:w-auto [&>svg]:hidden': isCollapsed },
      )"
    >
      <SelectValue placeholder="Select a project">
        <div class="flex items-center gap-3">
          <Icon class="size-4" :icon="selectedProject!.icon" />
          <span v-if="!isCollapsed">
            {{ selectedProject!.name }}
          </span>
        </div>
      </SelectValue>
    </SelectTrigger>
    <SelectContent>
      <SelectItem v-for="project of projects" :key="project.email" :value="project.email">
        <div class="flex items-center gap-3 [&_svg]:size-4 [&_svg]:shrink-0 [&_svg]:text-foreground">
          <Icon class="size-4" :icon="project.icon" />
          {{ project.email }}
        </div>
      </SelectItem>
    </SelectContent>
  </Select>
</template>
