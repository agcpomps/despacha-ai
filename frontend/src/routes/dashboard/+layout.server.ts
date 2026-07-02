import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

// Protege /dashboard e todas as sub-rotas (ex: /dashboard/editar/[id]) em
// qualquer navegação, incluindo client-side.
export const load: LayoutServerLoad = ({ locals, url }) => {
	if (!locals.isAuthenticated) {
		redirect(303, `/login?next=${encodeURIComponent(url.pathname + url.search)}`);
	}
};
