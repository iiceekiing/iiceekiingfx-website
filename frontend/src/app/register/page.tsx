'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Mail, Lock, User } from 'lucide-react';
import { RegisterRequest, AuthResponse } from '@/lib/auth';
import api from '@/lib/api';

export default function RegisterPage() {
  const [formData, setFormData] = useState<RegisterRequest>({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    // Basic client-side validation
    if (!formData.email || !formData.email.includes('@')) {
      setError('Please enter a valid email address.');
      setLoading(false);
      return;
    }

    if (formData.password.length < 6) {
      setError('Password must be at least 6 characters long.');
      setLoading(false);
      return;
    }

    if (!formData.first_name || !formData.last_name) {
      setError('Please enter both first and last name.');
      setLoading(false);
      return;
    }

    try {
      const response = await api.post<AuthResponse>('/auth/register', formData);
      const { access_token, refresh_token } = response.data;
      
      localStorage.setItem('access_token', access_token);
      localStorage.setItem('refresh_token', refresh_token);
      
      router.push('/dashboard');
    } catch (err: any) {
      let errorMessage = 'Registration failed';
      
      if (err.response?.data?.error) {
        if (err.response.data.error.includes('duplicate key') || err.response.data.error.includes('users_pkey')) {
          errorMessage = 'An account with this email already exists. Please sign in instead.';
        } else if (err.response.data.error.includes('password')) {
          errorMessage = 'Password is too weak. Please choose a stronger password.';
        } else if (err.response.data.error.includes('email')) {
          errorMessage = 'Please enter a valid email address.';
        } else {
          errorMessage = err.response.data.error;
        }
      }
      
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center fintech-gradient">
      <div className="glass-effect p-8 rounded-2xl w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-white mb-2">iiceekiingfx</h1>
          <p className="text-gray-400">Create your trading account</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid grid-cols-2 gap-4">
            <Input
              type="text"
              label="First Name"
              value={formData.first_name}
              onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
              icon={<User className="w-4 h-4" />}
              placeholder="Enter first name"
              required
            />
            <Input
              type="text"
              label="Last Name"
              value={formData.last_name}
              onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
              icon={<User className="w-4 h-4" />}
              placeholder="Enter last name"
              required
            />
          </div>

          <Input
            type="email"
            label="Email"
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            icon={<Mail className="w-4 h-4" />}
            placeholder="Enter your email"
            required
          />

          <Input
            type="password"
            label="Password"
            value={formData.password}
            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            icon={<Lock className="w-4 h-4" />}
            placeholder="Create a password"
            required
          />

          {error && (
            <div className="bg-red-500/20 border border-red-500 text-red-400 p-3 rounded-lg">
              {error}
              {error.includes('already exists') && (
                <div className="mt-2">
                  <a href="/login" className="text-blue-400 hover:text-blue-300 underline">
                    Sign in to your existing account
                  </a>
                </div>
              )}
            </div>
          )}

          <Button
            type="submit"
            className="w-full"
            loading={loading}
            disabled={loading}
          >
            {loading ? 'Creating Account...' : 'Create Account'}
          </Button>
        </form>

        <div className="text-center mt-6">
          <p className="text-gray-400">
            Already have an account?{' '}
            <a href="/login" className="text-blue-400 hover:text-blue-300 transition-colors">
              Sign in
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}
