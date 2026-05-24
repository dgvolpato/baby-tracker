export function roundToNearest30(date: Date): Date {
	const ms = 30 * 60 * 1000;
	return new Date(Math.round(date.getTime() / ms) * ms);
}

export function toLocalInputDate(date: Date): string {
	const y = date.getFullYear();
	const m = String(date.getMonth() + 1).padStart(2, '0');
	const d = String(date.getDate()).padStart(2, '0');
	return `${y}-${m}-${d}`;
}

export function roundedTimeValue(date: Date): string {
	const rounded = roundToNearest30(date);
	const h = String(rounded.getHours()).padStart(2, '0');
	const m = String(rounded.getMinutes()).padStart(2, '0');
	return `${h}:${m}`;
}

export function timeOptions(): { value: string; label: string }[] {
	const options = [];
	for (let h = 0; h < 24; h++) {
		for (const m of [0, 30]) {
			const hh = String(h).padStart(2, '0');
			const mm = String(m).padStart(2, '0');
			const period = h < 12 ? 'AM' : 'PM';
			const displayH = h === 0 ? 12 : h > 12 ? h - 12 : h;
			options.push({ value: `${hh}:${mm}`, label: `${displayH}:${mm} ${period}` });
		}
	}
	return options;
}

export function formatDateTime(isoStr: string): string {
	const d = new Date(isoStr);
	const date = d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
	const time = d.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' });
	return `${date}, ${time}`;
}

export function dayLabel(dateStr: string): string {
	const today = toLocalInputDate(new Date());
	const yesterday = toLocalInputDate(new Date(Date.now() - 86400000));
	if (dateStr === today) return 'Today';
	if (dateStr === yesterday) return 'Yesterday';
	const d = new Date(dateStr + 'T12:00:00');
	return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
}

export function fmtOz(oz: number): string {
	return oz % 1 === 0 ? `${oz} oz` : `${oz.toFixed(1)} oz`;
}
