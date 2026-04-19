'use client';

import { useState, useEffect } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { TradingAccount } from '@/lib/auth';
import api from '@/lib/api';
import { 
  TrendingUp, 
  TrendingDown, 
  DollarSign, 
  Plus,
  RefreshCw,
  Eye,
  EyeOff
} from 'lucide-react';

export default function PortfolioPage() {
  const [accounts, setAccounts] = useState<TradingAccount[]>([]);
  const [loading, setLoading] = useState(true);
  const [showConnectForm, setShowConnectForm] = useState(false);
  const [showCredentials, setShowCredentials] = useState<{[key: string]: boolean}>({});
  const [formData, setFormData] = useState({
    broker_name: '',
    account_id: '',
    login: '',
    password: '',
    server: '',
    account_type: 'Standard',
    currency: 'USD',
    leverage: 100,
  });

  useEffect(() => {
    fetchAccounts();
  }, []);

  const fetchAccounts = async () => {
    try {
      const response = await api.get('/portfolio/accounts');
      setAccounts(response.data.accounts || []);
    } catch (error) {
      console.error('Failed to fetch accounts:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleConnectAccount = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.post('/portfolio/connect', formData);
      setShowConnectForm(false);
      setFormData({
        broker_name: '',
        account_id: '',
        login: '',
        password: '',
        server: '',
        account_type: 'Standard',
        currency: 'USD',
        leverage: 100,
      });
      fetchAccounts();
    } catch (error: any) {
      console.error('Failed to connect account:', error);
    }
  };

  const toggleCredentials = (accountId: string) => {
    setShowCredentials(prev => ({
      ...prev,
      [accountId]: !prev[accountId]
    }));
  };

  const handleSync = async (accountId: string) => {
    try {
      await api.post(`/portfolio/accounts/${accountId}/sync`);
      fetchAccounts();
    } catch (error) {
      console.error('Failed to sync account:', error);
    }
  };

  const totalBalance = accounts.reduce((sum, account) => sum + account.balance, 0);
  const totalEquity = accounts.reduce((sum, account) => sum + account.equity, 0);

  if (loading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2">Portfolio</h1>
            <p className="text-gray-400">Manage your trading accounts and track performance</p>
          </div>
          <Button onClick={() => setShowConnectForm(true)}>
            <Plus className="w-4 h-4 mr-2" />
            Connect Account
          </Button>
        </div>

        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-gray-800 p-6 rounded-xl border border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-400 mb-1">Total Balance</p>
                <p className="text-2xl font-bold text-white">
                  ${totalBalance.toLocaleString()}
                </p>
              </div>
              <div className="p-3 bg-blue-500/20 rounded-lg">
                <DollarSign className="w-6 h-6 text-blue-400" />
              </div>
            </div>
          </div>

          <div className="bg-gray-800 p-6 rounded-xl border border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-400 mb-1">Total Equity</p>
                <p className="text-2xl font-bold text-white">
                  ${totalEquity.toLocaleString()}
                </p>
              </div>
              <div className="p-3 bg-green-500/20 rounded-lg">
                <TrendingUp className="w-6 h-6 text-green-400" />
              </div>
            </div>
          </div>

          <div className="bg-gray-800 p-6 rounded-xl border border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-400 mb-1">Active Accounts</p>
                <p className="text-2xl font-bold text-white">
                  {accounts.filter(acc => acc.is_active).length}
                </p>
              </div>
              <div className="p-3 bg-gray-600/20 rounded-lg">
                <RefreshCw className="w-6 h-6 text-gray-400" />
              </div>
            </div>
          </div>
        </div>

        {/* Connect Account Modal */}
        {showConnectForm && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-gray-800 p-6 rounded-xl w-full max-w-md">
              <h3 className="text-xl font-semibold text-white mb-4">Connect Trading Account</h3>
              <form onSubmit={handleConnectAccount} className="space-y-4">
                <Input
                  label="Broker Name"
                  value={formData.broker_name}
                  onChange={(e) => setFormData({...formData, broker_name: e.target.value})}
                  placeholder="e.g., MetaQuotes"
                  required
                />
                <Input
                  label="Account ID"
                  value={formData.account_id}
                  onChange={(e) => setFormData({...formData, account_id: e.target.value})}
                  placeholder="Your trading account ID"
                  required
                />
                <Input
                  label="Login"
                  value={formData.login}
                  onChange={(e) => setFormData({...formData, login: e.target.value})}
                  placeholder="Trading account login"
                  required
                />
                <Input
                  label="Password"
                  type="password"
                  value={formData.password}
                  onChange={(e) => setFormData({...formData, password: e.target.value})}
                  placeholder="Trading account password"
                  required
                />
                <Input
                  label="Server"
                  value={formData.server}
                  onChange={(e) => setFormData({...formData, server: e.target.value})}
                  placeholder="e.g., MetaQuotes-Demo"
                  required
                />
                <div className="grid grid-cols-2 gap-4">
                  <Input
                    label="Account Type"
                    value={formData.account_type}
                    onChange={(e) => setFormData({...formData, account_type: e.target.value})}
                  />
                  <Input
                    label="Currency"
                    value={formData.currency}
                    onChange={(e) => setFormData({...formData, currency: e.target.value})}
                  />
                </div>
                <Input
                  label="Leverage"
                  type="number"
                  value={formData.leverage}
                  onChange={(e) => setFormData({...formData, leverage: parseInt(e.target.value)})}
                />
                <div className="flex gap-4">
                  <Button type="submit" className="flex-1">
                    Connect Account
                  </Button>
                  <Button 
                    type="button" 
                    variant="outline" 
                    onClick={() => setShowConnectForm(false)}
                    className="flex-1"
                  >
                    Cancel
                  </Button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* Accounts List */}
        <div className="space-y-4">
          {accounts.map((account) => (
            <div key={account.id} className="bg-gray-800 p-6 rounded-xl border border-gray-700">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h3 className="text-lg font-semibold text-white">{account.broker_name}</h3>
                  <p className="text-gray-400">Account: {account.account_id}</p>
                  <p className="text-sm text-gray-500">
                    {account.account_type} · {account.currency} · {account.leverage}:1
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  <span className={`px-2 py-1 rounded-full text-xs ${
                    account.is_active 
                      ? 'bg-green-500/20 text-green-400' 
                      : 'bg-gray-600/20 text-gray-400'
                  }`}>
                    {account.is_active ? 'Active' : 'Inactive'}
                  </span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleSync(account.id)}
                  >
                    <RefreshCw className="w-4 h-4" />
                  </Button>
                </div>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                <div>
                  <p className="text-sm text-gray-400">Balance</p>
                  <p className="text-lg font-semibold text-white">
                    ${account.balance.toLocaleString()}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-gray-400">Equity</p>
                  <p className="text-lg font-semibold text-white">
                    ${account.equity.toLocaleString()}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-gray-400">Margin</p>
                  <p className="text-lg font-semibold text-white">
                    ${account.margin.toLocaleString()}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-gray-400">Free Margin</p>
                  <p className="text-lg font-semibold text-white">
                    ${account.free_margin.toLocaleString()}
                  </p>
                </div>
              </div>

              <div className="flex justify-between items-center">
                <p className="text-sm text-gray-400">
                  Last sync: {new Date(account.last_sync).toLocaleString()}
                </p>
              </div>
            </div>
          ))}

          {accounts.length === 0 && (
            <div className="text-center py-12">
              <p className="text-gray-400 mb-4">No trading accounts connected yet</p>
              <Button onClick={() => setShowConnectForm(true)}>
                <Plus className="w-4 h-4 mr-2" />
                Connect Your First Account
              </Button>
            </div>
          )}
        </div>
      </div>
    </DashboardLayout>
  );
}
