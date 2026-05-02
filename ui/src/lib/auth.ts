let csrfToken: string | null = null;

export async function fetchCsrfToken(): Promise<string> {
	if (csrfToken) return csrfToken;
	const res = await fetch("/api/auth/csrf-token", { credentials: "include" });
	const data = await res.json();
	csrfToken = data.csrfToken;
	return csrfToken;
}

async function authFetch(url: string, options: RequestInit = {}): Promise<Response> {
	const token = await fetchCsrfToken();
	return fetch(url, {
		...options,
		credentials: "include",
		headers: {
			"Content-Type": "application/json",
			"X-CSRF-Token": token,
			...options.headers,
		},
	});
}

export async function logout(): Promise<Response> {
	return authFetch("/api/auth/logout", { method: "POST" });
}

export async function getMe(): Promise<Response> {
	return fetch("/api/auth/me", {
		credentials: "include",
		headers: { "Content-Type": "application/json" },
	});
}
