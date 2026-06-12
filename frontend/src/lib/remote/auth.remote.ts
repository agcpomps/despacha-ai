import * as v from 'valibot';
import { form, query, getRequestEvent } from '$app/server';
import { redirect, invalid } from '@sveltejs/kit';
import { apiFetch, ApiError } from '$lib/server/api';
import { isValidAngolaPhone, normalizeAngolaPhone } from '$lib/utils';
import type { AuthResponse, User } from '$lib/types';

const phoneSchema = v.pipe(
	v.string(),
	v.check(isValidAngolaPhone, 'Insira um número válido de Angola. Ex: 923 456 789.')
);

function safeNext(next: string | undefined) {
	if (!next || !next.startsWith('/') || next.startsWith('//')) return '/';
	return next;
}

function setSessionCookie(token: string) {
	const event = getRequestEvent();
	event.cookies.set('access_token', token, {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		secure: event.url.protocol === 'https:',
		maxAge: 60 * 60 * 24
	});
}

export const getCurrentUser = query(async (): Promise<User | null> => {
	const { cookies } = getRequestEvent();
	if (!cookies.get('access_token')) return null;

	try {
		return await apiFetch<User>('/me', { auth: true });
	} catch (err) {
		if (err instanceof ApiError && (err.status === 401 || err.status === 404)) {
			cookies.delete('access_token', { path: '/' });
			return null;
		}
		throw err;
	}
});

export const login = form(
	v.object({
		phone: phoneSchema,
		password: v.pipe(v.string(), v.minLength(1, 'A palavra-passe é obrigatória.')),
		next: v.optional(v.string())
	}),
	async (data, issue) => {
		let auth: AuthResponse;
		try {
			auth = await apiFetch<AuthResponse>('/auth/login', {
				method: 'POST',
				body: { phone: normalizeAngolaPhone(data.phone), password: data.password }
			});
		} catch (err) {
			if (err instanceof ApiError && err.status === 401) {
				invalid(issue.password('Telemóvel ou palavra-passe incorrectos.'));
			}
			if (err instanceof ApiError && err.status === 403) {
				invalid(issue.phone('Esta conta encontra-se suspensa.'));
			}
			invalid(err instanceof ApiError ? err.message : 'Não foi possível iniciar sessão.');
		}

		setSessionCookie(auth.access_token);
		redirect(303, safeNext(data.next));
	}
);

export const register = form(
	v.object({
		name: v.pipe(v.string(), v.minLength(3, 'O nome deve ter pelo menos 3 caracteres.')),
		phone: phoneSchema,
		email: v.optional(v.string()),
		password: v.pipe(
			v.string(),
			v.minLength(6, 'A palavra-passe deve ter pelo menos 6 caracteres.')
		),
		next: v.optional(v.string())
	}),
	async (data, issue) => {
		const email = data.email?.trim();
		if (email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
			invalid(issue.email('Insira um email válido.'));
		}

		let auth: AuthResponse;
		try {
			auth = await apiFetch<AuthResponse>('/auth/register', {
				method: 'POST',
				body: {
					name: data.name.trim(),
					phone: normalizeAngolaPhone(data.phone),
					email: email || null,
					password: data.password
				}
			});
		} catch (err) {
			if (err instanceof ApiError && err.status === 409) {
				invalid(issue.phone('Já existe uma conta com este número.'));
			}
			invalid(err instanceof ApiError ? err.message : 'Não foi possível criar a conta.');
		}

		setSessionCookie(auth.access_token);
		redirect(303, safeNext(data.next));
	}
);

export const logout = form(async () => {
	const event = getRequestEvent();
	event.cookies.delete('access_token', { path: '/' });
	await getCurrentUser().refresh();
	redirect(303, '/');
});
