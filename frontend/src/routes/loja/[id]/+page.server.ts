import { error } from '@sveltejs/kit';
import { apiFetch, ApiError } from '$lib/server/api';
import type { PaginatedListings, PublicSeller } from '$lib/types';
import type { PageServerLoad } from './$types';

// Página pública da loja do vendedor. Carregada no servidor para que o link
// partilhado no WhatsApp tenha preview (nome, foto) e seja indexável.
export const load: PageServerLoad = async ({ params, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;

	try {
		const [seller, listings] = await Promise.all([
			apiFetch<PublicSeller>(`/users/${encodeURIComponent(params.id)}/public`),
			apiFetch<PaginatedListings>('/listings', {
				searchParams: { user_id: params.id, page, limit: 16, sort: 'newest' }
			})
		]);

		return { seller, listings };
	} catch (err) {
		if (err instanceof ApiError && (err.status === 404 || err.status === 400)) {
			error(404, 'Loja não encontrada.');
		}
		throw err;
	}
};
