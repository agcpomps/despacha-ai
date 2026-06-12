import { redirect, type Handle } from '@sveltejs/kit';

const PROTECTED_PATHS = ['/anunciar', '/dashboard', '/admin'];

function isProtected(pathname: string) {
	return PROTECTED_PATHS.some((path) => pathname === path || pathname.startsWith(`${path}/`));
}

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('access_token');
	event.locals.isAuthenticated = Boolean(token);

	if (isProtected(event.url.pathname) && !token) {
		const next = encodeURIComponent(`${event.url.pathname}${event.url.search}`);
		throw redirect(303, `/login?next=${next}`);
	}

	return resolve(event);
};
