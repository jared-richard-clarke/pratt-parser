import { assertSnapshot } from "https://deno.land/std@0.199.0/testing/snapshot.ts";
import { parse } from "../parser.js";

// Snapshots to quickly catch unexpected changes to parser output.

Deno.test("snapshot", async function (t) {
    const x = parse("1 + 2 * 3");
    await assertSnapshot(t, x);
})

Deno.test("snapshot", async function (t) {
    const x = parse("0.1 + 0.2");
    await assertSnapshot(t, x);
});

Deno.test("snapshot", async function(t) {
    const x = parse("1 / 0");
    await assertSnapshot(t, x);
});

Deno.test("snapshot", async function(t) {
    const x = parse("2 ^ 3 ^ 4");
    await assertSnapshot(t, x);
});

Deno.test("snapshot", async function(t) {
    const x = parse("5 ^ 2.5");
    await assertSnapshot(t, x);
});
