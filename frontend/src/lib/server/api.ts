import { env } from '$env/dynamic/public';
import { getRequestEvent } from '$app/server';
import { error } from '@sveltejs/kit';

const API_BASE = env.PUBLIC_API_BASE_URL || 'http://localhost:8080';

export class ApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.status = status;
	}
}

type ApiFetchOptions = {
	method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
	body?: unknown;
	formData?: FormData;
	auth?: boolean;
	searchParams?: Record<string, string | number | undefined>;
};

export function getToken() {
	return getRequestEvent().cookies.get('access_token');
}

export async function apiFetch<T>(path: string, options: ApiFetchOptions = {}): Promise<T> {
	const event = getRequestEvent();
	const url = new URL(`${API_BASE}/api/v1${path}`);

	for (const [key, value] of Object.entries(options.searchParams ?? {})) {
		if (value !== undefined && value !== '') url.searchParams.set(key, String(value));
	}

	const headers: Record<string, string> = {};

	if (options.auth) {
		const token = event.cookies.get('access_token');
		if (!token) throw error(401, 'Sessão expirada. Inicie sessão novamente.');
		headers.Authorization = `Bearer ${token}`;
	}

	let body: BodyInit | undefined;
	if (options.formData) {
		body = options.formData;
	} else if (options.body !== undefined) {
		headers['Content-Type'] = 'application/json';
		body = JSON.stringify(options.body);
	}

	let response: Response;
	try {
		response = await event.fetch(url, { method: options.method ?? 'GET', headers, body });
	} catch {
		throw new ApiError(503, 'Falha de ligação ao servidor.');
	}

	if (response.status === 204) return undefined as T;

	const data = await response.json().catch(() => ({}));

	if (!response.ok) {
		throw new ApiError(response.status, data.error ?? 'Ocorreu um erro inesperado.');
	}

	return data as T;
}
