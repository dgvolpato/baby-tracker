

export const index = 0;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_layout.svelte.js')).default;
export const universal = {
  "ssr": false,
  "prerender": false
};
export const universal_id = "src/routes/+layout.ts";
export const imports = ["_app/immutable/nodes/0.B5OX3KbE.js","_app/immutable/chunks/BSeaJ4Cw.js","_app/immutable/chunks/BNHIpwgD.js","_app/immutable/chunks/CQRjT9mT.js","_app/immutable/chunks/sQ4771Yp.js"];
export const stylesheets = ["_app/immutable/assets/0.P_JkNZKn.css"];
export const fonts = [];
