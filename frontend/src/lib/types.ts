export type UserRole = 'user' | 'moderator' | 'admin';

export type User = {
	id: string;
	name: string;
	phone: string;
	email?: string;
	avatar_url?: string;
	role: UserRole;
	status: string;
	is_verified: boolean;
	created_at?: string;
};

export type PaginatedUsers = {
	data: User[];
	page: number;
	limit: number;
	total: number;
	total_pages: number;
};

export type AuthResponse = {
	access_token: string;
	user: User;
};

export type Category = {
	id: string;
	name: string;
	slug: string;
	parent_id?: string;
	children?: Category[];
};

export type ListingImage = {
	id: string;
	image_url: string;
	position: number;
};

export type ListingSeller = {
	id: string;
	name: string;
	phone: string;
	avatar_url?: string;
};

export type ListingCategory = {
	id: string;
	name: string;
	slug: string;
};

export type ListingCondition = 'new' | 'used';
export type ListingStatus = 'active' | 'sold' | 'paused' | 'deleted';
export type ListingSort = 'newest' | 'oldest' | 'price_asc' | 'price_desc';

export type Listing = {
	id: string;
	user_id: string;
	category_id?: string;
	title: string;
	description: string;
	price: number;
	currency: string;
	province: string;
	city?: string;
	address_reference?: string;
	whatsapp_phone?: string;
	phone?: string;
	condition: ListingCondition;
	status: ListingStatus;
	views_count: number;
	is_featured: boolean;
	featured_until?: string;
	seller?: ListingSeller;
	category?: ListingCategory;
	images: ListingImage[];
	created_at: string;
	updated_at: string;
};

export type PaginatedListings = {
	data: Listing[];
	page: number;
	limit: number;
	total: number;
	total_pages: number;
};

export type ListingFilters = {
	search?: string;
	category_id?: string;
	province?: string;
	city?: string;
	min_price?: number;
	max_price?: number;
	sort?: ListingSort;
	status?: Exclude<ListingStatus, 'deleted'>;
	featured?: boolean;
	page?: number;
	limit?: number;
};
