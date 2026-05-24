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

export const api = {
	login: (username: string, password: string) =>
		request<{ token: string }>('POST', '/login', { username, password }),

	feedings: {
		list: () => request<Feeding[]>('GET', '/feedings'),
		create: (data: FeedingInput) => request<Feeding>('POST', '/feedings', data),
		update: (id: string, data: FeedingInput) => request<Feeding>('PUT', `/feedings/${id}`, data),
		delete: (id: string) => request<void>('DELETE', `/feedings/${id}`)
	}
};
