import { getMe, logout } from "$lib/auth";

export interface User {
	id: string;
	email: string;
	username: string;
}

class AuthStore {
	user: User | null = $state(null);
	loading: boolean = $state(true);

	async init() {
		try {
			const res = await getMe();
			if (res.ok) {
				const data = await res.json();
				this.user = data.user;
			} else {
				this.user = null;
			}
		} catch {
			this.user = null;
		} finally {
			this.loading = false;
		}
	}

	async logout() {
		try {
			await logout();
		} catch {
		} finally {
			this.user = null;
		}
	}
}

export const auth = new AuthStore();
