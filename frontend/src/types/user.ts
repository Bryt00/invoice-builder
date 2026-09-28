export interface Role {
  id: string
  name: string
  description?: string
}

export interface User {
  id: string
  name: string
  email: string
  role_id: string
  role: Role
  is_activated: boolean
  is_profile_complete: boolean
  created_at: string
  updated_at: string
}

export interface BusinessProfile {
  id?: string
  user_id?: string
  company_name: string
  company_email: string
  phone: string
  address: string
  logo_url: string
  tax_id: string
  default_currency: string
  bank_name: string
  account_number: string
  account_name: string
  routing_code: string
}

export interface AuthTokenPair {
  access_token: string
  refresh_token: string
}

export interface LoginResponse {
  status: string
  user: User
  profile?: BusinessProfile
  credits: number
  tokens: AuthTokenPair
}
