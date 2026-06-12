import type { RequestHandler } from './$types';
import { env } from '$env/dynamic/public';
import type { PaginatedListings, Category } from '$lib/types';

const API_BASE = env.PUBLIC_API_BASE_URL || 'http://localhost:8080';

export const GET: RequestHandler = async ({ url, fetch }) => {
	const origin = url.origin;

	const staticPages = ['', '/search', '/categories', '/anunciar', '/login', '/register'];

	let listings: PaginatedListings | null = null;
	let categories: Category[] = [];

	try {
		const [listingsRes, categoriesRes] = await Promise.all([
			fetch(`${API_BASE}/api/v1/listings?limit=50&sort=newest`),
			fetch(`${API_BASE}/api/v1/categories`)
		]);
		if (listingsRes.ok) listings = await listingsRes.json();
		if (categoriesRes.ok) categories = await categoriesRes.json();
	} catch {
		// sitemap continua só com as páginas estáticas
	}

	const urls: string[] = [];

	for (const path of staticPages) {
		urls.push(`<url><loc>${origin}${path}</loc><changefreq>daily</changefreq></url>`);
	}

	for (const category of categories) {
		urls.push(
			`<url><loc>${origin}/search?category_id=${category.id}</loc><changefreq>daily</changefreq></url>`
		);
		for (const child of category.children ?? []) {
			urls.push(
				`<url><loc>${origin}/search?category_id=${child.id}</loc><changefreq>daily</changefreq></url>`
			);
		}
	}

	for (const listing of listings?.data ?? []) {
		urls.push(
			`<url><loc>${origin}/anuncio/${listing.id}</loc><lastmod>${listing.updated_at.slice(0, 10)}</lastmod></url>`
		);
	}

	const xml = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${urls.join('\n')}
</urlset>`;

	return new Response(xml, {
		headers: {
			'Content-Type': 'application/xml',
			'Cache-Control': 'public, max-age=3600'
		}
	});
};
