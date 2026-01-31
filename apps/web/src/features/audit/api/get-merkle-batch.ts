import { apiClient } from "@/lib/api-client";
import type { MerkleBatch } from "./build-merkle-tree";

interface MerkleBatchResponse {
	batch: MerkleBatch;
}

export const getMerkleBatch = async (batchId: string): Promise<MerkleBatchResponse> => {
	return apiClient(`/api/audit/merkle/${batchId}`);
};

