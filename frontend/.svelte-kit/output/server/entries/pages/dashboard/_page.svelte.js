import { c as create_ssr_component } from "../../../chunks/ssr.js";
import { a as toLocalInputDate, d as dayLabel } from "../../../chunks/utils2.js";
const css = {
  code: "h1.svelte-1m8f6ip{font-size:22px;margin-bottom:20px}.error-text.svelte-1m8f6ip{color:var(--danger);font-size:14px}.list.svelte-1m8f6ip{display:flex;flex-direction:column;gap:10px}.day-card.svelte-1m8f6ip{display:flex;flex-direction:column;gap:8px}.day-header.svelte-1m8f6ip{display:flex;justify-content:space-between;align-items:baseline}.day-label.svelte-1m8f6ip{font-weight:600;font-size:16px}.day-total.svelte-1m8f6ip{font-size:20px;font-weight:700;color:var(--accent)}.day-meta.svelte-1m8f6ip{display:flex;align-items:center;gap:8px;flex-wrap:wrap}.tag.svelte-1m8f6ip{font-size:12px;padding:2px 8px;border-radius:12px;font-weight:500}.formula.svelte-1m8f6ip{background:#e8f0fe;color:#1a56db}.breast.svelte-1m8f6ip{background:#fde8f0;color:#b5477a}.count.svelte-1m8f6ip{font-size:12px;color:var(--text-muted);margin-left:auto}",
  map: `{"version":3,"file":"+page.svelte","sources":["+page.svelte"],"sourcesContent":["<script lang=\\"ts\\">import { onMount } from \\"svelte\\";\\nimport { api } from \\"$lib/api\\";\\nimport { dayLabel, fmtOz, toLocalInputDate } from \\"$lib/utils\\";\\nlet feedings = [];\\nlet loading = true;\\nlet error = \\"\\";\\n$: dayTotals = computeTotals(feedings);\\nfunction computeTotals(list) {\\n  const map = /* @__PURE__ */ new Map();\\n  for (const f of list) {\\n    const key = toLocalInputDate(new Date(f.timestamp));\\n    if (!map.has(key)) {\\n      map.set(key, { date: key, label: dayLabel(key), total: 0, formula: 0, breast: 0, count: 0 });\\n    }\\n    const day = map.get(key);\\n    day.total = Math.round((day.total + f.oz) * 10) / 10;\\n    day.count++;\\n    if (f.type === \\"formula\\") day.formula = Math.round((day.formula + f.oz) * 10) / 10;\\n    else day.breast = Math.round((day.breast + f.oz) * 10) / 10;\\n  }\\n  return Array.from(map.values()).sort((a, b) => b.date.localeCompare(a.date));\\n}\\nonMount(async () => {\\n  try {\\n    feedings = await api.feedings.list();\\n  } catch {\\n    error = \\"Failed to load feedings.\\";\\n  } finally {\\n    loading = false;\\n  }\\n});\\n<\/script>\\n\\n<div class=\\"page\\">\\n\\t<h1>Dashboard</h1>\\n\\n\\t{#if loading}\\n\\t\\t<p class=\\"muted\\">Loading…</p>\\n\\t{:else if error}\\n\\t\\t<p class=\\"error-text\\">{error}</p>\\n\\t{:else if dayTotals.length === 0}\\n\\t\\t<p class=\\"muted\\">No feedings logged yet. <a href=\\"/history\\">Add one →</a></p>\\n\\t{:else}\\n\\t\\t<div class=\\"list\\">\\n\\t\\t\\t{#each dayTotals as day (day.date)}\\n\\t\\t\\t\\t<div class=\\"day-card card\\">\\n\\t\\t\\t\\t\\t<div class=\\"day-header\\">\\n\\t\\t\\t\\t\\t\\t<span class=\\"day-label\\">{day.label}</span>\\n\\t\\t\\t\\t\\t\\t<span class=\\"day-total\\">{fmtOz(day.total)}</span>\\n\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t<div class=\\"day-meta\\">\\n\\t\\t\\t\\t\\t\\t{#if day.formula > 0}\\n\\t\\t\\t\\t\\t\\t\\t<span class=\\"tag formula\\">Formula {fmtOz(day.formula)}</span>\\n\\t\\t\\t\\t\\t\\t{/if}\\n\\t\\t\\t\\t\\t\\t{#if day.breast > 0}\\n\\t\\t\\t\\t\\t\\t\\t<span class=\\"tag breast\\">Breast {fmtOz(day.breast)}</span>\\n\\t\\t\\t\\t\\t\\t{/if}\\n\\t\\t\\t\\t\\t\\t<span class=\\"count\\">{day.count} feeding{day.count !== 1 ? 's' : ''}</span>\\n\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t</div>\\n\\t\\t\\t{/each}\\n\\t\\t</div>\\n\\t{/if}\\n</div>\\n\\n<style>\\n\\th1 {\\n\\t\\tfont-size: 22px;\\n\\t\\tmargin-bottom: 20px;\\n\\t}\\n\\t.error-text {\\n\\t\\tcolor: var(--danger);\\n\\t\\tfont-size: 14px;\\n\\t}\\n\\t.list {\\n\\t\\tdisplay: flex;\\n\\t\\tflex-direction: column;\\n\\t\\tgap: 10px;\\n\\t}\\n\\t.day-card {\\n\\t\\tdisplay: flex;\\n\\t\\tflex-direction: column;\\n\\t\\tgap: 8px;\\n\\t}\\n\\t.day-header {\\n\\t\\tdisplay: flex;\\n\\t\\tjustify-content: space-between;\\n\\t\\talign-items: baseline;\\n\\t}\\n\\t.day-label {\\n\\t\\tfont-weight: 600;\\n\\t\\tfont-size: 16px;\\n\\t}\\n\\t.day-total {\\n\\t\\tfont-size: 20px;\\n\\t\\tfont-weight: 700;\\n\\t\\tcolor: var(--accent);\\n\\t}\\n\\t.day-meta {\\n\\t\\tdisplay: flex;\\n\\t\\talign-items: center;\\n\\t\\tgap: 8px;\\n\\t\\tflex-wrap: wrap;\\n\\t}\\n\\t.tag {\\n\\t\\tfont-size: 12px;\\n\\t\\tpadding: 2px 8px;\\n\\t\\tborder-radius: 12px;\\n\\t\\tfont-weight: 500;\\n\\t}\\n\\t.formula {\\n\\t\\tbackground: #e8f0fe;\\n\\t\\tcolor: #1a56db;\\n\\t}\\n\\t.breast {\\n\\t\\tbackground: #fde8f0;\\n\\t\\tcolor: #b5477a;\\n\\t}\\n\\t.count {\\n\\t\\tfont-size: 12px;\\n\\t\\tcolor: var(--text-muted);\\n\\t\\tmargin-left: auto;\\n\\t}\\n</style>\\n"],"names":[],"mappings":"AAkEC,iBAAG,CACF,SAAS,CAAE,IAAI,CACf,aAAa,CAAE,IAChB,CACA,0BAAY,CACX,KAAK,CAAE,IAAI,QAAQ,CAAC,CACpB,SAAS,CAAE,IACZ,CACA,oBAAM,CACL,OAAO,CAAE,IAAI,CACb,cAAc,CAAE,MAAM,CACtB,GAAG,CAAE,IACN,CACA,wBAAU,CACT,OAAO,CAAE,IAAI,CACb,cAAc,CAAE,MAAM,CACtB,GAAG,CAAE,GACN,CACA,0BAAY,CACX,OAAO,CAAE,IAAI,CACb,eAAe,CAAE,aAAa,CAC9B,WAAW,CAAE,QACd,CACA,yBAAW,CACV,WAAW,CAAE,GAAG,CAChB,SAAS,CAAE,IACZ,CACA,yBAAW,CACV,SAAS,CAAE,IAAI,CACf,WAAW,CAAE,GAAG,CAChB,KAAK,CAAE,IAAI,QAAQ,CACpB,CACA,wBAAU,CACT,OAAO,CAAE,IAAI,CACb,WAAW,CAAE,MAAM,CACnB,GAAG,CAAE,GAAG,CACR,SAAS,CAAE,IACZ,CACA,mBAAK,CACJ,SAAS,CAAE,IAAI,CACf,OAAO,CAAE,GAAG,CAAC,GAAG,CAChB,aAAa,CAAE,IAAI,CACnB,WAAW,CAAE,GACd,CACA,uBAAS,CACR,UAAU,CAAE,OAAO,CACnB,KAAK,CAAE,OACR,CACA,sBAAQ,CACP,UAAU,CAAE,OAAO,CACnB,KAAK,CAAE,OACR,CACA,qBAAO,CACN,SAAS,CAAE,IAAI,CACf,KAAK,CAAE,IAAI,YAAY,CAAC,CACxB,WAAW,CAAE,IACd"}`
};
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let feedings = [];
  function computeTotals(list) {
    const map = /* @__PURE__ */ new Map();
    for (const f of list) {
      const key = toLocalInputDate(new Date(f.timestamp));
      if (!map.has(key)) {
        map.set(key, {
          date: key,
          label: dayLabel(key),
          total: 0,
          formula: 0,
          breast: 0,
          count: 0
        });
      }
      const day = map.get(key);
      day.total = Math.round((day.total + f.oz) * 10) / 10;
      day.count++;
      if (f.type === "formula") day.formula = Math.round((day.formula + f.oz) * 10) / 10;
      else day.breast = Math.round((day.breast + f.oz) * 10) / 10;
    }
    return Array.from(map.values()).sort((a, b) => b.date.localeCompare(a.date));
  }
  $$result.css.add(css);
  computeTotals(feedings);
  return `<div class="page"><h1 class="svelte-1m8f6ip" data-svelte-h="svelte-101alym">Dashboard</h1> ${`<p class="muted" data-svelte-h="svelte-jbr5m6">Loading…</p>`} </div>`;
});
export {
  Page as default
};
