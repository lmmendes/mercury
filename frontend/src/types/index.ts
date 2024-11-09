export interface Account {
  id: string
  name: string
  email: string
  icon: string
}

export interface Inbox {
  id: string
  accountId: string
  email: string
}

export interface Mail {
  id: string
  name: string
  email: string
  subject: string
  text: string
  date: string
  read: boolean
  labels: string[]
}
