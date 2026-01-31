import { apiClient } from "@/lib/api-client";
import type { MerkleBatch } from "./build-merkle-tree";

interface AnchorMerkleBatchResponse {
	batch: MerkleBatch;
}

export const anchorMerkleBatch = async (
	batchId: string,
): Promise<AnchorMerkleBatchResponse> => {
	return apiClient(`/api/audit/merkle/${batchId}/anchor`, { method: "POST" });
};

