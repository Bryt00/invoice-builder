export interface Payment {
  id: string
  user_id: string
  invoice_id?: string
  amount: number
  currency: string
  payment_method: string
  transaction_ref: string
  reference?: string
  status: string
  created_at: string
  updated_at: string
}

export interface CreditPackage {
  id: string
  name: string
  credits: number
  price: number
  currency: string
  is_active: boolean
  created_at?: string
}

export interface CreditTxn {
  id: string
  user_id: string
  amount: number
  type: string
  description: string
  created_at: string
}

export interface AuditLog {
  id: string
  user_id?: string
  action: string
  ip_address: string
  user_agent: string
  created_at: string
}

export interface WebhookLog {
  id: string
  event_type: string
  payload: string
  status: string
  created_at: string
}
