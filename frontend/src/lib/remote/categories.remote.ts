import { query } from '$app/server';
import { apiFetch } from '$lib/server/api';
import type { Category } from '$lib/types';

export const getCategories = query(async () => {
	return await apiFetch<Category[]>('/categories');
});
