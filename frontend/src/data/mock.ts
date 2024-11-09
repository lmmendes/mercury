import type { Account, Mail } from '../types';

export const mockAccounts: Account[] = [
  {
    id: '1',
    name: 'Personal',
    email: 'personal@example.com',
    icon: 'lucide:mail',
  },
  {
    id: '2',
    name: 'Work',
    email: 'work@example.com',
    icon: 'lucide:briefcase',
  },
];

export const mockMails: Mail[] = [
  {
    id: '1',
    name: 'John Doe',
    email: 'john@example.com',
    subject: 'Project Update',
    text: "Hi, I wanted to give you a quick update on the project. Everything is going well and we're on track to meet our deadlines...",
    date: new Date(2024, 11, 9).toISOString(),
    read: false,
    labels: ['work'],
  },
  {
    id: '2',
    name: 'Alice Smith',
    email: 'alice@example.com',
    subject: 'Weekend Plans',
    text: 'Hey! Are you free this weekend? I was thinking we could grab coffee and catch up...',
    date: new Date(2024, 11, 9).toISOString(),
    read: true,
    labels: ['personal'],
  },
  {
    id: '3',
    name: 'Bob Wilson',
    email: 'bob@example.com',
    subject: 'Meeting Notes',
    text: "Here are the notes from today's meeting. Please review and let me know if I missed anything important...",
    date: new Date(2024, 11, 9).toISOString(),
    read: false,
    labels: ['work', 'important'],
  },
];
