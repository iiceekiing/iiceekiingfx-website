-- Create trade_journal table
CREATE TABLE IF NOT EXISTS trade_journal (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pair VARCHAR(20) NOT NULL,
    entry_price DECIMAL(15,5) NOT NULL,
    exit_price DECIMAL(15,5) NOT NULL,
    lot_size DECIMAL(10,2) NOT NULL,
    result VARCHAR(20) NOT NULL CHECK (result IN ('win', 'loss', 'breakeven')),
    r_multiple DECIMAL(10,2) NOT NULL,
    notes TEXT,
    strategy_tag VARCHAR(100),
    trade_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create signals table
CREATE TABLE IF NOT EXISTS signals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pair VARCHAR(20) NOT NULL,
    action VARCHAR(10) NOT NULL CHECK (action IN ('BUY', 'SELL')),
    entry_price DECIMAL(15,5) NOT NULL,
    stop_loss DECIMAL(15,5) NOT NULL,
    take_profit DECIMAL(15,5) NOT NULL,
    risk_reward DECIMAL(5,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'closed', 'cancelled')),
    notes TEXT,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_trade_journal_user_id ON trade_journal(user_id);
CREATE INDEX IF NOT EXISTS idx_trade_journal_pair ON trade_journal(pair);
CREATE INDEX IF NOT EXISTS idx_trade_journal_trade_date ON trade_journal(trade_date);
CREATE INDEX IF NOT EXISTS idx_trade_journal_result ON trade_journal(result);
CREATE INDEX IF NOT EXISTS idx_signals_pair ON signals(pair);
CREATE INDEX IF NOT EXISTS idx_signals_status ON signals(status);
CREATE INDEX IF NOT EXISTS idx_signals_created_by ON signals(created_by);
