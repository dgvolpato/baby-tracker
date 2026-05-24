function roundToNearest30(date) {
  const ms = 30 * 60 * 1e3;
  return new Date(Math.round(date.getTime() / ms) * ms);
}
function toLocalInputDate(date) {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, "0");
  const d = String(date.getDate()).padStart(2, "0");
  return `${y}-${m}-${d}`;
}
function roundedTimeValue(date) {
  const rounded = roundToNearest30(date);
  const h = String(rounded.getHours()).padStart(2, "0");
  const m = String(rounded.getMinutes()).padStart(2, "0");
  return `${h}:${m}`;
}
function timeOptions() {
  const options = [];
  for (let h = 0; h < 24; h++) {
    for (const m of [0, 30]) {
      const hh = String(h).padStart(2, "0");
      const mm = String(m).padStart(2, "0");
      const period = h < 12 ? "AM" : "PM";
      const displayH = h === 0 ? 12 : h > 12 ? h - 12 : h;
      options.push({ value: `${hh}:${mm}`, label: `${displayH}:${mm} ${period}` });
    }
  }
  return options;
}
function dayLabel(dateStr) {
  const today = toLocalInputDate(/* @__PURE__ */ new Date());
  const yesterday = toLocalInputDate(new Date(Date.now() - 864e5));
  if (dateStr === today) return "Today";
  if (dateStr === yesterday) return "Yesterday";
  const d = /* @__PURE__ */ new Date(dateStr + "T12:00:00");
  return d.toLocaleDateString("en-US", { month: "short", day: "numeric" });
}
export {
  toLocalInputDate as a,
  dayLabel as d,
  roundedTimeValue as r,
  timeOptions as t
};
