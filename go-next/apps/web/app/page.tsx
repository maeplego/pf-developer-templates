async function loadProducts(): Promise<{ sku: string; name: string; priceMinor: number; currency: string }[]> {
  const base = process.env.API_URL?.trim() || "http://localhost:8080";
  try {
    const res = await fetch(`${base}/v1/products`, { cache: "no-store" });
    if (!res.ok) return [];
    const body = (await res.json()) as {
      products?: { sku: string; name: string; priceMinor: number; currency: string }[];
    };
    return body.products ?? [];
  } catch {
    return [];
  }
}

export default async function HomePage() {
  const products = await loadProducts();
  return (
    <main>
      <h1>{{PROJECT}}</h1>
      <p>Go API + Next.js from P04 workspace / P06 commerce patterns. OIDC login is a stub at /login.</p>
      <p>
        <a href="/health">health</a> · <a href="/ready">ready</a> · <a href="/login">login</a>
      </p>
      <ul>
        {products.map((p) => (
          <li key={p.sku}>
            {p.name} ({p.sku}) {p.priceMinor} {p.currency}
          </li>
        ))}
      </ul>
    </main>
  );
}
