<script lang="ts">
	let {
		open = $bindable(false),
		title,
		message,
		confirmLabel = 'Confirmar',
		cancelLabel = 'Cancelar',
		danger = false,
		onconfirm
	}: {
		open?: boolean;
		title: string;
		message: string;
		confirmLabel?: string;
		cancelLabel?: string;
		danger?: boolean;
		onconfirm: () => void;
	} = $props();

	function cancel() {
		open = false;
	}

	function confirm() {
		open = false;
		onconfirm();
	}

	function handleKeydown(event: KeyboardEvent) {
		if (open && event.key === 'Escape') cancel();
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
		<button
			type="button"
			class="absolute inset-0 cursor-default bg-neutral-900/50 backdrop-blur-sm"
			aria-label="Fechar"
			onclick={cancel}
		></button>

		<div
			role="dialog"
			aria-modal="true"
			aria-labelledby="confirm-title"
			class="relative w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl shadow-neutral-900/20"
		>
			<h2 id="confirm-title" class="text-lg font-bold text-neutral-900">{title}</h2>
			<p class="mt-2 text-sm leading-relaxed text-neutral-600">{message}</p>

			<div class="mt-6 flex justify-end gap-3">
				<button
					type="button"
					onclick={cancel}
					class="h-10 rounded-full border border-neutral-200 px-5 text-sm font-medium text-neutral-700 transition hover:border-neutral-300"
				>
					{cancelLabel}
				</button>
				<button
					type="button"
					onclick={confirm}
					class="h-10 rounded-full px-5 text-sm font-bold text-white transition {danger
						? 'bg-red-600 hover:bg-red-700'
						: 'bg-brand-700 hover:bg-brand-800'}"
				>
					{confirmLabel}
				</button>
			</div>
		</div>
	</div>
{/if}
