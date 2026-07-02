import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

// Corre em todas as navegações (servidor e client-side), garantindo que o
// utilizador é enviado para o login ANTES de ver o formulário — senão preenchia
// tudo e só perdia os dados ao submeter. O hooks.server.ts só protege pedidos
// diretos ao servidor, não a navegação dentro da app.
export const load: PageServerLoad = ({ locals, url }) => {
	if (!locals.isAuthenticated) {
		redirect(303, `/login?next=${encodeURIComponent(url.pathname + url.search)}`);
	}
};
