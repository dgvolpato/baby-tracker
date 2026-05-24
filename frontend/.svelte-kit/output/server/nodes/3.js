

export const index = 3;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/dashboard/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/3.B9K9OPkp.js","_app/immutable/chunks/BSeaJ4Cw.js","_app/immutable/chunks/BNHIpwgD.js","_app/immutable/chunks/BESZb65m.js","_app/immutable/chunks/DZM41zAk.js"];
export const stylesheets = ["_app/immutable/assets/3.DmUPhoiY.css"];
export const fonts = [];
