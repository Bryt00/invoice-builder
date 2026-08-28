export interface Client {
  id: string
  user_id?: string
  name: string
  email: string
  company?: string
  phone?: string
  tax_id?: string
  address?: string
  created_at?: string
  updated_at?: string
}

export interface LineItem {
  id?: string
  invoice_id?: string
  description: string
  quantity: number
  unit_price: number
  amount: number
}

export type InvoiceStatus = 'draft' | 'sent' | 'paid' | 'overdue' | 'cancelled'

export interface Invoice {
  id: string
  user_id: string
  client_id: string
  client?: Client
  invoice_number: string
  status: InvoiceStatus
  issue_date: string
  due_date: string
  currency: string
  subtotal: number
  tax_rate: number
  tax: number
  discount: number
  total: number
  notes?: string
  public_token?: string
  items: LineItem[]
  created_at: string
  updated_at: string
}

export interface Receipt {
  id: string
  invoice_id: string
  receipt_number: string
  amount_paid: number
  payment_method: string
  created_at: string
}
