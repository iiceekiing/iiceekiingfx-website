'use client';

import { useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { 
  BarChart3, 
  BookOpen, 
  TrendingUp, 
  FileText, 
  Signal, 
  Calculator,
  User,
  LogOut
} from 'lucide-react';
import { Button } from '@/components/ui/Button';

const menuItems = [
  { href: '/dashboard', label: 'Overview', icon: BarChart3 },
  { href: '/portfolio', label: 'Portfolio', icon: TrendingUp },
  { href: '/journal', label: 'Trade Journal', icon: FileText },
  { href: '/courses', label: 'My Courses', icon: BookOpen },
  { href: '/signals', label: 'Signals', icon: Signal },
  { href: '/calculator', label: 'Calculator', icon: Calculator },
  { href: '/profile', label: 'Profile', icon: User },
];

export function Sidebar() {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const pathname = usePathname();

  const handleLogout = () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    window.location.href = '/login';
  };

  return (
    <div className={`bg-gray-800 border-r border-gray-700 transition-all duration-300 ${
      isCollapsed ? 'w-16' : 'w-64'
    }`}>
      <div className="flex flex-col h-full">
        {/* Logo */}
        <div className="p-4 border-b border-gray-700">
          <div className="flex items-center space-x-3">
            <div className="w-8 h-8 bg-gradient-to-r from-blue-500 to-amber-500 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">FX</span>
            </div>
            {!isCollapsed && (
              <span className="text-white font-semibold">iiceekiingfx</span>
            )}
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 p-4 space-y-2">
          {menuItems.map((item) => {
            const Icon = item.icon;
            const isActive = pathname === item.href;
            
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`flex items-center space-x-3 px-3 py-2 rounded-lg transition-colors ${
                  isActive
                    ? 'bg-blue-600 text-white'
                    : 'text-gray-300 hover:bg-gray-700 hover:text-white'
                }`}
              >
                <Icon className="w-5 h-5 flex-shrink-0" />
                {!isCollapsed && <span>{item.label}</span>}
              </Link>
            );
          })}
        </nav>

        {/* User Actions */}
        <div className="p-4 border-t border-gray-700">
          <Button
            variant="ghost"
            size="sm"
            onClick={handleLogout}
            className="w-full justify-start"
          >
            <LogOut className="w-4 h-4 mr-2" />
            {!isCollapsed && 'Logout'}
          </Button>
        </div>

        {/* Collapse Toggle */}
        <button
          onClick={() => setIsCollapsed(!isCollapsed)}
          className="absolute -right-3 top-8 w-6 h-6 bg-gray-600 rounded-full flex items-center justify-center text-gray-300 hover:bg-gray-500"
        >
          <span className="text-xs">
            {isCollapsed ? '→' : '←'}
          </span>
        </button>
      </div>
    </div>
  );
}
