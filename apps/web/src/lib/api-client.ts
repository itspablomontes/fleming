export const apiClient = async (
    endpoint: string,
    { body, ...customConfig }: Omit<RequestInit, "body"> & { body?: any } = {},
) => {
    const isFormData = body instanceof FormData;
    const headers: HeadersInit = isFormData
        ? {}
        : { "Content-Type": "application/json" };

    const config: RequestInit = {
        method: body ? "POST" : "GET",
        ...customConfig,
        headers: {
            ...headers,
            ...customConfig.headers,
        },
        credentials: "include",
    };

    if (body) {
        config.body = isFormData ? body : JSON.stringify(body);
    }

    const response = await fetch(endpoint, config);

    if (!response.ok) {
        const errorMessage = await response.text();
        return Promise.reject(new Error(errorMessage || response.statusText));
    }

    if (response.status === 204) {
        return null;
    }

    return response.json();
};
