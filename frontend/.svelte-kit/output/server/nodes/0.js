

export const index = 0;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_layout.svelte.js')).default;
export const universal = {
  "ssr": false,
  "prerender": false
};
export const universal_id = "src/routes/+layout.ts";
export const imports = ["_app/immutable/nodes/0.Dt71I2n_.js","_app/immutable/chunks/BSeaJ4Cw.js","_app/immutable/chunks/BNHIpwgD.js","_app/immutable/chunks/BUx1JyEs.js","_app/immutable/chunks/B2sYrIGw.js"];
export const stylesheets = ["_app/immutable/assets/0.P_JkNZKn.css"];
export const fonts = [];
