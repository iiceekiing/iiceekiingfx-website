'use client';

import { useState, useEffect } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Signal } from '@/lib/auth';
import api from '@/lib/api';
import { 
  TrendingUp, 
  TrendingDown, 
  Filter,
  Search,
  Clock,
  Target
} from 'lucide-react';

export default function SignalsPage() {
  const [signals, setSignals] = useState<Signal[]>([]);
  const [filteredSignals, setFilteredSignals] = useState<Signal[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedPair, setSelectedPair] = useState('all');
  const [selectedAction, setSelectedAction] = useState('all');

  useEffect(() => {
    fetchSignals();
  }, []);

  useEffect(() => {
    filterSignals();
  }, [signals, searchTerm, selectedPair, selectedAction]);

  const fetchSignals = async () => {
    try {
      const response = await api.get('/signals/');
      setSignals(response.data.signals || []);
    } catch (error) {
      console.error('Failed to fetch signals:', error);
    } finally {
      setLoading(false);
    }
  };

  const filterSignals = () => {
    let filtered = signals;

    // Filter by search term
    if (searchTerm) {
      filtered = filtered.filter(signal =>
        signal.pair.toLowerCase().includes(searchTerm.toLowerCase()) ||
        signal.notes.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Filter by pair
    if (selectedPair !== 'all') {
      filtered = filtered.filter(signal => signal.pair === selectedPair);
    }

    // Filter by action
    if (selectedAction !== 'all') {
      filtered = filtered.filter(signal => signal.action === selectedAction);
    }

    setFilteredSignals(filtered);
  };

  const uniquePairs = Array.from(new Set(signals.map(s => s.pair)));

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
          <h1 className="text-3xl font-bold text-white mb-2">Trading Signals</h1>
          <p className="text-gray-400">Professional trading signals with risk management</p>
        </div>

        {/* Filters */}
        <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <Input
                placeholder="Search signals..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10"
              />
            </div>
            
            <select
              value={selectedPair}
              onChange={(e) => setSelectedPair(e.target.value)}
              className="px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="all">All Pairs</option>
              {uniquePairs.map(pair => (
                <option key={pair} value={pair}>{pair}</option>
              ))}
            </select>

            <select
              value={selectedAction}
              onChange={(e) => setSelectedAction(e.target.value)}
              className="px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="all">All Actions</option>
              <option value="BUY">Buy Only</option>
              <option value="SELL">Sell Only</option>
            </select>

            <Button variant="outline" onClick={() => {
              setSearchTerm('');
              setSelectedPair('all');
              setSelectedAction('all');
            }}>
              <Filter className="w-4 h-4 mr-2" />
              Reset Filters
            </Button>
          </div>
        </div>

        {/* Signals List */}
        <div className="space-y-4">
          {filteredSignals.map((signal) => (
            <div key={signal.id} className="bg-gray-800 p-6 rounded-xl border border-gray-700">
              <div className="flex justify-between items-start mb-4">
                <div className="flex items-center space-x-4">
                  <div className={`p-3 rounded-lg ${
                    signal.action === 'BUY' 
                      ? 'bg-green-500/20' 
                      : 'bg-red-500/20'
                  }`}>
                    {signal.action === 'BUY' ? (
                      <TrendingUp className="w-6 h-6 text-green-400" />
                    ) : (
                      <TrendingDown className="w-6 h-6 text-red-400" />
                    )}
                  </div>
                  <div>
                    <h3 className="text-xl font-semibold text-white">{signal.pair}</h3>
                    <p className="text-gray-400">
                      <span className={`px-2 py-1 rounded text-xs ${
                        signal.action === 'BUY' 
                          ? 'bg-green-500/20 text-green-400' 
                          : 'bg-red-500/20 text-red-400'
                      }`}>
                        {signal.action}
                      </span>
                    </p>
                  </div>
                </div>

                <div className="text-right">
                  <span className={`px-2 py-1 rounded-full text-xs ${
                    signal.status === 'active' 
                      ? 'bg-green-500/20 text-green-400' 
                      : 'bg-gray-600/20 text-gray-400'
                  }`}>
                    {signal.status}
                  </span>
                  <p className="text-sm text-gray-400 mt-1">
                    {new Date(signal.created_at).toLocaleString()}
                  </p>
                </div>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                <div>
                  <p className="text-sm text-gray-400">Entry Price</p>
                  <p className="text-lg font-semibold text-white">
                    ${signal.entry_price}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-gray-400">Stop Loss</p>
                  <p className="text-lg font-semibold text-red-400">
                    ${signal.stop_loss}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-gray-400">Take Profit</p>
                  <p className="text-lg font-semibold text-green-400">
                    ${signal.take_profit}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-gray-400">Risk/Reward</p>
                  <p className="text-lg font-semibold text-blue-400">
                    {signal.risk_reward}:1
                  </p>
                </div>
              </div>

              {signal.notes && (
                <div className="bg-gray-700 p-3 rounded-lg">
                  <p className="text-sm text-gray-300">{signal.notes}</p>
                </div>
              )}

              <div className="flex items-center justify-between mt-4">
                <div className="flex items-center text-sm text-gray-400">
                  <Target className="w-4 h-4 mr-1" />
                  <span>Potential: ${Math.abs(signal.take_profit - signal.entry_price).toFixed(2)}</span>
                </div>
                <div className="flex items-center text-sm text-gray-400">
                  <Clock className="w-4 h-4 mr-1" />
                  <span>Created {new Date(signal.created_at).toLocaleDateString()}</span>
                </div>
              </div>
            </div>
          ))}

          {filteredSignals.length === 0 && (
            <div className="text-center py-12">
              <p className="text-gray-400">
                {searchTerm || selectedPair !== 'all' || selectedAction !== 'all'
                  ? 'No signals found matching your filters'
                  : 'No signals available at the moment'
                }
              </p>
            </div>
          )}
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-gray-800 p-6 rounded-xl border border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-400 mb-1">Total Signals</p>
                <p className="text-2xl font-bold text-white">{filteredSignals.length}</p>
              </div>
              <div className="p-3 bg-blue-500/20 rounded-lg">
                <Target className="w-6 h-6 text-blue-400" />
              </div>
            </div>
          </div>

          <div className="bg-gray-800 p-6 rounded-xl border border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-400 mb-1">Buy Signals</p>
                <p className="text-2xl font-bold text-green-400">
                  {filteredSignals.filter(s => s.action === 'BUY').length}
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
                <p className="text-sm text-gray-400 mb-1">Sell Signals</p>
                <p className="text-2xl font-bold text-red-400">
                  {filteredSignals.filter(s => s.action === 'SELL').length}
                </p>
              </div>
              <div className="p-3 bg-red-500/20 rounded-lg">
                <TrendingDown className="w-6 h-6 text-red-400" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
