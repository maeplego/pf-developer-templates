import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "{{PROJECT}}",
  description: "Generated from P04/P06 patterns via pf-dev new",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
