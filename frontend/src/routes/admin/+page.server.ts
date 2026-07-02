import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals, url }) => {
	if (!locals.isAuthenticated) {
		redirect(303, `/login?next=${encodeURIComponent(url.pathname + url.search)}`);
	}
};
