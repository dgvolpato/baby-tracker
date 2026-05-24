const BASE = import.meta.env.VITE_API_URL ?? '';

function getToken(): string {
	return localStorage.getItem('token') ?? '';
}

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
	const res = await fetch(`${BASE}${path}`, {
		method,
		headers: {
			'Content-Type': 'application/json',
			Authorization: `Bearer ${getToken()}`
		},
		body: body !== undefined ? JSON.stringify(body) : undefined
	});

	if (res.status === 401) {
		localStorage.removeItem('token');
		window.location.href = '/login';
		throw new Error('unauthorized');
	}

	if (res.status === 204) return undefined as T;

	if (!res.ok) {
		const err = await res.json().catch(() => ({}));
		throw new Error((err as { error?: string }).error ?? 'request failed');
	}

	return res.json();
}

export interface Feeding {
	id: string;
	timestamp: string;
	type: 'formula' | 'breast';
	oz: number;
	createdBy?: string;
}

export type FeedingInput = Omit<Feeding, 'id' | 'createdBy'>;

export interface Diaper {
	id: string;
	timestamp: string;
	type: 'wet' | 'poop' | 'both';
	createdBy?: string;
}

export type DiaperInput = Omit<Diaper, 'id' | 'createdBy'>;

export const api = {
	login: (username: string, password: string) =>
		request<{ token: string }>('POST', '/login', { username, password }),

	feedings: {
		list: () => request<Feeding[]>('GET', '/feedings'),
		create: (data: FeedingInput) => request<Feeding>('POST', '/feedings', data),
		update: (id: string, data: FeedingInput) => request<Feeding>('PUT', `/feedings/${id}`, data),
		delete: (id: string) => request<void>('DELETE', `/feedings/${id}`)
	},

	diapers: {
		list: () => request<Diaper[]>('GET', '/diapers'),
		create: (data: DiaperInput) => request<Diaper>('POST', '/diapers', data),
		update: (id: string, data: DiaperInput) => request<Diaper>('PUT', `/diapers/${id}`, data),
		delete: (id: string) => request<void>('DELETE', `/diapers/${id}`)
	},

	measurements: {
		list: () => request<Measurement[]>('GET', '/measurements'),
		create: (data: MeasurementInput) => request<Measurement>('POST', '/measurements', data),
		update: (id: string, data: MeasurementInput) =>
			request<Measurement>('PUT', `/measurements/${id}`, data),
		delete: (id: string) => request<void>('DELETE', `/measurements/${id}`)
	}
};

export interface Measurement {
	id: string;
	timestamp: string;
	weightLbs?: number;
	heightIn?: number;
	createdBy?: string;
}

export type MeasurementInput = Omit<Measurement, 'id' | 'createdBy'>;
