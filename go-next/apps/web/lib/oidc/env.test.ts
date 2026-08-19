import assert from "node:assert/strict";
import test from "node:test";

import { cookieKey, oidcEnabled, postLogoutRedirectUri } from "./env.ts";

test("oidcEnabled is false without issuer and client", () => {
  delete process.env.OIDC_ISSUER;
  delete process.env.OIDC_CLIENT_ID;
  assert.equal(oidcEnabled(), false);
  assert.equal(cookieKey("rp_state"), "rp_state");
});

test("oidcEnabled and cookieKey when configured", () => {
  process.env.OIDC_ISSUER = "http://localhost:8081";
  process.env.OIDC_CLIENT_ID = "demo-web";
  process.env.OIDC_REDIRECT_URI = "http://localhost:3000/callback";
  assert.equal(oidcEnabled(), true);
  assert.equal(cookieKey("rp_state"), "rp_state_demo_web");
  assert.equal(postLogoutRedirectUri(), "http://localhost:3000/logged-out");
});
