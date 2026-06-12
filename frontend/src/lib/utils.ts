export function formatPrice(value: number) {
	return `${Math.round(value).toLocaleString('pt-PT').replace(/\s/g, '.')} Kz`;
}

export function locationText(province: string, city?: string | null) {
	return city ? `${city}, ${province}` : province;
}

export function timeAgo(iso: string) {
	const seconds = Math.floor((Date.now() - new Date(iso).getTime()) / 1000);
	if (seconds < 60) return 'agora mesmo';
	const minutes = Math.floor(seconds / 60);
	if (minutes < 60) return `há ${minutes} min`;
	const hours = Math.floor(minutes / 60);
	if (hours < 24) return `há ${hours} h`;
	const days = Math.floor(hours / 24);
	if (days < 30) return days === 1 ? 'há 1 dia' : `há ${days} dias`;
	const months = Math.floor(days / 30);
	if (months < 12) return months === 1 ? 'há 1 mês' : `há ${months} meses`;
	const years = Math.floor(months / 12);
	return years === 1 ? 'há 1 ano' : `há ${years} anos`;
}

export function normalizeAngolaPhone(phone: string) {
	const digits = phone.replace(/\D/g, '');
	const local = digits.startsWith('244') ? digits.slice(3) : digits;
	return `+244${local}`;
}

export function isValidAngolaPhone(phone: string) {
	const digits = phone.replace(/\D/g, '');
	const local = digits.startsWith('244') ? digits.slice(3) : digits;
	return /^9\d{8}$/.test(local);
}

export const PROVINCES = [
	'Bengo',
	'Benguela',
	'Bié',
	'Cabinda',
	'Cuando Cubango',
	'Cuanza Norte',
	'Cuanza Sul',
	'Cunene',
	'Huambo',
	'Huíla',
	'Luanda',
	'Lunda Norte',
	'Lunda Sul',
	'Malanje',
	'Moxico',
	'Namibe',
	'Uíge',
	'Zaire'
];

export const CONDITION_LABELS: Record<string, string> = {
	new: 'Novo',
	used: 'Usado'
};

export const STATUS_LABELS: Record<string, string> = {
	active: 'Activo',
	sold: 'Vendido',
	paused: 'Pausado'
};
