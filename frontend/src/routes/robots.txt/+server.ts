import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ url }) => {
	const body = `User-agent: *
Allow: /
Disallow: /dashboard
Disallow: /admin
Disallow: /login
Disallow: /register

Sitemap: ${url.origin}/sitemap.xml
`;

	return new Response(body, {
		headers: { 'Content-Type': 'text/plain' }
	});
};
