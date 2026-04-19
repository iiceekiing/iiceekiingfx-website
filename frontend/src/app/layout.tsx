import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "iiceekiingfx - Forex Trading Platform",
  description: "Premium Forex education, analytics, and trading signals platform",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${inter.className} antialiased`}>
      <body className="bg-gray-900 text-white min-h-screen">{children}</body>
    </html>
  );
}
