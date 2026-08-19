import type { NextRequest, NextResponse } from "next/server";

import { cookieKey } from "./env";

const week = 60 * 60 * 24 * 7;

const base = {
  httpOnly: true,
  sameSite: "lax" as const,
  path: "/",
  secure: false,
};

export { cookieKey };

export function setOn(res: NextResponse, name: string, value: string, maxAge = week) {
  res.cookies.set(cookieKey(name), value, { ...base, maxAge });
}

export function readRequestCookie(req: NextRequest, name: string): string | undefined {
  return req.cookies.get(cookieKey(name))?.value;
}
