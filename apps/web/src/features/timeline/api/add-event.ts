import { apiClient } from "@/lib/api-client";
import type { UploadResponse } from "../types";

export interface UploadPayload {
	file?: File | Blob;
	eventType: string;
	title: string;
	description: string;
	provider?: string;
	date: string;
	metadata?: Record<string, unknown>;

	isEncrypted?: boolean;
	wrappedKey?: string;
	blobRef?: string;
}

export const addEvent = (payload: UploadPayload): Promise<UploadResponse> => {
	const formData = new FormData();
	if (payload.file) {
		formData.append("file", payload.file);
	}
	formData.append("eventType", payload.eventType);
	formData.append("title", payload.title);
	formData.append("description", payload.description);
	formData.append("date", payload.date);
	
	if (payload.provider) {
		formData.append("provider", payload.provider);
	}

	if (payload.metadata) {
		formData.append("metadata", JSON.stringify(payload.metadata));
	}
	
	if (payload.isEncrypted) {
		formData.append("isEncrypted", "true");
	}
	
	if (payload.wrappedKey) {
		formData.append("wrappedKey", payload.wrappedKey);
	}

	return apiClient("/api/timeline/events", {
		method: "POST",
		body: formData,
	});
};

export interface CorrectionPayload extends Omit<UploadPayload, "file"> {
	id: string;
	file?: File | Blob;
}

export const correctEvent = (payload: CorrectionPayload): Promise<UploadResponse> => {
	const formData = new FormData();

	if (payload.file) {
		formData.append("file", payload.file);
	}
	
	formData.append("eventType", payload.eventType);
	formData.append("title", payload.title);
	formData.append("description", payload.description);
	formData.append("date", payload.date);
	
	if (payload.provider) {
		formData.append("provider", payload.provider);
	}

	if (payload.metadata) {
		formData.append("metadata", JSON.stringify(payload.metadata));
	}
	
	if (payload.isEncrypted) {
		formData.append("isEncrypted", "true");
	}
	
	if (payload.wrappedKey) {
		formData.append("wrappedKey", payload.wrappedKey);
	}

	return apiClient(`/api/timeline/events/${payload.id}/correction`, {
		method: "POST",
		body: formData,
	});
};
