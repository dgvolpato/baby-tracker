

export const index = 0;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_layout.svelte.js')).default;
export const universal = {
  "ssr": false,
  "prerender": false
};
export const universal_id = "src/routes/+layout.ts";
export const imports = ["_app/immutable/nodes/0.JjhjDPJC.js","_app/immutable/chunks/BSeaJ4Cw.js","_app/immutable/chunks/BNHIpwgD.js","_app/immutable/chunks/BFEHRNIQ.js","_app/immutable/chunks/nUxS9Bqo.js"];
export const stylesheets = ["_app/immutable/assets/0.g93I6oGV.css"];
export const fonts = [];
