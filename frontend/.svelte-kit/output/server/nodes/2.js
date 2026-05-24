

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/2.DdZF4TF7.js","_app/immutable/chunks/BSeaJ4Cw.js","_app/immutable/chunks/BNHIpwgD.js","_app/immutable/chunks/B2sYrIGw.js"];
export const stylesheets = [];
export const fonts = [];
