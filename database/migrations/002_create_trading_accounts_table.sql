-- Create trading_accounts table
CREATE TABLE IF NOT EXISTS trading_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    broker_name VARCHAR(100) NOT NULL,
    account_id VARCHAR(100) NOT NULL,
    account_type VARCHAR(50) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    leverage INTEGER NOT NULL DEFAULT 100,
    balance DECIMAL(15,2) DEFAULT 0.00,
    equity DECIMAL(15,2) DEFAULT 0.00,
    margin DECIMAL(15,2) DEFAULT 0.00,
    free_margin DECIMAL(15,2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT true,
    last_sync TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, account_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_trading_accounts_user_id ON trading_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_trading_accounts_is_active ON trading_accounts(is_active);
CREATE INDEX IF NOT EXISTS idx_trading_accounts_last_sync ON trading_accounts(last_sync);
