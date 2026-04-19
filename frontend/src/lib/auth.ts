export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: 'student' | 'admin';
  created_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}

export interface Course {
  id: string;
  title: string;
  description: string;
  price: number;
  image_url: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Signal {
  id: string;
  pair: string;
  action: 'BUY' | 'SELL';
  entry_price: number;
  stop_loss: number;
  take_profit: number;
  risk_reward: number;
  status: 'active' | 'closed';
  notes: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface TradingAccount {
  id: string;
  user_id: string;
  broker_name: string;
  account_id: string;
  account_type: string;
  currency: string;
  leverage: number;
  balance: number;
  equity: number;
  margin: number;
  free_margin: number;
  is_active: boolean;
  last_sync: string;
  created_at: string;
  updated_at: string;
}

export interface DashboardOverview {
  total_profit: number;
  win_rate: number;
  active_accounts: number;
  total_trades: number;
  current_balance: number;
  current_equity: number;
}
