import * as v from 'valibot';
import { command, form, query, getRequestEvent } from '$app/server';
import { error, redirect, invalid } from '@sveltejs/kit';
import { apiFetch, ApiError } from '$lib/server/api';
import type { Listing, PaginatedListings } from '$lib/types';

const filtersSchema = v.object({
	search: v.optional(v.string()),
	category_id: v.optional(v.string()),
	province: v.optional(v.string()),
	city: v.optional(v.string()),
	min_price: v.optional(v.number()),
	max_price: v.optional(v.number()),
	sort: v.optional(v.picklist(['newest', 'oldest', 'price_asc', 'price_desc'])),
	status: v.optional(v.picklist(['active', 'sold', 'paused'])),
	featured: v.optional(v.boolean()),
	page: v.optional(v.number()),
	limit: v.optional(v.number())
});

type Filters = v.InferOutput<typeof filtersSchema>;

function toSearchParams(filters: Filters) {
	return {
		search: filters.search,
		category_id: filters.category_id,
		province: filters.province,
		city: filters.city,
		min_price: filters.min_price,
		max_price: filters.max_price,
		sort: filters.sort,
		status: filters.status,
		featured: filters.featured ? 'true' : undefined,
		page: filters.page,
		limit: filters.limit
	};
}

export const getListings = query(filtersSchema, async (filters) => {
	return await apiFetch<PaginatedListings>('/listings', {
		searchParams: toSearchParams(filters)
	});
});

export const getListing = query(v.string(), async (id) => {
	try {
		return await apiFetch<Listing>(`/listings/${encodeURIComponent(id)}`);
	} catch (err) {
		if (err instanceof ApiError && err.status === 404) {
			error(404, 'Anúncio não encontrado.');
		}
		throw err;
	}
});

export const getMyListings = query(filtersSchema, async (filters) => {
	return await apiFetch<PaginatedListings>('/me/listings', {
		auth: true,
		searchParams: toSearchParams(filters)
	});
});

// an empty <input type="file"> still submits one empty File entry — drop those
const imagesSchema = v.optional(
	v.pipe(
		v.union([v.instance(File), v.array(v.instance(File))]),
		v.transform((value) =>
			(Array.isArray(value) ? value : [value]).filter((file) => file.size > 0 && file.name)
		),
		v.maxLength(8, 'Máximo de 8 imagens por anúncio.')
	)
);

const priceSchema = v.pipe(
	v.string(),
	v.transform((value) => Number(value.replace(/\s/g, '').replace(',', '.'))),
	v.number('O preço deve ser um número válido.'),
	v.minValue(1, 'O preço deve ser maior que zero.')
);

export const createListing = form(
	v.object({
		title: v.pipe(v.string(), v.minLength(3, 'O título deve ter pelo menos 3 caracteres.')),
		description: v.pipe(
			v.string(),
			v.minLength(10, 'A descrição deve ter pelo menos 10 caracteres.')
		),
		price: priceSchema,
		category_id: v.optional(v.string()),
		province: v.pipe(v.string(), v.minLength(1, 'Selecione a província.')),
		city: v.optional(v.string()),
		address_reference: v.optional(v.string()),
		whatsapp_phone: v.optional(v.string()),
		phone: v.optional(v.string()),
		condition: v.picklist(['new', 'used'], 'Selecione o estado do artigo.'),
		images: imagesSchema
	}),
	async (data) => {
		const formData = new FormData();
		formData.set('title', data.title.trim());
		formData.set('description', data.description.trim());
		formData.set('price', String(data.price));
		formData.set('province', data.province);
		formData.set('condition', data.condition);
		if (data.category_id) formData.set('category_id', data.category_id);
		if (data.city) formData.set('city', data.city);
		if (data.address_reference) formData.set('address_reference', data.address_reference);
		if (data.whatsapp_phone) formData.set('whatsapp_phone', data.whatsapp_phone);
		if (data.phone) formData.set('phone', data.phone);
		// O File reconstruído pelo protocolo binário das remote functions não é um
		// Blob nativo que o fetch reconheça — ao serializar o FormData seria convertido
		// para a string "[object Object]" (15 bytes). Reconstruímos um File nativo a
		// partir dos bytes reais para garantir o upload correcto.
		for (const file of data.images ?? []) {
			const buffer = await file.arrayBuffer();
			formData.append('images', new File([buffer], file.name, { type: file.type }));
		}

		let listing: Listing;
		try {
			listing = await apiFetch<Listing>('/listings', {
				method: 'POST',
				auth: true,
				formData
			});
		} catch (err) {
			invalid(err instanceof ApiError ? err.message : 'Não foi possível publicar o anúncio.');
		}

		redirect(303, `/anuncio/${listing.id}`);
	}
);

export const updateListing = form(
	v.object({
		id: v.string(),
		title: v.pipe(v.string(), v.minLength(3, 'O título deve ter pelo menos 3 caracteres.')),
		description: v.pipe(
			v.string(),
			v.minLength(10, 'A descrição deve ter pelo menos 10 caracteres.')
		),
		price: priceSchema,
		province: v.pipe(v.string(), v.minLength(1, 'Selecione a província.')),
		city: v.optional(v.string()),
		address_reference: v.optional(v.string()),
		whatsapp_phone: v.optional(v.string()),
		phone: v.optional(v.string()),
		condition: v.picklist(['new', 'used'], 'Selecione o estado do artigo.')
	}),
	async (data) => {
		try {
			await apiFetch<Listing>(`/listings/${encodeURIComponent(data.id)}`, {
				method: 'PUT',
				auth: true,
				body: {
					title: data.title.trim(),
					description: data.description.trim(),
					price: data.price,
					province: data.province,
					city: data.city || null,
					address_reference: data.address_reference || null,
					whatsapp_phone: data.whatsapp_phone || null,
					phone: data.phone || null,
					condition: data.condition
				}
			});
		} catch (err) {
			invalid(err instanceof ApiError ? err.message : 'Não foi possível actualizar o anúncio.');
		}

		redirect(303, '/dashboard');
	}
);

export const setListingStatus = command(
	v.object({
		id: v.string(),
		status: v.picklist(['active', 'sold', 'paused'])
	}),
	async ({ id, status }) => {
		await apiFetch<Listing>(`/listings/${encodeURIComponent(id)}`, {
			method: 'PUT',
			auth: true,
			body: { status }
		});
	}
);

// --- admin-only promotion commands (backend enforces the admin role) ---

export const featureListing = command(
	v.object({ id: v.string(), days: v.number() }),
	async ({ id, days }) => {
		await apiFetch(`/admin/listings/${encodeURIComponent(id)}/feature`, {
			method: 'PUT',
			auth: true,
			body: { days }
		});
		await getListing(id).refresh();
	}
);

export const unfeatureListing = command(v.string(), async (id) => {
	await apiFetch(`/admin/listings/${encodeURIComponent(id)}/feature`, {
		method: 'DELETE',
		auth: true
	});
	await getListing(id).refresh();
});

export const bumpListing = command(v.string(), async (id) => {
	await apiFetch(`/admin/listings/${encodeURIComponent(id)}/bump`, {
		method: 'PUT',
		auth: true
	});
	await getListing(id).refresh();
});

export const deleteListing = command(v.string(), async (id) => {
	const { cookies } = getRequestEvent();
	if (!cookies.get('access_token')) error(401, 'Sessão expirada.');

	await apiFetch<void>(`/listings/${encodeURIComponent(id)}`, {
		method: 'DELETE',
		auth: true
	});
});
