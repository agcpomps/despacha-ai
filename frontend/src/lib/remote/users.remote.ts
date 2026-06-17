import * as v from 'valibot';
import { command, query } from '$app/server';
import { apiFetch } from '$lib/server/api';
import type { PaginatedUsers, User } from '$lib/types';

export const getAdminUsers = query(
	v.object({
		search: v.optional(v.string()),
		page: v.optional(v.number()),
		limit: v.optional(v.number())
	}),
	async (filters) => {
		return await apiFetch<PaginatedUsers>('/admin/users', {
			auth: true,
			searchParams: {
				search: filters.search,
				page: filters.page,
				limit: filters.limit
			}
		});
	}
);

export const resetUserPassword = command(v.string(), async (id) => {
	return await apiFetch<{ password: string }>(
		`/admin/users/${encodeURIComponent(id)}/password`,
		{ method: 'PUT', auth: true }
	);
});

export const setUserRole = command(
	v.object({
		id: v.string(),
		role: v.picklist(['user', 'moderator', 'admin'])
	}),
	async ({ id, role }) => {
		return await apiFetch<User>(`/admin/users/${encodeURIComponent(id)}/role`, {
			method: 'PUT',
			auth: true,
			body: { role }
		});
	}
);
