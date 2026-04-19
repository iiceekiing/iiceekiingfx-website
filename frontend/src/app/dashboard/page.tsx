'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { Button } from '@/components/ui/Button';
import { DashboardOverview } from '@/lib/auth';
import api from '@/lib/api';
import { 
  TrendingUp, 
  TrendingDown, 
  DollarSign, 
  BarChart3,
  Activity
} from 'lucide-react';

export default function DashboardPage() {
  const [overview, setOverview] = useState<DashboardOverview | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchOverview = async () => {
      try {
        const response = await api.get<DashboardOverview>('/dashboard/overview');
        setOverview(response.data);
      } catch (error) {
        console.error('Failed to fetch dashboard overview:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchOverview();
  }, []);

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
        <div>
          <h1 className="text-3xl font-bold text-white mb-2">Dashboard Overview</h1>
          <p className="text-gray-400">Welcome back! Here's your trading performance summary.</p>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div className="bg-gray-800 p-6 rounded-xl border border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-400 mb-1">Total Profit</p>
                <p className="text-2xl font-bold {(overview?.total_profit ?? 0) >= 0 ? 'text-green-400' : 'text-red-400'}">
                  ${(overview?.total_profit ?? 0).toLocaleString()}
                </p>
              </div>
              <div className={`p-3 rounded-lg ${(overview?.total_profit ?? 0) >= 0 ? 'bg-green-500/20' : 'bg-red-500/20'}`}>
                {(overview?.total_profit ?? 0) >= 0 ? (
                  <TrendingUp className="w-6 h-6 text-green-400" />
                ) : (
                  <TrendingDown className="w-6 h-6 text-red-400" />
                )}
              </div>
            </div>
          </div>

          <div className="bg-gray-800 p-6 rounded-xl border border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-400 mb-1">Win Rate</p>
                <p className="text-2xl font-bold text-blue-400">
                  {overview?.win_rate?.toFixed(1) || '0.0'}%
                </p>
              </div>
              <div className="p-3 bg-blue-500/20 rounded-lg">
                <BarChart3 className="w-6 h-6 text-blue-400" />
              </div>
            </div>
          </div>

          <div className="bg-gray-800 p-6 rounded-xl border border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-400 mb-1">Active Accounts</p>
                <p className="text-2xl font-bold text-white">
                  {overview?.active_accounts || 0}
                </p>
              </div>
              <div className="p-3 bg-gray-600/20 rounded-lg">
                <Activity className="w-6 h-6 text-gray-400" />
              </div>
            </div>
          </div>

          <div className="bg-gray-800 p-6 rounded-xl border border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-400 mb-1">Total Trades</p>
                <p className="text-2xl font-bold text-white">
                  {overview?.total_trades || 0}
                </p>
              </div>
              <div className="p-3 bg-gray-600/20 rounded-lg">
                <DollarSign className="w-6 h-6 text-gray-400" />
              </div>
            </div>
          </div>
        </div>

        {/* Balance & Equity */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="bg-gray-800 p-6 rounded-xl border border-gray-700">
            <h3 className="text-lg font-semibold text-white mb-4">Current Balance</h3>
            <p className="text-3xl font-bold text-green-400">
              ${overview?.current_balance?.toLocaleString() || '0'}
            </p>
          </div>

          <div className="bg-gray-800 p-6 rounded-xl border border-gray-700">
            <h3 className="text-lg font-semibold text-white mb-4">Current Equity</h3>
            <p className="text-3xl font-bold text-blue-400">
              ${overview?.current_equity?.toLocaleString() || '0'}
            </p>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="flex flex-wrap gap-4">
          <Link href="/portfolio">
            <Button>View Portfolio</Button>
          </Link>
          <Link href="/journal">
            <Button variant="outline">Add Trade</Button>
          </Link>
          <Link href="/signals">
            <Button variant="secondary">View Signals</Button>
          </Link>
        </div>
      </div>
    </DashboardLayout>
  );
}
