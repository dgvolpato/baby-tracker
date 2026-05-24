

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/2.DpOB3gBc.js","_app/immutable/chunks/BSeaJ4Cw.js","_app/immutable/chunks/BNHIpwgD.js","_app/immutable/chunks/sQ4771Yp.js"];
export const stylesheets = [];
export const fonts = [];
