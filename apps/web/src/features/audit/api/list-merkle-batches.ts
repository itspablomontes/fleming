import { apiClient } from "@/lib/api-client";
import type { MerkleBatch } from "./build-merkle-tree";

interface ListMerkleBatchesResponse {
	batches: MerkleBatch[];
}

export const listMerkleBatches = async (params?: {
	limit?: number;
	offset?: number;
}): Promise<ListMerkleBatchesResponse> => {
	const search = new URLSearchParams();
	if (params?.limit !== undefined) search.set("limit", String(params.limit));
	if (params?.offset !== undefined) search.set("offset", String(params.offset));
	const qs = search.toString();
	return apiClient(`/api/audit/merkle/batches${qs ? `?${qs}` : ""}`);
};

