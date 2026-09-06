import { error } from '@sveltejs/kit';
import { apiFetch, ApiError } from '$lib/server/api';
import type { Listing } from '$lib/types';
import type { PageServerLoad } from './$types';

// Carrega o anúncio no servidor para que as meta tags (og:*, JSON-LD) saiam
// no HTML inicial — os crawlers do WhatsApp/Facebook não executam JavaScript,
// por isso o preview de partilha depende disto.
export const load: PageServerLoad = async ({ params }) => {
	try {
		const listing = await apiFetch<Listing>(`/listings/${encodeURIComponent(params.id)}`);
		return { listing };
	} catch (err) {
		if (err instanceof ApiError && err.status === 404) {
			error(404, 'Anúncio não encontrado.');
		}
		throw err;
	}
};
