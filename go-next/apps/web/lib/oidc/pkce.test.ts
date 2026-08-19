import assert from "node:assert/strict";
import test from "node:test";

import { s256 } from "./pkce.ts";

test("s256 is 43-char base64url without padding", () => {
  const d = s256("demo-verifier");
  assert.equal(d.length, 43);
  assert.match(d, /^[A-Za-z0-9_-]+$/);
  assert.notEqual(d, s256("other"));
});
